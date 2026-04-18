package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/stashctl/internal/config"
	"github.com/user/stashctl/internal/export"
	"github.com/user/stashctl/internal/format"
	"github.com/user/stashctl/internal/search"
	"github.com/user/stashctl/internal/snippet"
	"github.com/user/stashctl/internal/store"
)

func newRootCmd(cfg *config.Config, st *store.Store) *cobra.Command {
	root := &cobra.Command{
		Use:   "stashctl",
		Short: "Manage and search local code snippets",
	}

	root.AddCommand(newAddCmd(st))
	root.AddCommand(newListCmd(cfg, st))
	root.AddCommand(newSearchCmd(cfg, st))
	root.AddCommand(newDeleteCmd(st))
	root.AddCommand(newExportCmd(st))

	return root
}

func newAddCmd(st *store.Store) *cobra.Command {
	var tags []string
	var lang string
	cmd := &cobra.Command{
		Use:   "add <title> <body>",
		Short: "Add a new snippet",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			sn := snippet.New(args[0], args[1], lang, tags)
			if err := st.Add(sn); err != nil {
				return fmt.Errorf("add: %w", err)
			}
			fmt.Printf("Added snippet %s\n", sn.ID)
			return nil
		},
	}
	cmd.Flags().StringSliceVarP(&tags, "tags", "t", nil, "Comma-separated tags")
	cmd.Flags().StringVarP(&lang, "lang", "l", "", "Language")
	return cmd
}

func newListCmd(cfg *config.Config, st *store.Store) *cobra.Command {
	var tags []string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List snippets",
		RunE: func(cmd *cobra.Command, args []string) error {
			snippets, err := st.FilterByTags(tags)
			if err != nil {
				return err
			}
			for _, sn := range snippets {
				fmt.Println(format.SnippetSummary(sn, cfg.ColorOutput))
			}
			return nil
		},
	}
	cmd.Flags().StringSliceVarP(&tags, "tags", "t", nil, "Filter by tags")
	return cmd
}

func newSearchCmd(cfg *config.Config, st *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "search <query>",
		Short: "Search snippets by query",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			all, err := st.List()
			if err != nil {
				return err
			}
			results := search.ByQuery(all, args[0])
			for _, sn := range results {
				fmt.Println(format.SnippetSummary(sn, cfg.ColorOutput))
			}
			return nil
		},
	}
}

func newDeleteCmd(st *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a snippet by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := st.Delete(args[0]); err != nil {
				return fmt.Errorf("delete: %w", err)
			}
			fmt.Printf("Deleted snippet %s\n", args[0])
			return nil
		},
	}
}

func newExportCmd(st *store.Store) *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export snippets to stdout",
		RunE: func(cmd *cobra.Command, args []string) error {
			all, err := st.List()
			if err != nil {
				return err
			}
			out, err := export.Snippets(all, strings.ToLower(format))
			if err != nil {
				return err
			}
			fmt.Print(out)
			return nil
		},
	}
	cmd.Flags().StringVarP(&format, "format", "f", "json", "Output format: json or markdown")
	return cmd
}
