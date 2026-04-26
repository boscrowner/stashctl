package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/stashctl/internal/snippet"
	"github.com/user/stashctl/internal/store"
)

func newWorkspaceCmd(s *store.Store) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "Manage workspaces",
	}
	cmd.AddCommand(
		newWorkspaceCreateCmd(s),
		newWorkspaceAddCmd(s),
		newWorkspaceRemoveCmd(s),
		newWorkspaceListCmd(s),
	)
	return cmd
}

func newWorkspaceCreateCmd(s *store.Store) *cobra.Command {
	var description string
	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new workspace",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			w, err := snippet.NewWorkspace(args[0], description)
			if err != nil {
				return err
			}
			if err := s.SaveWorkspace(w); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "workspace %q created (%s)\n", w.Name, w.ID)
			return nil
		},
	}
	cmd.Flags().StringVarP(&description, "description", "d", "", "workspace description")
	return cmd
}

func newWorkspaceAddCmd(s *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "add <workspace-id> <snippet-id>",
		Short: "Add a snippet to a workspace",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			w, err := s.GetWorkspace(args[0])
			if err != nil {
				return err
			}
			if err := w.AddSnippet(args[1]); err != nil {
				return err
			}
			if err := s.SaveWorkspace(w); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "snippet %q added to workspace %q\n", args[1], w.Name)
			return nil
		},
	}
}

func newWorkspaceRemoveCmd(s *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <workspace-id> <snippet-id>",
		Short: "Remove a snippet from a workspace",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			w, err := s.GetWorkspace(args[0])
			if err != nil {
				return err
			}
			if err := w.RemoveSnippet(args[1]); err != nil {
				return err
			}
			if err := s.SaveWorkspace(w); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "snippet %q removed from workspace %q\n", args[1], w.Name)
			return nil
		},
	}
}

func newWorkspaceListCmd(s *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all workspaces",
		RunE: func(cmd *cobra.Command, args []string) error {
			workspaces, err := s.ListWorkspaces()
			if err != nil {
				return err
			}
			if len(workspaces) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no workspaces found")
				return nil
			}
			for _, w := range workspaces {
				snippetCount := len(w.SnippetIDs)
				desc := ""
				if w.Description != "" {
					desc = " — " + w.Description
				}
				fmt.Fprintf(cmd.OutOrStdout(), "[%s] %s%s (%s)\n",
					w.ID[:8],
					w.Name,
					desc,
					snippetCountLabel(snippetCount),
				)
			}
			return nil
		},
	}
}

func snippetCountLabel(n int) string {
	if n == 1 {
		return "1 snippet"
	}
	return fmt.Sprintf("%d snippets", n)
}

// ensure strings import is used
var _ = strings.TrimSpace
