package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/user/stashctl/internal/snippet"
)

func newFavoriteCmd(env *cmdEnv) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "favorite",
		Short: "Manage favorite snippets",
	}
	cmd.AddCommand(newFavoriteAddCmd(env), newFavoriteRemoveCmd(env), newFavoriteListCmd(env))
	return cmd
}

func newFavoriteAddCmd(env *cmdEnv) *cobra.Command {
	return &cobra.Command{
		Use:   "add <id>",
		Short: "Mark a snippet as favorite",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := env.store.Get(args[0])
			if err != nil {
				return err
			}
			if s.Favorite {
				fmt.Fprintf(cmd.OutOrStdout(), "%s is already a favorite\n", s.ID)
				return nil
			}
			snippet.Favorite(s)
			if err := env.store.Update(s); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "marked %s as favorite\n", s.ID)
			return nil
		},
	}
}

func newFavoriteRemoveCmd(env *cmdEnv) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <id>",
		Short: "Unmark a snippet as favorite",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s, err := env.store.Get(args[0])
			if err != nil {
				return err
			}
			if !s.Favorite {
				fmt.Fprintf(cmd.OutOrStdout(), "%s is not a favorite\n", s.ID)
				return nil
			}
			snippet.Unfavorite(s)
			if err := env.store.Update(s); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "removed %s from favorites\n", s.ID)
			return nil
		},
	}
}

func newFavoriteListCmd(env *cmdEnv) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List favorite snippets",
		RunE: func(cmd *cobra.Command, args []string) error {
			all, err := env.store.List()
			if err != nil {
				return err
			}
			favs := snippet.Favorites(all)
			if len(favs) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no favorite snippets found")
				return nil
			}
			for _, s := range favs {
				fmt.Fprintln(cmd.OutOrStdout(), env.formatter.SnippetSummary(s))
			}
			return nil
		},
	}
}
