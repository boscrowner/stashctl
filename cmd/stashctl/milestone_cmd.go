package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/user/stashctl/internal/snippet"
	"github.com/user/stashctl/internal/store"
)

func newMilestoneCmd(st *store.Store) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "milestone",
		Short: "Manage snippet milestones",
	}
	cmd.AddCommand(
		newMilestoneAddCmd(st),
		newMilestoneListCmd(st),
		newMilestoneCompleteCmd(st),
		newMilestoneRemoveCmd(st),
	)
	return cmd
}

func newMilestoneAddCmd(st *store.Store) *cobra.Command {
	var dueStr, note string
	cmd := &cobra.Command{
		Use:   "add <snippet-id> <name>",
		Short: "Add a milestone to a snippet",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			var due time.Time
			if dueStr != "" {
				var err error
				due, err = time.Parse(time.DateOnly, dueStr)
				if err != nil {
					return fmt.Errorf("invalid due date %q: %w", dueStr, err)
				}
			}
			m, err := snippet.NewMilestone(args[0], args[1], note, due)
			if err != nil {
				return err
			}
			if err := st.AddMilestone(m); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "milestone added: %s\n", m.ID)
			return nil
		},
	}
	cmd.Flags().StringVar(&dueStr, "due", "", "due date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&note, "note", "", "optional note")
	return cmd
}

func newMilestoneListCmd(st *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "list <snippet-id>",
		Short: "List milestones for a snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			all, err := st.ListMilestones()
			if err != nil {
				return err
			}
			ms := snippet.MilestonesFor(all, args[0])
			if len(ms) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no milestones")
				return nil
			}
			for _, m := range ms {
				done := " "
				if m.Done {
					done = "✓"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "[%s] %s  %s\n", done, m.ID[:8], m.Name)
			}
			return nil
		},
	}
}

func newMilestoneCompleteCmd(st *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "complete <milestone-id>",
		Short: "Mark a milestone as done",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return st.CompleteMilestone(args[0])
		},
	}
}

func newMilestoneRemoveCmd(st *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <milestone-id>",
		Short: "Remove a milestone",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return st.RemoveMilestone(args[0])
		},
	}
}
