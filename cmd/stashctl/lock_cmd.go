package main

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/user/stashctl/internal/snippet"
)

func newLockCmd(env *cmdEnv) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock",
		Short: "Manage snippet edit locks",
	}
	cmd.AddCommand(
		newLockAcquireCmd(env),
		newLockReleaseCmd(env),
		newLockStatusCmd(env),
	)
	return cmd
}

func newLockAcquireCmd(env *cmdEnv) *cobra.Command {
	var owner string
	var ttl time.Duration

	cmd := &cobra.Command{
		Use:   "acquire <snippet-id>",
		Short: "Acquire an exclusive edit lock on a snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			snippetID := args[0]
			existing, ok := snippet.FindLock(env.store.Locks(), snippetID)
			if ok {
				return fmt.Errorf("%w: held by %s until %s",
					snippet.ErrLockConflict, existing.Owner,
					existing.ExpiresAt.Format(time.RFC3339))
			}
			l, err := snippet.NewLock(snippetID, owner, ttl)
			if err != nil {
				return err
			}
			if err := env.store.AddLock(l); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "lock acquired: %s (expires %s)\n",
				l.ID, l.ExpiresAt.Format(time.RFC3339))
			return nil
		},
	}
	cmd.Flags().StringVar(&owner, "owner", "", "lock owner identifier (required)")
	cmd.Flags().DurationVar(&ttl, "ttl", 30*time.Minute, "lock time-to-live")
	_ = cmd.MarkFlagRequired("owner")
	return cmd
}

func newLockReleaseCmd(env *cmdEnv) *cobra.Command {
	return &cobra.Command{
		Use:   "release <lock-id>",
		Short: "Release a previously acquired lock",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := env.store.RemoveLock(args[0]); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "lock released")
			return nil
		},
	}
}

func newLockStatusCmd(env *cmdEnv) *cobra.Command {
	return &cobra.Command{
		Use:   "status <snippet-id>",
		Short: "Show the active lock for a snippet, if any",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			l, ok := snippet.FindLock(env.store.Locks(), args[0])
			if !ok {
				fmt.Fprintln(cmd.OutOrStdout(), "no active lock")
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "locked by %s until %s (id: %s)\n",
				l.Owner, l.ExpiresAt.Format(time.RFC3339), l.ID)
			return nil
		},
	}
}
