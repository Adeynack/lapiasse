package main

import (
	"os"

	"adeynack.net/lapiasse/pkg/app"
	"github.com/spf13/cobra"
)

var (
	cliFlags app.CliFlags
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lapiasse",
	Short: "Personal finances manager",
	Long: `
This is a multi-mode application. It offers a TUI for interacting with
its model. It also offers an API server, for integration and scripting
against its business logic.

When no command is passed, the TUI is launched automatically.
	`,
	RunE: executeTui,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	p := rootCmd.PersistentFlags()
	p.StringVar(&cliFlags.Config, "config", "", "Optional config file path. Will default to the proper user configuration directory of the running platform.")
	p.BoolVar(&cliFlags.ServeWeb, "serve-web", false, "Enable the web server, making the web API accessible via HTTP.")
	p.StringVar(&cliFlags.Data, "data", "", "Optional path to the data directory. Will default to the proper user data directory of the running platform.")
	p.BoolVar(&cliFlags.DryStart, "dry-start", false, "Enable dry start mode. Will initialize the specified command, but exit then immediately.")
}
