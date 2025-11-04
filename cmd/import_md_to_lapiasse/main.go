package main

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	mdImporter moneydanceImporter
)

func main() {
	rootCmd := cobra.Command{
		Use:   "import_md_to_lapiasse",
		Short: "Import Moneydance exported data into lapiasse",
		RunE: func(cmd *cobra.Command, args []string) error {
			return mdImporter.Start()
		},
	}

	rootCmd.Flags().StringVar(&mdImporter.ApiEndpoint, "api", "http://localhost:8080", "Lapiasse API endpoint")
	rootCmd.Flags().StringVar(&mdImporter.MoneydanceExportPath, "md-source", "", "Path to Moneydance export file (JSON)")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
