package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/user/stashctl/internal/snippet"
	"github.com/user/stashctl/internal/store"
)

func newVersionCmd(s *store.Store) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Manage snippet versions",
	}
	cmd.AddCommand(newVersionSaveCmd(s))
	cmd.AddCommand(newVersionListCmd(s))
	cmd.AddCommand(newVersionShowCmd(s))
	return cmd
}

func newVersionSaveCmd(s *store.Store) *cobra.Command {
	var message string
	cmd := &cobra.Command{
		Use:   "save <snippet-id>",
		Short: "Save a version snapshot of a snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			snippets, err := s.List()
			if err != nil {
				return err
			}
			var found *snippet.Snippet
			for i := range snippets {
				if snippets[i].ID == args[0] {
					found = &snippets[i]
					break
				}
			}
			if found == nil {
				return fmt.Errorf("snippet %q not found", args[0])
			}
			v, err := snippet.NewVersion(found.ID, found.Content, message)
			if err != nil {
				return err
			}
			if err := s.AddVersion(v); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "version %s saved\n", v.ID)
			return nil
		},
	}
	cmd.Flags().StringVarP(&message, "message", "m", "", "optional version message")
	return cmd
}

func newVersionListCmd(s *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "list <snippet-id>",
		Short: "List all versions for a snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			versions, err := s.ListVersions(args[0])
			if err != nil {
				return err
			}
			if len(versions) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no versions found")
				return nil
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tCREATED\tMESSAGE")
			for _, v := range versions {
				fmt.Fprintf(w, "%s\t%s\t%s\n", v.ID, v.CreatedAt.Format("2006-01-02 15:04"), v.Message)
			}
			return w.Flush()
		},
	}
}

func newVersionShowCmd(s *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "show <version-id>",
		Short: "Show the content of a specific version",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			versions, err := s.LoadVersions()
			if err != nil {
				return err
			}
			for _, v := range versions {
				if v.ID == args[0] {
					fmt.Fprintln(cmd.OutOrStdout(), v.Content)
					return nil
				}
			}
			fmt.Fprintf(os.Stderr, "version %q not found\n", args[0])
			return fmt.Errorf("version not found")
		},
	}
}
