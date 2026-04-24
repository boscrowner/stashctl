package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newLabelCmd(app *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "label",
		Short: "Manage colour labels",
	}
	cmd.AddCommand(
		newLabelCreateCmd(app),
		newLabelListCmd(app),
		newLabelRemoveCmd(app),
	)
	return cmd
}

func newLabelCreateCmd(app *appState) *cobra.Command {
	var color string
	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "Create a new label",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			label, err := snippetpkg.NewLabel(args[0], color)
			if err != nil {
				return err
			}
			app.store.Labels = append(app.store.Labels, label)
			if err := app.store.Save(); err != nil {
				return fmt.Errorf("save: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "created label %q (%s)\n", label.Name, label.ID)
			return nil
		},
	}
	cmd.Flags().StringVar(&color, "color", "", "hex colour code (default #cccccc)")
	return cmd
}

func newLabelListCmd(app *appState) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all labels",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if len(app.store.Labels) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no labels defined")
				return nil
			}
			for _, l := range app.store.Labels {
				fmt.Fprintf(cmd.OutOrStdout(), "%s\t%s\t%s\n", l.ID[:8], l.Color, l.Name)
			}
			return nil
		},
	}
}

func newLabelRemoveCmd(app *appState) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a label by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			l, ok := snippetpkg.FindLabel(app.store.Labels, args[0])
			if !ok {
				return fmt.Errorf("label %q not found", args[0])
			}
			app.store.Labels = snippetpkg.RemoveLabel(app.store.Labels, l.ID)
			if err := app.store.Save(); err != nil {
				return fmt.Errorf("save: %w", err)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "removed label %q\n", l.Name)
			return nil
		},
	}
}
