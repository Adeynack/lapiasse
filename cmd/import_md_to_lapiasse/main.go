package main

import (
	"fmt"
	"os"

	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

var (
	mdImporter moneydanceImporter
)

func main() {
	rootCmd := cobra.Command{
		Use:   "import_md_to_lapiasse",
		Short: "Import Moneydance exported data into lapiasse",
		Run: func(cmd *cobra.Command, args []string) {
			if err := mdImporter.Start(cmd.Context()); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		},
	}

	rootCmd.Flags().StringVar(&mdImporter.ApiEndpoint, "api", "http://localhost:8080", "Lapiasse API endpoint")
	rootCmd.Flags().StringVar(&mdImporter.ApiToken, "api-token", "", "Lapiasse API bearer token")
	lo.Must0(rootCmd.MarkFlagRequired("api-token"))
	rootCmd.Flags().StringVar(&mdImporter.MoneydanceExportPath, "md-source", "", "Path to Moneydance export file (JSON)")
	lo.Must0(rootCmd.MarkFlagRequired("md-source"))
	rootCmd.Flags().StringVar(&mdImporter.BookName, "book-name", "", "Name of the book to create in La Piasse. If empty, the name will be determined from the Moneydance export JSON (metadata->file_name).")
	rootCmd.Flags().BoolVar(&mdImporter.DeleteBookIfExists, "delete-book-if-exists", false, "Deletes the book if it already exists (by its name) in La Piasse before importing.")
	rootCmd.Flags().BoolVar(&mdImporter.SkipSuffixToBookName, "skip-suffix-to-book-name", false, `Whether to skip adding a suffix to the book name to ensure uniqueness (e.g.: "(imported on 2025-11-12 at 17:53:20)")`)
	rootCmd.Flags().StringVar(&mdImporter.BookCurrencyIsoCode, "cur", "EUR", "ISO code of the imported book's currency")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
