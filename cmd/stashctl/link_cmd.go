package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/stashctl/internal/snippet"
)

func newLinkCmd(env *cmdEnv) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "link",
		Short: "Manage URL links attached to snippets",
	}
	cmd.AddCommand(
		newLinkAddCmd(env),
		newLinkListCmd(env),
		newLinkRemoveCmd(env),
	)
	return cmd
}

func newLinkAddCmd(env *cmdEnv) *cobra.Command {
	var title string
	cmd := &cobra.Command{
		Use:   "add <snippet-id> <url>",
		Short: "Attach a URL link to a snippet",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			snippetID, rawURL := args[0], args[1]
			if _, err := env.store.Get(snippetID); err != nil {
				return fmt.Errorf("snippet %q not found", snippetID)
			}
			l, err := snippet.NewLink(snippetID, rawURL, title)
			if err != nil {
				return err
			}
			if err := env.store.AddLink(l); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "link %s added\n", l.ID)
			return nil
		},
	}
	cmd.Flags().StringVar(&title, "title", "", "Optional human-readable title for the link")
	return cmd
}

func newLinkListCmd(env *cmdEnv) *cobra.Command {
	return &cobra.Command{
		Use:   "list <snippet-id>",
		Short: "List URL links for a snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			links, err := env.store.ListLinks(args[0])
			if err != nil {
				return err
			}
			if len(links) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no links found")
				return nil
			}
			for _, l := range links {
				title := l.Title
				if title == "" {
					title = "(no title)"
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s  %s  %s\n", l.ID, l.URL, title)
			}
			return nil
		},
	}
}

func newLinkRemoveCmd(env *cmdEnv) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <link-id>",
		Short: "Remove a URL link by its ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := env.store.RemoveLink(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "link %s removed\n", args[0])
			return nil
		},
	}
}
