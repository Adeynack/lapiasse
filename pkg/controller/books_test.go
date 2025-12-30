//go:build test

package controller_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"

	"adeynack.net/lapiasse/pkg/api"
	"adeynack.net/lapiasse/pkg/app"
	"adeynack.net/lapiasse/pkg/controller"
	"adeynack.net/lapiasse/pkg/model"
	"adeynack.net/lapiasse/pkg/platform/ctxval"
	"adeynack.net/lapiasse/pkg/repository"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestBooksIndex(t *testing.T) {
	seedMoreBooksThanMaxPageSize := func(ctx context.Context) {
		db := ctxval.MustResolve[*gorm.DB](ctx)
		for _, i := range rand.Perm(controller.DefaultPageSize + 10) {
			lo.Must0(gorm.G[model.Book](db).Create(ctx, &model.Book{Name: fmt.Sprintf("Book %04d", i+1)}))
		}
	}

	responseBookNames := func(resp api.BooksIndex200JSONResponse) []string {
		return lo.Map(resp.Data, func(b api.BookShow, _ int) string { return b.Name })
	}

	for _, tc := range []struct {
		name         string
		seed         func(ctx context.Context)
		request      api.BooksIndexRequestObject
		expecting200 func(t *testing.T, resp api.BooksIndex200JSONResponse)
	}{
		{
			name: "no books exist, response is empty",
			expecting200: func(t *testing.T, resp api.BooksIndex200JSONResponse) {
				require.Empty(t, resp.Data)
			},
		},
		{
			name: "some books exist, response contains limited number of books (default page size)",
			seed: seedMoreBooksThanMaxPageSize,
			expecting200: func(t *testing.T, resp api.BooksIndex200JSONResponse) {
				expectedBookNames := lo.Map(lo.RangeFrom(1, controller.DefaultPageSize), func(i int, _ int) string {
					return fmt.Sprintf("Book %04d", i)
				})
				require.Equal(t, expectedBookNames, responseBookNames(resp))
			},
		},
		{
			name: "some books exist, request specifies limit, response contains correct books",
			seed: seedMoreBooksThanMaxPageSize,
			request: api.BooksIndexRequestObject{
				Params: api.BooksIndexParams{
					PageSize: lo.ToPtr(3),
				},
			},
			expecting200: func(t *testing.T, resp api.BooksIndex200JSONResponse) {
				require.Equal(t, []string{"Book 0001", "Book 0002", "Book 0003"}, responseBookNames(resp))
			},
		},
		{
			name: "some books exist, request specifies limit and page, response contains correct books",
			seed: seedMoreBooksThanMaxPageSize,
			request: api.BooksIndexRequestObject{
				Params: api.BooksIndexParams{
					PageSize: lo.ToPtr(3),
					Page:     lo.ToPtr(2),
				},
			},
			expecting200: func(t *testing.T, resp api.BooksIndex200JSONResponse) {
				require.Equal(t, []string{"Book 0004", "Book 0005", "Book 0006"}, responseBookNames(resp))
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := app.CreateTestAppCtx(t)
			if tc.seed != nil {
				tc.seed(ctx)
			}

			ctrl := &controller.ApplicationController{}
			resp, err := ctrl.BooksIndex(ctx, tc.request)

			require.NoError(t, err)

			if tc.expecting200 != nil {
				resp200, ok := resp.(api.BooksIndex200JSONResponse)
				require.True(t, ok, "expected BooksIndex200JSONResponse, got %T", resp)
				tc.expecting200(t, resp200)
			}
		})
	}
}

func TestBooksCreate(t *testing.T) {
	for name, tc := range map[string]struct {
		seed         func(ctx context.Context, db *gorm.DB)
		requestBook  *api.BooksCreateJSONRequestBody
		expecting201 func(t *testing.T, resp api.BooksCreate201JSONResponse)
		expecting422 func(t *testing.T, resp api.ValidationError)
	}{
		"simple book creation": {
			requestBook: &api.BooksCreateJSONRequestBody{
				Name:                   "My Book",
				DefaultCurrencyIsoCode: "EUR",
			},
		},
		"wrong currency ISO code": {
			requestBook: &api.BooksCreateJSONRequestBody{
				Name:                   "My Book",
				DefaultCurrencyIsoCode: "INVALID",
			},
			expecting422: func(t *testing.T, resp api.ValidationError) {
				lenValidErr, ok := lo.Find(resp.ValidationErrors, func(fe api.FieldValidationError) bool { return fe.Validation == "len" })
				require.True(t, ok, "expected len validation error")
				require.Equal(t, "Book.DefaultCurrencyIsoCode", lenValidErr.Field)
				require.Equal(t, "currency ISO code must be 3 characters long", lenValidErr.Message)
				require.Equal(t, "3", lo.FromPtr(lenValidErr.Param))
			},
		},
		"duplicate name": {
			seed: func(ctx context.Context, db *gorm.DB) {
				repository.MustCreate0(ctx, &model.Book{Name: "My Book", DefaultCurrencyIsoCode: "EUR"})
			},
			requestBook: &api.BooksCreateJSONRequestBody{
				Name:                   "My Book",
				DefaultCurrencyIsoCode: "USD",
			},
			expecting422: func(t *testing.T, resp api.ValidationError) {
				uniqueValidErr, ok := lo.Find(resp.ValidationErrors, func(fe api.FieldValidationError) bool { return fe.Validation == "unique" })
				require.True(t, ok, "expected unique validation error")
				require.Equal(t, "Book.Name", uniqueValidErr.Field)
				require.Equal(t, "book name must be unique", uniqueValidErr.Message)
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctx := app.CreateTestAppCtx(t)
			db := ctxval.MustResolve[*gorm.DB](ctx)

			if tc.seed != nil {
				tc.seed(ctx, db)
			}

			ctrl := &controller.ApplicationController{}
			request := api.BooksCreateRequestObject{
				Body: tc.requestBook,
			}
			resp, err := ctrl.BooksCreate(ctx, request)

			require.NoError(t, err)

			if tc.expecting201 != nil {
				resp201, ok := resp.(api.BooksCreate201JSONResponse)
				require.True(t, ok, "expected BooksCreate201JSONResponse, got %T", resp)
				tc.expecting201(t, resp201)

				// Check if book was really created in DB
				db := ctxval.MustResolve[*gorm.DB](ctx)
				book, err := gorm.G[model.Book](db).Where("id = ?", resp201.Id).First(ctx)
				require.NoError(t, err)
				require.Equal(t, tc.requestBook.Name, book.Name)
			}

			if tc.expecting422 != nil {
				resp422, ok := resp.(api.ValidationError)
				require.True(t, ok, "expected ValidationError, got %T", resp)
				tc.expecting422(t, resp422)
			}
		})
	}
}
