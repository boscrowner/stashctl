package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/user/stashctl/internal/snippet"
)

func newDependencyCmd(app *appContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dep",
		Short: "Manage snippet dependencies",
	}
	cmd.AddCommand(
		newDepAddCmd(app),
		newDepListCmd(app),
		newDepRemoveCmd(app),
	)
	return cmd
}

func newDepAddCmd(app *appContext) *cobra.Command {
	var note string
	cmd := &cobra.Command{
		Use:   "add <source-id> <target-id>",
		Short: "Add a dependency from source snippet to target snippet",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			dep, err := snippet.NewDependency(args[0], args[1], note)
			if err != nil {
				return err
			}
			if err := app.store.AddDependency(dep); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "dependency added: %s -> %s (id: %s)\n", dep.SourceID, dep.TargetID, dep.ID)
			return nil
		},
	}
	cmd.Flags().StringVar(&note, "note", "", "Optional note describing the dependency")
	return cmd
}

func newDepListCmd(app *appContext) *cobra.Command {
	var targetOf bool
	cmd := &cobra.Command{
		Use:   "list <snippet-id>",
		Short: "List dependencies for a snippet",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := app.store.LoadDependencies()
			if err != nil {
				return err
			}
			var results []snippet.Dependency
			if targetOf {
				results = snippet.DependentsOf(args[0], deps)
			} else {
				results = snippet.DependenciesFor(args[0], deps)
			}
			if len(results) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no dependencies found")
				return nil
			}
			for _, d := range results {
				note := ""
				if d.Note != "" {
					note = fmt.Sprintf(" (%s)", d.Note)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%s -> %s%s [%s]\n", d.SourceID, d.TargetID, note, d.ID)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&targetOf, "dependents", false, "List snippets that depend on this snippet instead")
	return cmd
}

func newDepRemoveCmd(app *appContext) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <dependency-id>",
		Short: "Remove a dependency by its ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := app.store.DeleteDependency(args[0]); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "dependency %s removed\n", args[0])
			return nil
		},
	}
}
