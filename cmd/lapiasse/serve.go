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

	// running this command implies:
	configuration.Configuration.Web.Expose = true // exposing the web interface
	configuration.Configuration.Log.UILess = true // running in UI-less mode

	appInstance, err := app.NewInstance(configuration)
	if err != nil {
		return fmt.Errorf("initializing application instance: %w", err)
	}
	defer appInstance.Close()

	// Listen for interrupt signal (eg: Ctrl-C) to gracefully shutdown the server (managed by appInstance).
	if configuration.Configuration.DryStart {
		applog.Info(appInstance.Context(), "Dry start enabled, exiting now.")
		return nil
	}

	applog.Info(appInstance.Context(), "Server is running. Press Ctrl-C to stop.")
	<-appInstance.Context().Done()

	return nil
}
