package main

import (
	"log/slog"
	"os"
)

type moneydanceImporter struct {
	ApiEndpoint          string
	MoneydanceExportPath string
}

func (mdi *moneydanceImporter) Start() error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info(
		"Importing MoneyDance Export JSON to La Piasse",
		slog.String("api-endpoint", mdi.ApiEndpoint),
		slog.String("md-source", mdi.MoneydanceExportPath),
	)

	return nil
}
