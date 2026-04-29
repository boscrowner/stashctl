package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/user/stashctl/internal/snippet"
)

func newWebhookCmd(app *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "Manage webhooks for snippet lifecycle events",
	}
	cmd.AddCommand(
		newWebhookAddCmd(app),
		newWebhookListCmd(app),
		newWebhookRemoveCmd(app),
	)
	return cmd
}

func newWebhookAddCmd(app *appState) *cobra.Command {
	var events []string
	var secret string
	cmd := &cobra.Command{
		Use:   "add <snippet-id> <url>",
		Short: "Register a webhook for a snippet",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			snippetID, rawURL := args[0], args[1]
			var evs []snippet.WebhookEvent
			for _, e := range events {
				evs = append(evs, snippet.WebhookEvent(strings.TrimSpace(e)))
			}
			w, err := snippet.NewWebhook(snippetID, rawURL, evs, secret)
			if err != nil {
				return err
			}
			if err := app.store.AddWebhook(w); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "webhook registered: %s\n", w.ID)
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&events, "events", []string{"created"}, "comma-separated events (created,updated,deleted)")
	cmd.Flags().StringVar(&secret, "secret", "", "optional signing secret")
	return cmd
}

func newWebhookListCmd(app *appState) *cobra.Command {
	return &cobra.Command{
		Use:   "list <snippet-id>",
		Short: "List webhooks for a snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hooks, err := app.store.ListWebhooks(args[0])
			if err != nil {
				return err
			}
			if len(hooks) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no webhooks registered")
				return nil
			}
			for _, w := range hooks {
				evStrs := make([]string, len(w.Events))
				for i, e := range w.Events {
					evStrs[i] = string(e)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t[%s]\n", w.ID, w.URL, strings.Join(evStrs, ","))
			}
			return nil
		},
	}
}

func newWebhookRemoveCmd(app *appState) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <webhook-id>",
		Short: "Remove a registered webhook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.store.DeleteWebhook(args[0]); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "webhook removed")
			return nil
		},
	}
}
