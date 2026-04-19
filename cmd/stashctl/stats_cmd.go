package main

import (
	"fmt"
	"sort"

	"github.com/spf13/cobra"
	"github.com/user/stashctl/internal/snippet"
	"github.com/user/stashctl/internal/store"
)

func newStatsCmd(st *store.Store) *cobra.Command {
	return &cobra.Command{
		Use:   "stats",
		Short: "Show statistics about your snippet stash",
		RunE: func(cmd *cobra.Command, args []string) error {
			snippets, err := st.List()
			if err != nil {
				return err
			}
			s := snippet.ComputeStats(snippets)
			fmt.Fprintf(cmd.OutOrStdout(), "Total snippets: %d\n", s.Total)
			if s.Total == 0 {
				return nil
			}
			fmt.Fprintln(cmd.OutOrStdout(), "\nBy language:")
			langs := sortedKeys(s.ByLanguage)
			for _, l := range langs {
				fmt.Fprintf(cmd.OutOrStdout(), "  %-16s %d\n", l, s.ByLanguage[l])
			}
			fmt.Fprintln(cmd.OutOrStdout(), "\nBy tag:")
			tags := sortedKeys(s.ByTag)
			for _, t := range tags {
				fmt.Fprintf(cmd.OutOrStdout(), "  %-16s %d\n", t, s.ByTag[t])
			}
			fmt.Fprintf(cmd.OutOrStdout(), "\nOldest created: %s\n", s.OldestCreated.Format("2006-01-02"))
			fmt.Fprintf(cmd.OutOrStdout(), "Last updated:   %s\n", s.NewestUpdated.Format("2006-01-02"))
			return nil
		},
	}
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
