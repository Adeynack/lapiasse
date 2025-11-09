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
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestGetBooks(t *testing.T) {
	seedMoreBooksThanMaxPageSize := func(ctx context.Context) {
		db := ctxval.MustResolve[*gorm.DB](ctx)
		for _, i := range rand.Perm(controller.DefaultPageSize + 10) {
			lo.Must0(gorm.G[model.Book](db).Create(ctx, &model.Book{Name: fmt.Sprintf("Book %04d", i+1)}))
		}
	}

	responseBookNames := func(resp api.GetBooks200JSONResponse) []string {
		return lo.Map(resp.Books, func(b api.Book, _ int) string { return b.Name })
	}

	for _, tc := range []struct {
		name         string
		seed         func(ctx context.Context)
		request      api.GetBooksRequestObject
		expecting200 func(t *testing.T, resp api.GetBooks200JSONResponse)
	}{
		{
			name: "no books exist, response is empty",
			expecting200: func(t *testing.T, resp api.GetBooks200JSONResponse) {
				require.Empty(t, resp.Books)
			},
		},
		{
			name: "some books exist, response contains limited number of books (default page size)",
			seed: seedMoreBooksThanMaxPageSize,
			expecting200: func(t *testing.T, resp api.GetBooks200JSONResponse) {
				expectedBookNames := lo.Map(lo.RangeFrom(1, controller.DefaultPageSize), func(i int, _ int) string {
					return fmt.Sprintf("Book %04d", i)
				})
				require.Equal(t, expectedBookNames, responseBookNames(resp))
			},
		},
		{
			name: "some books exist, request specifies limit, response contains correct books",
			seed: seedMoreBooksThanMaxPageSize,
			request: api.GetBooksRequestObject{
				Params: api.GetBooksParams{
					PageSize: lo.ToPtr(3),
				},
			},
			expecting200: func(t *testing.T, resp api.GetBooks200JSONResponse) {
				require.Equal(t, []string{"Book 0001", "Book 0002", "Book 0003"}, responseBookNames(resp))
			},
		},
		{
			name: "some books exist, request specifies limit and page, response contains correct books",
			seed: seedMoreBooksThanMaxPageSize,
			request: api.GetBooksRequestObject{
				Params: api.GetBooksParams{
					PageSize: lo.ToPtr(3),
					Page:     lo.ToPtr(2),
				},
			},
			expecting200: func(t *testing.T, resp api.GetBooks200JSONResponse) {
				require.Equal(t, []string{"Book 0004", "Book 0005", "Book 0006"}, responseBookNames(resp))
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ctx := app.CreateTestAppCtx(t)
			if tc.seed != nil {
				tc.seed(ctx)
			}

			ctrl := &controller.BooksController{}
			resp, err := ctrl.GetBooks(ctx, tc.request)

			require.NoError(t, err)

			if tc.expecting200 != nil {
				resp200, ok := resp.(api.GetBooks200JSONResponse)
				require.True(t, ok, "expected GetBooks200JSONResponse, got %T", resp)
				tc.expecting200(t, resp200)
			}
		})
	}
}

func TestCreateBook(t *testing.T) {
	for name, tc := range map[string]struct {
		seed         func(ctx context.Context)
		requestBook  api.BookProperties
		expecting201 func(t *testing.T, resp api.CreateBook201JSONResponse)
		expecting422 func(t *testing.T, resp api.CreateBook422JSONResponse)
	}{
		"simple book creation": {
			requestBook: api.BookProperties{
				Name:                   "My Book",
				DefaultCurrencyIsoCode: "EUR",
			},
		},
		"wrong currency ISO code": {
			requestBook: api.BookProperties{
				Name:                   "My Book",
				DefaultCurrencyIsoCode: "INVALID",
			},
			expecting422: func(t *testing.T, resp api.CreateBook422JSONResponse) {
				lenValidErr, ok := lo.Find(resp.ValidationErrors, func(fe api.FieldValidationError) bool { return fe.Validation == "len" })
				require.True(t, ok, "expected len validation error")
				require.Equal(t, "Book.DefaultCurrencyIsoCode", lenValidErr.Field)
				require.Equal(t, "currency ISO code must be 3 characters long", lenValidErr.Message)
				require.Equal(t, "3", lo.FromPtr(lenValidErr.Param))
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctx := app.CreateTestAppCtx(t)

			if tc.seed != nil {
				tc.seed(ctx)
			}

			ctrl := &controller.BooksController{}
			request := api.CreateBookRequestObject{
				Body: &api.CreateBookJSONRequestBody{
					Book: tc.requestBook,
				},
			}
			resp, err := ctrl.CreateBook(ctx, request)

			require.NoError(t, err)

			if tc.expecting201 != nil {
				resp201, ok := resp.(api.CreateBook201JSONResponse)
				require.True(t, ok, "expected CreateBook201JSONResponse, got %T", resp)
				tc.expecting201(t, resp201)

				// Check if book was really created in DB
				db := ctxval.MustResolve[*gorm.DB](ctx)
				book, err := gorm.G[model.Book](db).Where("id = ?", resp201.Book.Id).First(ctx)
				require.NoError(t, err)
				require.Equal(t, tc.requestBook.Name, book.Name)
			}

			if tc.expecting422 != nil {
				resp422, ok := resp.(api.CreateBook422JSONResponse)
				require.True(t, ok, "expected CreateBook422JSONResponse, got %T", resp)
				tc.expecting422(t, resp422)
			}
		})
	}
}
