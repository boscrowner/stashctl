package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/stashctl/internal/snippet"
)

func newSubscriptionCmd(env *cmdEnv) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscription",
		Short: "Manage snippet subscriptions",
	}
	cmd.AddCommand(
		newSubscriptionAddCmd(env),
		newSubscriptionListCmd(env),
		newSubscriptionRemoveCmd(env),
	)
	return cmd
}

func newSubscriptionAddCmd(env *cmdEnv) *cobra.Command {
	var note string
	cmd := &cobra.Command{
		Use:   "add <snippet-id> <subscriber> <event>",
		Short: "Subscribe to events on a snippet",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			sub, err := snippet.NewSubscription(args[0], args[1], snippet.SubscriptionEvent(args[2]), note)
			if err != nil {
				return err
			}
			if err := env.store.AddSubscription(sub); err != nil {
				return err
			}
			fmt.Fprintf(env.out, "subscription added: %s\n", sub.ID)
			return nil
		},
	}
	cmd.Flags().StringVar(&note, "note", "", "optional note")
	return cmd
}

func newSubscriptionListCmd(env *cmdEnv) *cobra.Command {
	return &cobra.Command{
		Use:   "list <snippet-id>",
		Short: "List subscriptions for a snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			subs, err := env.store.ListSubscriptions(args[0])
			if err != nil {
				return err
			}
			if len(subs) == 0 {
				fmt.Fprintln(env.out, "no subscriptions")
				return nil
			}
			for _, s := range subs {
				fmt.Fprintf(env.out, "%s  %s  %s  %s\n", s.ID, s.Subscriber, s.Event, s.Note)
			}
			return nil
		},
	}
}

func newSubscriptionRemoveCmd(env *cmdEnv) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <id>",
		Short: "Remove a subscription by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := env.store.DeleteSubscription(args[0]); err != nil {
				return err
			}
			fmt.Fprintln(env.out, "subscription removed")
			return nil
		},
	}
}
