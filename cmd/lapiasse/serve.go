package main

import (
	"fmt"

	"adeynack.net/lapiasse/pkg/app"
	"adeynack.net/lapiasse/pkg/applog"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the API and Web interface over HTTP",
	RunE:  executeServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func executeServe(cmd *cobra.Command, args []string) error {
	configuration, err := app.InitializeConfiguration(cliFlags)
	if err != nil {
		return fmt.Errorf("initializing configuration: %w", err)
	}

	configuration.Configuration.Web.Expose = true
	configuration.Configuration.Log.UILess = true

	appInstance, err := app.NewInstance(configuration)
	if err != nil {
		return fmt.Errorf("initializing application instance: %w", err)
	}
	defer appInstance.Close()

	// Listen for interrupt signal (eg: Ctrl-C) to gracefully shutdown the server (managed by appInstance).
	applog.Info(appInstance.Context(), "Server is running. Press Ctrl-C to stop.")
	<-appInstance.Context().Done()

	return nil
}
