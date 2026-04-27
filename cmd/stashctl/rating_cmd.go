package main

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/user/stashctl/internal/snippet"
	"github.com/user/stashctl/internal/store"
)

func newRatingCmd(st *store.Store) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rating",
		Short: "Manage snippet ratings",
	}
	cmd.AddCommand(
		newRatingAddCmd(st),
		newRatingListCmd(st),
		newRatingRemoveCmd(st),
	)
	return cmd
}

func newRatingAddCmd(st *store.Store) *cobra.Command {
	var note string
	cmd := &cobra.Command{
		Use:   "add <snippet-id> <score>",
		Short: "Rate a snippet (1–5)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			snippetID := args[0]
			score, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("score must be an integer: %w", err)
			}
			r, err := snippet.NewRating(snippetID, score, note)
			if err != nil {
				return err
			}
			if err := st.AddRating(r); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "rated snippet %s with score %d\n", snippetID, score)
			return nil
		},
	}
	cmd.Flags().StringVar(&note, "note", "", "optional note for the rating")
	return cmd
}

func newRatingListCmd(st *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "list <snippet-id>",
		Short: "List ratings for a snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			all, err := st.ListRatings()
			if err != nil {
				return err
			}
			ratings := snippet.RatingsFor(args[0], all)
			if len(ratings) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no ratings found")
				return nil
			}
			avg := snippet.AverageScore(ratings)
			fmt.Fprintf(cmd.OutOrStdout(), "ratings: %d  average: %.1f\n", len(ratings), avg)
			for _, r := range ratings {
				line := fmt.Sprintf("  [%s] score=%d", r.ID[:8], r.Score)
				if r.Note != "" {
					line += "  note: " + r.Note
				}
				fmt.Fprintln(cmd.OutOrStdout(), line)
			}
			return nil
		},
	}
}

func newRatingRemoveCmd(st *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <rating-id>",
		Short: "Remove a rating by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := st.DeleteRating(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "rating %s removed\n", args[0])
			return nil
		},
	}
}
