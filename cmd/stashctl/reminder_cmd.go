package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/user/stashctl/internal/snippet"
	"github.com/user/stashctl/internal/store"
)

func newReminderCmd(st *store.Store) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reminder",
		Short: "Manage snippet review reminders",
	}
	cmd.AddCommand(newReminderSetCmd(st))
	cmd.AddCommand(newReminderListCmd(st))
	cmd.AddCommand(newReminderRemoveCmd(st))
	return cmd
}

func newReminderSetCmd(st *store.Store) *cobra.Command {
	var inDuration string
	var note string
	cmd := &cobra.Command{
		Use:   "set <snippet-id>",
		Short: "Set a reminder for a snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			d, err := time.ParseDuration(inDuration)
			if err != nil {
				return fmt.Errorf("invalid duration %q: %w", inDuration, err)
			}
			r, err := snippet.NewReminder(args[0], time.Now().Add(d), note)
			if err != nil {
				return err
			}
			if err := st.AddReminder(r); err != nil {
				return fmt.Errorf("saving reminder: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Reminder set for snippet %s due in %s\n", r.SnippetID, inDuration)
			return nil
		},
	}
	cmd.Flags().StringVarP(&inDuration, "in", "i", "24h", "Duration until reminder fires (e.g. 48h, 30m)")
	cmd.Flags().StringVarP(&note, "note", "n", "", "Optional note for the reminder")
	return cmd
}

func newReminderListCmd(st *store.Store) *cobra.Command {
	var dueOnly bool
	return &cobra.Command{
		Use:   "list",
		Short: "List reminders",
		RunE: func(cmd *cobra.Command, args []string) error {
			reminders, err := st.ListReminders()
			if err != nil {
				return err
			}
			if dueOnly {
				reminders = snippet.DueReminders(reminders, time.Now())
			}
			if len(reminders) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No reminders found.")
				return nil
			}
			for _, r := range reminders {
				note := ""
				if r.Note != "" {
					note = fmt.Sprintf(" — %s", r.Note)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "[%s] due %s%s\n", r.SnippetID, r.DueAt.Format(time.RFC822), note)
			}
			return nil
		},
	}
}

func newReminderRemoveCmd(st *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <snippet-id>",
		Short: "Remove a reminder for a snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := st.RemoveReminder(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Reminder for snippet %s removed.\n", args[0])
			return nil
		},
	}
}
