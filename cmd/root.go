package cmd

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
	cliFlags.Config = p.String("config", "", "config file (default is $HOME/.lapiasse.yaml)")
	cliFlags.ServeWeb = p.Bool("serve-web", false, "enable the web server, making the web API accessible via HTTP")
	cliFlags.Data = p.String("data", "", "specify where the data for the application is located")
}
