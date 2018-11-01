package cmd

import (
	"github.com/edio/randrctl2/lib"
	"github.com/edio/randrctl2/x"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"
)

type Context struct {
	Display     string
	ProfilesDir string
	Stdout      io.Writer
}

func RootCmd(vpr *viper.Viper, ctx *Context) *cobra.Command {
	rootCmd := &cobra.Command{
		Use: "randrctl",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Usage()
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctx.Stdout = cmd.OutOrStdout()
			ctx.Display = ""
			return nil
		},
		SilenceUsage:  true,
		SilenceErrors: false,
	}
	return rootCmd
}

func Execute() {
	home, _ := homedir.Dir()
	configDir := filepath.Join(home, ".config", "randrctl2")
	profilesDir := filepath.Join(configDir, "profiles")
	os.MkdirAll(profilesDir, 0755)

	vpr := viper.New()
	vpr.SetConfigType("yaml")
	vpr.SetConfigName("config")
	vpr.AddConfigPath(configDir)

	log.SetLevel(log.WarnLevel)

	ctx := &Context{
		ProfilesDir: filepath.Join(configDir, "profiles"),
	}
	rootCmd := RootCmd(vpr, ctx)
	rootCmd.AddCommand(CatCmd(ctx))
	rootCmd.AddCommand(ListCmd(ctx))
	rootCmd.AddCommand(VersionCmd(ctx))

	if err := rootCmd.Execute(); err != nil {
		switch err.(type) {
		case lib.SimpleError:
			os.Exit(2)
		case *x.XError:
			os.Exit(64)
		default:
			rootCmd.Usage()
			os.Exit(1)
		}
	}
}
