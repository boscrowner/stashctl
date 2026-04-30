package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/user/stashctl/internal/snippet"
	"github.com/user/stashctl/internal/store"
)

func newWorkflowCmd(st *store.Store) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workflow",
		Short: "Manage snippet workflow states",
	}
	cmd.AddCommand(
		newWorkflowTransitionCmd(st),
		newWorkflowLogCmd(st),
		newWorkflowStatusCmd(st),
	)
	return cmd
}

func newWorkflowTransitionCmd(st *store.Store) *cobra.Command {
	var actor, note string
	cmd := &cobra.Command{
		Use:   "transition <snippet-id> <from> <to>",
		Short: "Record a workflow state transition for a snippet",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			if actor == "" {
				actor = os.Getenv("USER")
			}
			tr, err := snippet.NewWorkflowTransition(
				args[0],
				snippet.WorkflowState(args[1]),
				snippet.WorkflowState(args[2]),
				actor, note,
			)
			if err != nil {
				return err
			}
			if err := st.AddWorkflowTransition(tr); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "transitioned %s: %s → %s\n", args[0], tr.From, tr.To)
			return nil
		},
	}
	cmd.Flags().StringVar(&actor, "actor", "", "actor performing the transition (defaults to $USER)")
	cmd.Flags().StringVar(&note, "note", "", "optional note for the transition")
	return cmd
}

func newWorkflowLogCmd(st *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "log <snippet-id>",
		Short: "Show workflow transition history for a snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			all, err := st.ListWorkflowTransitions()
			if err != nil {
				return err
			}
			transitions := snippet.TransitionsFor(args[0], all)
			if len(transitions) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no transitions recorded")
				return nil
			}
			w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "TIME\tFROM\tTO\tACTOR\tNOTE")
			for _, tr := range transitions {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					tr.At.Format("2006-01-02 15:04"), tr.From, tr.To, tr.Actor, tr.Note)
			}
			return w.Flush()
		},
	}
}

func newWorkflowStatusCmd(st *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "status <snippet-id>",
		Short: "Show current workflow state of a snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			all, err := st.ListWorkflowTransitions()
			if err != nil {
				return err
			}
			state := snippet.CurrentState(args[0], all)
			fmt.Fprintf(cmd.OutOrStdout(), "%s: %s\n", args[0], state)
			return nil
		},
	}
}
