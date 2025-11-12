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
			return mdImporter.Start(cmd.Context())
		},
	}

	rootCmd.Flags().StringVar(&mdImporter.ApiEndpoint, "api", "http://localhost:8080", "Lapiasse API endpoint")
	rootCmd.Flags().StringVar(&mdImporter.MoneydanceExportPath, "md-source", "", "Path to Moneydance export file (JSON)")
	rootCmd.Flags().StringVar(&mdImporter.BookName, "book-name", "", "Name of the book to create in La Piasse. If empty, the name will be determined from the Moneydance export JSON (metadata->file_name).")
	rootCmd.Flags().BoolVar(&mdImporter.AddSuffixToBookName, "add-suffix", true, `Whether to add a suffix to the book name to ensure uniqueness (e.g.: "(imported on 2025-11-12 at 17:53:20)")`)
	rootCmd.Flags().StringVar(&mdImporter.BookCurrencyIsoCode, "cur", "EUR", "ISO code of the imported book's currency")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
