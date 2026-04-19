package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/user/stashctl/internal/format"
	"github.com/user/stashctl/internal/snippet"
	"github.com/user/stashctl/internal/store"
)

func newRecentCmd(st *store.Store, color bool) *cobra.Command {
	var limit int
	var sinceDays string

	cmd := &cobra.Command{
		Use:   "recent",
		Short: "Show recently updated snippets",
		RunE: func(cmd *cobra.Command, args []string) error {
			snippets, err := st.List()
			if err != nil {
				return fmt.Errorf("listing snippets: %w", err)
			}

			opts := snippet.RecentOptions{Limit: limit}
			if sinceDays != "" {
				d, err := strconv.Atoi(sinceDays)
				if err != nil {
					return fmt.Errorf("invalid --since value: %w", err)
				}
				opts.Since = time.Now().AddDate(0, 0, -d)
			}

			results := snippet.Recent(snippets, opts)
			if len(results) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no recent snippets found")
				return nil
			}
			for _, s := range results {
				fmt.Fprintln(cmd.OutOrStdout(), format.SnippetSummary(s, color))
			}
			return nil
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "n", 10, "maximum number of snippets to show")
	cmd.Flags().StringVar(&sinceDays, "since", "", "only show snippets updated within N days")
	return cmd
}
