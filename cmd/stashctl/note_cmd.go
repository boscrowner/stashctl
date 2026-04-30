package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/user/stashctl/internal/snippet"
	"github.com/user/stashctl/internal/store"
)

func newNoteCmd(st *store.Store) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "note",
		Short: "Manage notes attached to snippets",
	}
	cmd.AddCommand(
		newNoteAddCmd(st),
		newNoteListCmd(st),
		newNoteRemoveCmd(st),
	)
	return cmd
}

func newNoteAddCmd(st *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "add <snippet-id> <body>",
		Short: "Add a note to a snippet",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			snippetID, body := args[0], args[1]
			if _, err := st.Get(snippetID); err != nil {
				return fmt.Errorf("snippet %q not found", snippetID)
			}
			n, err := snippet.NewNote(snippetID, body)
			if err != nil {
				return err
			}
			if err := st.AddNote(n); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "note %s added\n", n.ID)
			return nil
		},
	}
}

func newNoteListCmd(st *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "list <snippet-id>",
		Short: "List notes for a snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			notes, err := st.ListNotes(args[0])
			if err != nil {
				return err
			}
			if len(notes) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no notes")
				return nil
			}
			for _, n := range notes {
				fmt.Fprintf(cmd.OutOrStdout(), "[%s] %s\n", n.ID[:8], n.Body)
			}
			return nil
		},
	}
}

func newNoteRemoveCmd(st *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <note-id>",
		Short: "Remove a note by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := st.RemoveNote(args[0]); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "note removed")
			return nil
		},
	}
}
