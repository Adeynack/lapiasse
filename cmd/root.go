package cmd

import (
	"os"

	"adeynack.net/lapiasse/pkg/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configFilePath string       // Path to the configuration file
	configuration  *viper.Viper // Configuration of the application
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
	Run: executeTui,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
		configuration, err = app.InitializeConfiguration(cmd)

		return err
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	f := rootCmd.Flags()
	f.StringVar(&configFilePath, "config", "", "config file (default is $HOME/.lapiasse.yaml)")

	p := rootCmd.PersistentFlags()
	p.Bool("serve-web", false, "enable the web server, making the web API accessible via HTTP")
	p.String("data", "", "specify where the data for the application is located")
}
