package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/stashctl/internal/snippet"
)

func newReviewCmd(s storeIface) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "review",
		Short: "Manage snippet reviews",
	}
	cmd.AddCommand(
		newReviewAddCmd(s),
		newReviewListCmd(s),
		newReviewRemoveCmd(s),
	)
	return cmd
}

func newReviewAddCmd(s storeIface) *cobra.Command {
	var status string
	var comment string
	cmd := &cobra.Command{
		Use:   "add <snippet-id> <reviewer>",
		Short: "Add a review to a snippet",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			r, err := snippet.NewReview(args[0], args[1], comment, snippet.ReviewStatus(status))
			if err != nil {
				return err
			}
			if err := s.AddReview(r); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "review %s added\n", r.ID)
			return nil
		},
	}
	cmd.Flags().StringVar(&status, "status", "pending", "Review status: pending, approved, rejected")
	cmd.Flags().StringVar(&comment, "comment", "", "Review comment")
	return cmd
}

func newReviewListCmd(s storeIface) *cobra.Command {
	return &cobra.Command{
		Use:   "list <snippet-id>",
		Short: "List reviews for a snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			reviews, err := s.ReviewsForSnippet(args[0])
			if err != nil {
				return err
			}
			if len(reviews) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no reviews found")
				return nil
			}
			for _, r := range reviews {
				fmt.Fprintf(cmd.OutOrStdout(), "[%s] %s by %s: %s\n", r.Status, r.ID[:8], r.Reviewer, r.Comment)
			}
			return nil
		},
	}
}

func newReviewRemoveCmd(s storeIface) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <review-id>",
		Short: "Remove a review by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := s.DeleteReview(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "review %s removed\n", args[0])
			return nil
		},
	}
}
