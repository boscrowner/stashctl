package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/user/stashctl/internal/snippet"
)

func newBadgeCmd(env *cmdEnv) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "badge",
		Short: "Manage snippet badges",
	}
	cmd.AddCommand(
		newBadgeGrantCmd(env),
		newBadgeListCmd(env),
		newBadgeRemoveCmd(env),
	)
	return cmd
}

func newBadgeGrantCmd(env *cmdEnv) *cobra.Command {
	var icon, color, note, grantedBy string
	cmd := &cobra.Command{
		Use:   "grant <snippet-id> <name>",
		Short: "Grant a badge to a snippet",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			snippetID, name := args[0], args[1]
			if grantedBy == "" {
				grantedBy = env.cfg.Author
			}
			b, err := snippet.NewBadge(snippetID, name, icon, color, grantedBy, note)
			if err != nil {
				return err
			}
			if err := env.store.AddBadge(b); err != nil {
				return fmt.Errorf("grant badge: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "badge %s granted to snippet %s\n", b.ID, snippetID)
			return nil
		},
	}
	cmd.Flags().StringVar(&icon, "icon", "⭐", "badge icon (emoji or text)")
	cmd.Flags().StringVar(&color, "color", "", "badge color (red, green, blue, yellow, purple, orange, gray, gold)")
	cmd.Flags().StringVar(&note, "note", "", "optional note")
	cmd.Flags().StringVar(&grantedBy, "by", "", "granter identity (defaults to config author)")
	return cmd
}

func newBadgeListCmd(env *cmdEnv) *cobra.Command {
	return &cobra.Command{
		Use:   "list <snippet-id>",
		Short: "List badges for a snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			all, err := env.store.ListBadges()
			if err != nil {
				return fmt.Errorf("list badges: %w", err)
			}
			badges := snippet.BadgesFor(args[0], all)
			if len(badges) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no badges")
				return nil
			}
			for _, b := range badges {
				fmt.Fprintf(cmd.OutOrStdout(), "%s  %s %s  [%s] by %s\n",
					b.ID, b.Icon, b.Name, b.Color, b.GrantedBy)
			}
			return nil
		},
	}
}

func newBadgeRemoveCmd(env *cmdEnv) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <badge-id>",
		Short: "Remove a badge by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := env.store.RemoveBadge(args[0]); err != nil {
				return fmt.Errorf("remove badge: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "badge %s removed\n", args[0])
			return nil
		},
	}
}
