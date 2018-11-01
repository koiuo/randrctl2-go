package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

// set by linker
var version string = "0.0.0"

func VersionCmd(ctx *Context) *cobra.Command {
	return &cobra.Command{
		Use: "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintln(ctx.Stdout, version)
		},
	}
}
