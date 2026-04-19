package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/stashctl/internal/format"
	"github.com/user/stashctl/internal/snippet"
)

func newLanguagesCmd() *cobra.Command {
	var check string

	cmd := &cobra.Command{
		Use:   "languages",
		Short: "List supported languages or check if a language is recognized",
		RunE: func(cmd *cobra.Command, args []string) error {
			if check != "" {
				norm := snippet.NormalizeLanguage(check)
				if snippet.IsKnownLanguage(norm) {
					fmt.Fprintf(cmd.OutOrStdout(), "%q is a recognized language.\n", norm)
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "%q is NOT a recognized language.\n", norm)
				}
				return nil
			}
			fmt.Fprint(cmd.OutOrStdout(), format.LanguageList())
			return nil
		},
	}

	cmd.Flags().StringVar(&check, "check", "", "check if a specific language is supported")
	return cmd
}
