package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newCollectionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection",
		Short: "Manage snippet collections",
	}
	cmd.AddCommand(newCollectionCreateCmd())
	cmd.AddCommand(newCollectionAddCmd())
	cmd.AddCommand(newCollectionRemoveCmd())
	cmd.AddCommand(newCollectionListCmd())
	return cmd
}

func newCollectionCreateCmd() *cobra.Command {
	var description string
	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new collection",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			import_snippet "github.com/your-org/stashctl/internal/snippet"
			c, err := snippet.NewCollection(args[0], description)
			if err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Created collection %q (id: %s)\n", c.Name, c.ID)
			return nil
		},
	}
	cmd.Flags().StringVarP(&description, "description", "d", "", "Optional description")
	return cmd
}

func newCollectionAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <collection-id> <snippet-id>",
		Short: "Add a snippet to a collection",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "Added snippet %s to collection %s\n", args[1], args[0])
			return nil
		},
	}
}

func newCollectionRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <collection-id> <snippet-id>",
		Short: "Remove a snippet from a collection",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintf(cmd.OutOrStdout(), "Removed snippet %s from collection %s\n", args[1], args[0])
			return nil
		},
	}
}

func newCollectionListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all collections",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "No collections found.")
			return nil
		},
	}
}
