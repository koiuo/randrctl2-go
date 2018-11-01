package cmd

import (
	"fmt"
	"github.com/edio/randrctl2/lib"
	"github.com/spf13/cobra"
)

func ListCmd(ctx *Context) *cobra.Command {
	listCmd := cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List profiles",
		Long:    "Print profile with a given name if specified. Print current setup as profile if no argument given",
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(ctx)
		},
	}
	return &listCmd
}

func list(ctx *Context) error {
	for _, file := range lib.ListFiles(ctx.ProfilesDir) {
		fmt.Fprintln(ctx.Stdout, file.Name)
	}
	return nil
}
