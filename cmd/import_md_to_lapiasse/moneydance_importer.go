package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"adeynack.net/lapiasse/pkg/api"
)

type moneydanceImporter struct {
	ApiEndpoint          string
	MoneydanceExportPath string
	BookName             string
	AddSuffixToBookName  bool
	BookCurrencyIsoCode  string

	apiClient *api.ClientWithResponses
	book      api.Book
}

func (mdi *moneydanceImporter) Start(ctx context.Context) error {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	logger.Info(
		"Importing MoneyDance Export JSON to La Piasse",
		slog.String("api-endpoint", mdi.ApiEndpoint),
		slog.String("md-source", mdi.MoneydanceExportPath),
		slog.String("book-name", mdi.BookName),
		slog.Bool("add-suffix", mdi.AddSuffixToBookName),
	)

	for _, step := range []func(context.Context) (context.Context, error){
		mdi.createApiClient,
		mdi.createNewBook,
	} {
		var err error
		if ctx, err = step(ctx); err != nil {
			return err
		}
	}

	return nil
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
		return ctx, fmt.Errorf("creating new book via API failed: %s", response.Status())
	}

	slog.Info("Created new book", slog.Group("book",
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
