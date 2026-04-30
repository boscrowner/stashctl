package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/stashctl/internal/snippet"
	"github.com/user/stashctl/internal/store"
)

func newCommentCmd(st *store.Store) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "comment",
		Short: "Manage comments on snippets",
	}
	cmd.AddCommand(
		newCommentAddCmd(st),
		newCommentListCmd(st),
		newCommentRemoveCmd(st),
	)
	return cmd
}

func newCommentAddCmd(st *store.Store) *cobra.Command {
	var author string
	cmd := &cobra.Command{
		Use:   "add <snippet-id> <body>",
		Short: "Add a comment to a snippet",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := snippet.NewComment(args[0], author, args[1])
			if err != nil {
				return err
			}
			if err := st.AddComment(c); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "comment %s added\n", c.ID)
			return nil
		},
	}
	cmd.Flags().StringVarP(&author, "author", "a", "anonymous", "Comment author")
	return cmd
}

func newCommentListCmd(st *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "list <snippet-id>",
		Short: "List comments for a snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			comments, err := st.ListComments(args[0])
			if err != nil {
				return err
			}
			if len(comments) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no comments")
				return nil
			}
			for _, c := range comments {
				fmt.Fprintf(cmd.OutOrStdout(), "[%s] %s: %s\n", c.ID[:8], c.Author, c.Body)
			}
			return nil
		},
	}
}

func newCommentRemoveCmd(st *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <comment-id>",
		Short: "Remove a comment by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := st.RemoveComment(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "comment %s removed\n", args[0])
			return nil
		},
	}
}
