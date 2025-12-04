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
)

type moneydanceImporter struct {
	ApiEndpoint          string
	MoneydanceExportPath string
	BookName             string
	AddSuffixToBookName  bool
	BookCurrencyIsoCode  string

	apiClient *api.ClientWithResponses
	book      api.Book
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

	for _, step := range []func(context.Context) (context.Context, error){
		mdi.loadMoneydanceExport,
		mdi.createApiClient,
		mdi.createNewBook,
		mdi.createRegisters,
	} {
		var err error
		if ctx, err = step(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (mdi *moneydanceImporter) loadMoneydanceExport(ctx context.Context) (context.Context, error) {
	file, err := os.Open(mdi.MoneydanceExportPath)
	if err != nil {
		return ctx, fmt.Errorf("opening Moneydance export file: %w", err)
	}
	defer file.Close()

	var mdExport moneydance.Export
	if err := json.UnmarshalRead(file, &mdExport); err != nil {
		return ctx, fmt.Errorf("decoding Moneydance export file: %w", err)
	}

	mdi.md = &mdExport

	return ctx, nil
}

func (mdi *moneydanceImporter) createApiClient(ctx context.Context) (context.Context, error) {
	apiClient, err := api.NewClientWithResponses(mdi.ApiEndpoint)
	if err != nil {
		return ctx, fmt.Errorf("creating the API client: %w", err)
	}
	mdi.apiClient = apiClient

	return ctx, nil
}

func (mdi *moneydanceImporter) createNewBook(ctx context.Context) (context.Context, error) {
	bookName := mdi.determineBookName()

	response, err := mdi.apiClient.CreateBookWithResponse(ctx, api.CreateBookJSONRequestBody{
		Book: api.BookProperties{
			DefaultCurrencyIsoCode: strings.ToUpper(mdi.BookCurrencyIsoCode),
			Name:                   bookName,
		},
	})
	if err != nil {
		return ctx, fmt.Errorf("creating new book via API: %w", err)
	}
	if response.JSON422 != nil {
		return ctx, fmt.Errorf("creating new book via API: %v", response.JSON422)
	}
	if response.JSON201 == nil {
		return ctx, fmt.Errorf("creating new book via API: %s", response.Status())
	}

	applog.Info(ctx, "Created new book", slog.Group("book",
		"id", response.JSON201.Book.Id,
		"name", bookName,
	))
	mdi.book = response.JSON201.Book

	return ctx, nil
}

func (mdi *moneydanceImporter) determineBookName() string {
	if mdi.AddSuffixToBookName {
		return fmt.Sprintf("%s (imported on %s)", mdi.BookName, time.Now().Format("2006-01-02 at 15:04:05"))
	}

	return mdi.BookName
}

func (mdi *moneydanceImporter) createRegisters(ctx context.Context) (context.Context, error) {
	r := registerImport{
		apiClient: mdi.apiClient,
		book:      mdi.book,
		md:        mdi.md,
	}

	if err := r.run(ctx); err != nil {
		return ctx, fmt.Errorf("importing registers: %w", err)
	}

	return ctx, nil
}
