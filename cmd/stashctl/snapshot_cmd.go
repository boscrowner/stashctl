package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/user/stashctl/internal/snippet"
)

func newSnapshotCmd(env *cmdEnv) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Manage snippet snapshots",
	}
	cmd.AddCommand(
		newSnapshotTakeCmd(env),
		newSnapshotListCmd(env),
		newSnapshotRemoveCmd(env),
	)
	return cmd
}

func newSnapshotTakeCmd(env *cmdEnv) *cobra.Command {
	var label string
	cmd := &cobra.Command{
		Use:   "take <snippet-id>",
		Short: "Take a snapshot of a snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			snippets, err := env.store.List()
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
			sn, err := snippet.NewSnapshot(*found, label)
			if err != nil {
				return err
			}
			if err := env.store.SaveSnapshot(sn); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "snapshot %s taken\n", sn.ID)
			return nil
		},
	}
	cmd.Flags().StringVarP(&label, "label", "l", "", "optional label for the snapshot")
	return cmd
}

func newSnapshotListCmd(env *cmdEnv) *cobra.Command {
	return &cobra.Command{
		Use:   "list <snippet-id>",
		Short: "List snapshots for a snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			all, err := env.store.ListSnapshots(args[0])
			if err != nil {
				return err
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "ID\tLABEL\tCREATED")
			for _, sn := range all {
				fmt.Fprintf(w, "%s\t%s\t%s\n", sn.ID, sn.Label, sn.CreatedAt.Format("2006-01-02 15:04"))
			}
			return w.Flush()
		},
	}
}

func newSnapshotRemoveCmd(env *cmdEnv) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <snapshot-id>",
		Short: "Remove a snapshot",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := env.store.RemoveSnapshot(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "snapshot %s removed\n", args[0])
			return nil
		},
	}
}

// ensure snapshot_cmd compiles even without os import warning
var _ = os.Stderr
