package main

import (
	"fmt"

	"adeynack.net/lapiasse/pkg/app"
	"adeynack.net/lapiasse/pkg/env"
	"adeynack.net/lapiasse/pkg/platform/loex"
	"adeynack.net/lapiasse/pkg/tui"
	"github.com/spf13/cobra"
)

// tuiCmd represents the tui command
var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Launch the Terminal User Interface (TUI)",
	Long: `
For the terminal lovers, this application offers a text based version
of its user interface.
	`,
	RunE: executeTui,
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

func executeTui(cmd *cobra.Command, args []string) (err error) {
	ctx := env.AutoRegisterEnvironments(cmd.Context())

	configuration, err := app.InitializeConfiguration(ctx, cliFlags)
	if err != nil {
		return fmt.Errorf("initializing configuration: %w", err)
	}

	appInstance, err := app.NewInstance(ctx, configuration)
	if err != nil {
		return fmt.Errorf("initializing application instance: %w", err)
	}
	defer loex.OnErrJoin(&err, appInstance.Close)

	tuiInstance, err := tui.NewInstance(appInstance)
	if err != nil {
		return fmt.Errorf("initializing TUI instance: %w", err)
	}
	defer loex.OnErrJoin(&err, tuiInstance.Close)

	err = tuiInstance.Run()
	if err != nil {
		return fmt.Errorf("running TUI instance: %w", err)
	}

	return nil
}
