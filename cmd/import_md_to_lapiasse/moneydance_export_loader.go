package main

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"os"
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

	var mdExport MdExport
	if err := json.UnmarshalRead(file, &mdExport); err != nil {
		return fmt.Errorf("decoding Moneydance export file: %w", err)
	}

	return nil
}
