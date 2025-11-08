package main

import (
	"fmt"
	"os"

	"adeynack.net/lapiasse/pkg/app"
	"adeynack.net/lapiasse/pkg/env"
	"github.com/spf13/cobra"
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Output the detected runtime environment and the loaded configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := env.AutoRegisterEnvironments(cmd.Context())

		fmt.Println()
		fmt.Printf("Environment: %s\n", env.GetRunEnv(ctx).String())

		configuration, err := app.InitializeConfiguration(ctx, cliFlags)
		if err != nil {
			return fmt.Errorf("initializing configuration: %w", err)
		}

		fmt.Println()
		fmt.Println("Configuration file path:", configuration.Path)

		fmt.Println()
		fmt.Println("Effective configuration:")
		if err := configuration.WriteTo(os.Stdout); err != nil {
			return err
		}
		fmt.Print("\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
}
