/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"adeynack.net/lapiasse/pkg/app"
	"github.com/spf13/cobra"
)

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Output the detected runtime environment and the loaded configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Environment: %s\n", app.RunEnv)
		fmt.Println("Effective configuration:")
		if err := configuration.WriteConfigTo(os.Stdout); err != nil {
			return err
		}
		fmt.Print("\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(envCmd)
}
