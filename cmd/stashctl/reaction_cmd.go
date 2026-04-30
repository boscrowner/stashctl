package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/stashctl/internal/snippet"
)

func newReactionCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reaction",
		Short: "Manage emoji reactions on snippets",
	}
	cmd.AddCommand(
		newReactionAddCmd(app),
		newReactionListCmd(app),
		newReactionRemoveCmd(app),
	)
	return cmd
}

func newReactionAddCmd(app *appContext) *cobra.Command {
	return &cobra.Command{
		Use:   "add <snippet-id> <emoji>",
		Short: "Add a reaction to a snippet",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			snippetID, emoji := args[0], args[1]
			actor := app.cfg.DefaultActor
			if actor == "" {
				actor = "user"
			}
			r, err := snippet.NewReaction(snippetID, actor, emoji)
			if err != nil {
				return err
			}
			if err := app.store.AddReaction(r); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "reaction %s added to snippet %s\n", emoji, snippetID)
			return nil
		},
	}
}

func newReactionListCmd(app *appContext) *cobra.Command {
	return &cobra.Command{
		Use:   "list <snippet-id>",
		Short: "List reactions for a snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reactions, err := app.store.ListReactions(args[0])
			if err != nil {
				return err
			}
			if len(reactions) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no reactions")
				return nil
			}
			counts := snippet.CountByEmoji(reactions)
			for emoji, count := range counts {
				fmt.Fprintf(cmd.OutOrStdout(), ":%s: %s\n", emoji, strings.Repeat("|", count))
			}
			return nil
		},
	}
}

func newReactionRemoveCmd(app *appContext) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <reaction-id>",
		Short: "Remove a reaction by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.store.DeleteReaction(args[0]); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "reaction removed")
			return nil
		},
	}
}
