package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/user/stashctl/internal/snippet"
)

func newTagGroupCmd(app *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tag-group",
		Short: "Manage named tag groups",
	}
	cmd.AddCommand(
		newTagGroupCreateCmd(app),
		newTagGroupListCmd(app),
		newTagGroupRemoveCmd(app),
	)
	return cmd
}

func newTagGroupCreateCmd(app *appState) *cobra.Command {
	var description string
	cmd := &cobra.Command{
		Use:   "create <name> <tag,...>",
		Short: "Create a new tag group",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			tags := snippet.ParseTags(args[1])
			g, err := snippet.NewTagGroup(args[0], description, tags)
			if err != nil {
				return err
			}
			app.tagGroups = append(app.tagGroups, g)
			fmt.Fprintf(cmd.OutOrStdout(), "created tag group %q (%s)\n", g.Name, g.ID)
			return nil
		},
	}
	cmd.Flags().StringVarP(&description, "description", "d", "", "optional description")
	return cmd
}

func newTagGroupListCmd(app *appState) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all tag groups",
		RunE: func(cmd *cobra.Command, _ []string) error {
			groups := snippet.TagGroupsFor(app.tagGroups)
			if len(groups) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "no tag groups defined")
				return nil
			}
			w := cmd.OutOrStdout()
			for _, g := range groups {
				fmt.Fprintf(w, "%-20s  [%s]\n", g.Name, strings.Join(g.Tags, ", "))
			}
			return nil
		},
	}
}

func newTagGroupRemoveCmd(app *appState) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a tag group by name",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			g, ok := snippet.FindTagGroup(app.tagGroups, args[0])
			if !ok {
				return fmt.Errorf("tag group %q not found", args[0])
			}
			app.tagGroups = snippet.RemoveTagGroup(app.tagGroups, g.ID)
			fmt.Fprintf(cmd.OutOrStdout(), "removed tag group %q\n", g.Name)
			return nil
		},
	}
}
