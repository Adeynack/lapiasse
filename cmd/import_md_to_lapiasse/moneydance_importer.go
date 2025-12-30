package main

import (
	"context"
	"encoding/json/v2"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"adeynack.net/lapiasse/cmd/import_md_to_lapiasse/moneydance"
	"adeynack.net/lapiasse/pkg/api"
	"adeynack.net/lapiasse/pkg/applog"
	"adeynack.net/lapiasse/pkg/platform/errorex"
)

type moneydanceImporter struct {
	ApiEndpoint          string
	MoneydanceExportPath string
	BookName             string
	AddSuffixToBookName  bool
	BookCurrencyIsoCode  string

	apiClient *api.ClientWithResponses
	book      *api.BookShow
	md        *moneydance.Export
}

func (mdi *moneydanceImporter) Start(ctx context.Context) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	ctx = applog.WithLogger(ctx, logger)

	applog.Info(
		ctx,
		"Importing MoneyDance Export JSON to La Piasse",
		slog.String("api-endpoint", mdi.ApiEndpoint),
		slog.String("md-source", mdi.MoneydanceExportPath),
		slog.String("book-name", mdi.BookName),
		slog.Bool("add-suffix", mdi.AddSuffixToBookName),
	)

	for _, step := range []func(context.Context) error{
		mdi.loadMoneydanceExport,
		mdi.createApiClient,
		mdi.createNewBook,
		mdi.createRegisters,
	} {
		var err error
		if err = step(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (mdi *moneydanceImporter) loadMoneydanceExport(ctx context.Context) (err error) {
	file, err := os.Open(mdi.MoneydanceExportPath)
	if err != nil {
		return fmt.Errorf("opening Moneydance export file: %w", err)
	}
	defer errorex.CallJoinErr(&err, file.Close)

	var mdExport moneydance.Export
	if err := json.UnmarshalRead(file, &mdExport); err != nil {
		return fmt.Errorf("decoding Moneydance export file: %w", err)
	}

	mdi.md = &mdExport

	return nil
}

func (mdi *moneydanceImporter) createApiClient(ctx context.Context) error {
	apiClient, err := api.NewClientWithResponses(mdi.ApiEndpoint)
	if err != nil {
		return fmt.Errorf("creating the API client: %w", err)
	}
	mdi.apiClient = apiClient

	return nil
}

func (mdi *moneydanceImporter) createNewBook(ctx context.Context) error {
	bookName := mdi.determineBookName()

	response, err := mdi.apiClient.CreateBookWithResponse(ctx, api.CreateBookJSONRequestBody{
		DefaultCurrencyIsoCode: strings.ToUpper(mdi.BookCurrencyIsoCode),
		Name:                   bookName,
	})
	switch {
	case err != nil:
		return fmt.Errorf("creating new book via API: %w", err)
	case response.JSON422 != nil:
		return fmt.Errorf("creating new book via API: %v", response.JSON422)
	case response.JSON201 == nil:
		return fmt.Errorf("creating new book via API: %s", response.Status())
	}

	book := response.JSON201
	applog.Info(ctx, "Created new book", slog.Group("book",
		"id", book.Id,
		"name", book.Name,
	))

	mdi.book = book

	return nil
}

func (mdi *moneydanceImporter) determineBookName() string {
	if mdi.AddSuffixToBookName {
		return fmt.Sprintf("%s (imported on %s)", mdi.BookName, time.Now().Format("2006-01-02 at 15:04:05"))
	}

	return mdi.BookName
}

func (mdi *moneydanceImporter) createRegisters(ctx context.Context) error {
	// r := registerImport{
	// 	apiClient: mdi.apiClient,
	// 	book:      mdi.book,
	// 	md:        mdi.md,
	// }

	// if err := r.run(ctx); err != nil {
	// 	return fmt.Errorf("importing registers: %w", err)
	// }

	return nil
}
