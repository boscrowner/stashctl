package main

import (
	"fmt"

	"github.com/spf13/cobra"
	importer "github.com/user/stashctl/internal/import"
	"github.com/user/stashctl/internal/store"
)

func newImportCmd(st *store.Store) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "import <file>",
		Short: "Import snippets from a JSON file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := args[0]
			snippets, res, err := importer.FromFile(path)
			if err != nil {
				return fmt.Errorf("import failed: %w", err)
			}
			for _, s := range snippets {
				if err := st.Add(s); err != nil {
					res.Skipped++
					res.Errors = append(res.Errors, fmt.Sprintf("store error for %q: %v", s.Title, err))
					continue
				}
				res.Imported++
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Imported: %d, Skipped: %d\n", res.Imported, res.Skipped)
			for _, e := range res.Errors {
				fmt.Fprintf(cmd.ErrOrStderr(), "  warning: %s\n", e)
			}
			return nil
		},
	}
	return cmd
}
