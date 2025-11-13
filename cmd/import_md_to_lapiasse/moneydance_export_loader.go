package main

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"os"

	"adeynack.net/lapiasse/cmd/import_md_to_lapiasse/moneydance"
)

type moneydanceExportLoader struct {
	FilePath string
}

func (loader *moneydanceExportLoader) Load(ctx context.Context) error {
	file, err := os.Open(loader.FilePath)
	if err != nil {
		return fmt.Errorf("opening Moneydance export file: %w", err)
	}
	defer file.Close()

	var mdExport moneydance.Export
	if err := json.UnmarshalRead(file, &mdExport); err != nil {
		return fmt.Errorf("decoding Moneydance export file: %w", err)
	}

	return nil
}
