package cmd

import (
	"github.com/edio/randrctl2/lib"
	"github.com/edio/randrctl2/profile"
	"github.com/edio/randrctl2/x"
	"github.com/spf13/cobra"
	"io"
	"os"
)

func CatCmd(ctx *Context) *cobra.Command {
	var raw bool
	catCmd := cobra.Command{
		Use:   "cat [PROFILE]",
		Short: "Print profile",
		Long:  "Print profile with a given name if specified. Print current setup as profile if no argument given",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 || args[0] == "." || len(args[0]) == 0 {
				return catActive(ctx)
			} else {
				if raw {
					return catSaved(ctx, args[0], asRaw)
				} else {
					return catSaved(ctx, args[0], asParsed)
				}
			}
		},
	}
	catCmd.Flags().BoolVarP(&raw, "raw", "r", false, "do not parse profile, just print file contents as is")
	return &catCmd
}

func asParsed(writer io.Writer, reader io.Reader) error {
	parsedProfile, err := profile.Read(reader)
	if err != nil {
		return err
	}
	profile.Write(writer, parsedProfile)
	return nil
}

func asRaw(writer io.Writer, reader io.Reader) error {
	_, err := io.Copy(writer, reader)
	return err
}

func catSaved(ctx *Context, profileName string, writeAs func(writer io.Writer, reader io.Reader) error) error {
	for _, file := range lib.ListFiles(ctx.ProfilesDir) {
		if file.Name == profileName {
			profileFile, err := os.OpenFile(file.Path, os.O_RDONLY, 0)
			if err != nil {
				return err
			}
			return writeAs(ctx.Stdout, profileFile)
		}
	}
	return lib.SimpleErrorf("%s: no such profile", profileName)
}

func catActive(ctx *Context) error {
	x.Connect(ctx.Display)
	defer x.Disconnect()
	connected, err := x.GetConnectedOutputs()
	if err != nil {
		return err
	}
	_, primary, err := x.FindPrimary(connected)
	if err != nil {
		return err
	}

	pr := lib.ToProfile(connected, primary)
	profile.Write(ctx.Stdout, pr)
	return nil
}
