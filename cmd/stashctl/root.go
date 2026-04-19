package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/user/stashctl/internal/config"
	"github.com/user/stashctl/internal/export"
	"github.com/user/stashctl/internal/format"
	"github.com/user/stashctl/internal/search"
	"github.com/user/stashctl/internal/snippet"
	"github.com/user/stashctl/internal/store"
)

func newRootCmd() *cobra.Command {
	cfg, _ := config.Load(config.DefaultPath())

	root := &cobra.Command{
		Use:   "stashctl",
		Short: "Manage and search local code snippets",
	}

	root.AddCommand(
		newAddCmd(cfg),
		newListCmd(cfg),
		newSearchCmd(cfg),
		newDeleteCmd(cfg),
		newExportCmd(cfg),
		newEditCmd(cfg),
	)
	return root
}

func newAddCmd(cfg *config.Config) *cobra.Command {
	var tags, language string
	cmd := &cobra.Command{
		Use:   "add <title> <content>",
		Short: "Add a new snippet",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := store.New(cfg.StorePath)
			snip := snippet.New(args[0], args[1], language, snippet.ParseTags(tags))
			if err := s.Add(snip); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Added snippet %s\n", snip.ID)
			return nil
		},
	}
	cmd.Flags().StringVarP(&tags, "tags", "t", "", "comma-separated tags")
	cmd.Flags().StringVarP(&language, "lang", "l", "", "language")
	return cmd
}

func newListCmd(cfg *config.Config) *cobra.Command {
	var tags string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List snippets",
		RunE: func(cmd *cobra.Command, args []string) error {
			s := store.New(cfg.StorePath)
			snippets, err := s.List()
			if err != nil {
				return err
			}
			if tags != "" {
				snippets, err = s.FilterByTags(snippet.ParseTags(tags))
				if err != nil {
					return err
				}
			}
			for _, snip := range snippets {
				fmt.Fprintln(cmd.OutOrStdout(), format.SnippetSummary(snip, false))
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&tags, "tags", "t", "", "filter by tags")
	return cmd
}

func newSearchCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "search <query>",
		Short: "Search snippets by query",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := store.New(cfg.StorePath)
			snippets, err := s.List()
			if err != nil {
				return err
			}
			results := search.ByQuery(args[0], snippets)
			for _, r := range results {
				fmt.Fprintln(cmd.OutOrStdout(), format.SnippetSummary(r.Snippet, false))
			}
			return nil
		},
	}
}

func newDeleteCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a snippet by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := store.New(cfg.StorePath)
			if err := s.Delete(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Deleted %s\n", args[0])
			return nil
		},
	}
}

func newExportCmd(cfg *config.Config) *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export snippets to stdout",
		RunE: func(cmd *cobra.Command, args []string) error {
			s := store.New(cfg.StorePath)
			snippets, err := s.List()
			if err != nil {
				return err
			}
			out, err := export.Snippets(snippets, format)
			if err != nil {
				return err
			}
			fmt.Fprint(cmd.OutOrStdout(), out)
			return nil
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", "json", "output format: json|markdown")
	return cmd
}

func newEditCmd(cfg *config.Config) *cobra.Command {
	var tags, language, title, content string
	cmd := &cobra.Command{
		Use:   "edit <id>",
		Short: "Edit an existing snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := store.New(cfg.StorePath)
			snip, err := s.Get(args[0])
			if err != nil {
				return err
			}
			newTitle := snip.Title
			if title != "" {
				newTitle = title
			}
			newContent := snip.Content
			if content != "" {
				newContent = content
			}
			newLang := snip.Language
			if language != "" {
				newLang = language
			}
			newTags := snip.Tags
			if tags != "" {
				newTags = snippet.ParseTags(tags)
			}
			snip.Update(newTitle, newContent, newLang, newTags)
			if err := s.Update(snip); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Updated snippet %s\n", snip.ID)
			return nil
		},
	}
	cmd.Flags().StringVarP(&title, "title", "n", "", "new title")
	cmd.Flags().StringVarP(&content, "content", "c", "", "new content")
	cmd.Flags().StringVarP(&tags, "tags", "t", "", "new comma-separated tags")
	cmd.Flags().StringVarP(&language, "lang", "l", "", "new language")
	_ = strings.TrimSpace("") // keep import
	_ = os.Stderr
	return cmd
}
