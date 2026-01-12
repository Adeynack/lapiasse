//go:build test

package controller_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"testing"
	"testing/synctest"
	"time"

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

func TestListBooks(t *testing.T) {
	seedMoreBooksThanMaxPageSize := func(ctx context.Context) {
		db := ctxval.MustResolve[*gorm.DB](ctx)
		for _, i := range rand.Perm(controller.DefaultPageSize + 10) {
			lo.Must0(gorm.G[model.Book](db).Create(ctx, &model.Book{Name: fmt.Sprintf("Book %04d", i+1)}))
		}
	}

	responseBookNames := func(resp api.ListBooks200JSONResponse) []string {
		return lo.Map(resp.Books, func(b api.BookShow, _ int) string { return b.Name })
	}

	for _, tc := range []struct {
		name         string
		seed         func(ctx context.Context)
		request      api.ListBooksRequestObject
		expecting200 func(t *testing.T, resp api.ListBooks200JSONResponse)
	}{
		{
			name: "no books exist, response is empty",
			expecting200: func(t *testing.T, resp api.ListBooks200JSONResponse) {
				require.Empty(t, resp.Books)
			},
		},
		{
			name: "some books exist, response contains limited number of books (default page size)",
			seed: seedMoreBooksThanMaxPageSize,
			expecting200: func(t *testing.T, resp api.ListBooks200JSONResponse) {
				expectedBookNames := lo.Map(lo.RangeFrom(1, controller.DefaultPageSize), func(i int, _ int) string {
					return fmt.Sprintf("Book %04d", i)
				})
				require.Equal(t, expectedBookNames, responseBookNames(resp))
			},
		},
		{
			name: "some books exist, request specifies limit, response contains correct books",
			seed: seedMoreBooksThanMaxPageSize,
			request: api.ListBooksRequestObject{
				Params: api.ListBooksParams{
					PageSize: lo.ToPtr(3),
				},
			},
			expecting200: func(t *testing.T, resp api.ListBooks200JSONResponse) {
				require.Equal(t, []string{"Book 0001", "Book 0002", "Book 0003"}, responseBookNames(resp))
			},
		},
		{
			name: "some books exist, request specifies limit and page, response contains correct books",
			seed: seedMoreBooksThanMaxPageSize,
			request: api.ListBooksRequestObject{
				Params: api.ListBooksParams{
					PageSize: lo.ToPtr(3),
					Page:     lo.ToPtr(2),
				},
			},
			expecting200: func(t *testing.T, resp api.ListBooks200JSONResponse) {
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
			resp, err := ctrl.ListBooks(ctx, tc.request)

			require.NoError(t, err)

			if tc.expecting200 != nil {
				resp200, ok := resp.(api.ListBooks200JSONResponse)
				require.True(t, ok, "expected ListBooks200JSONResponse, got %T", resp)
				tc.expecting200(t, resp200)
			}
		})
	}
}

func TestCreateBook(t *testing.T) {
	timeAtTestStart := time.Now()

	for name, tc := range map[string]struct {
		seed         func(ctx context.Context, db *gorm.DB)
		requestBook  *api.CreateBookJSONRequestBody
		expecting201 func(t *testing.T, resp api.CreateBook201JSONResponse, book model.Book)
		expecting422 func(t *testing.T, resp api.ValidationError)
	}{
		"simple book creation": {
			requestBook: &api.CreateBookJSONRequestBody{
				Name:                   "My Book",
				DefaultCurrencyIsoCode: "EUR",
			},
			expecting201: func(t *testing.T, resp api.CreateBook201JSONResponse, book model.Book) {
				require.GreaterOrEqual(t, resp.CreatedAt, timeAtTestStart)
				require.Equal(t, resp.UpdatedAt, resp.CreatedAt)
				require.Equal(t, "My Book", book.Name)
				require.Equal(t, "EUR", book.DefaultCurrencyIsoCode)

				require.GreaterOrEqual(t, book.CreatedAt, timeAtTestStart)
				require.Equal(t, book.UpdatedAt, book.CreatedAt)
				require.Equal(t, "My Book", book.Name)
				require.Equal(t, "EUR", book.DefaultCurrencyIsoCode)
			},
		},
		"wrong currency ISO code": {
			requestBook: &api.CreateBookJSONRequestBody{
				Name:                   "My Book",
				DefaultCurrencyIsoCode: "INVALID",
			},
			expecting422: func(t *testing.T, resp api.ValidationError) {
				lenValidErr, ok := lo.Find(resp.ValidationErrors, func(fe api.FieldValidationFailure) bool { return fe.Validation == "len" })
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
			requestBook: &api.CreateBookJSONRequestBody{
				Name:                   "My Book",
				DefaultCurrencyIsoCode: "USD",
			},
			expecting422: func(t *testing.T, resp api.ValidationError) {
				uniqueValidErr, ok := lo.Find(resp.ValidationErrors, func(fe api.FieldValidationFailure) bool { return fe.Validation == "unique" })
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
			request := api.CreateBookRequestObject{
				Body: tc.requestBook,
			}
			resp, err := ctrl.CreateBook(ctx, request)

			require.NoError(t, err)

			if tc.expecting422 != nil {
				resp422, ok := resp.(api.ValidationError)
				require.True(t, ok, "expected ValidationError, got %T", resp)
				tc.expecting422(t, resp422)

				return
			}

			require.NotNil(t, tc.expecting201)

			resp201, ok := resp.(api.CreateBook201JSONResponse)
			require.True(t, ok, "expected CreateBook201JSONResponse, got %T", resp)

			// Check if book was really created in DB
			book, err := gorm.G[model.Book](db).Where("id = ?", resp201.Id).First(ctx)
			require.NoError(t, err)

			tc.expecting201(t, resp201, book)
		})
	}
}

func TestUpdateBook(t *testing.T) {
	seedMyBook := func(ctx context.Context, db *gorm.DB) *model.Book {
		return repository.MustCreate(ctx, &model.Book{
			Base:                   model.Base{ID: model.ID(1)},
			Name:                   "My Book",
			DefaultCurrencyIsoCode: "EUR",
		})
	}

	for name, tc := range map[string]struct {
		seed         func(ctx context.Context, db *gorm.DB) *model.Book
		request      api.UpdateBookRequestObject
		expecting200 func(t *testing.T, resp api.UpdateBook200JSONResponse, original, updated *model.Book)
		expecting404 func(t *testing.T, resp api.NotFoundError)
		expecting422 func(t *testing.T, resp api.ValidationError)
	}{
		"invalid book ID": {
			request: api.UpdateBookRequestObject{
				BookId: api.ID("does not exist"),
				Body:   &api.BookEdit{},
			},
			expecting404: func(t *testing.T, resp api.NotFoundError) {
				require.Equal(t, lo.ToPtr(`Book with ID "does not exist" not found`), resp.Detail)
			},
		},
		"no change still updates the updated_at": {
			seed: seedMyBook,
			request: api.UpdateBookRequestObject{
				BookId: api.ID("1"),
				Body: &api.BookEdit{
					Name:                   "My Book",
					DefaultCurrencyIsoCode: "EUR",
				},
			},
			expecting200: func(t *testing.T, resp api.UpdateBook200JSONResponse, original, updated *model.Book) {
				require.Equal(t, original.CreatedAt.Local(), updated.CreatedAt.Local(), "created_at should not have changed")
				require.Greater(t, updated.UpdatedAt.Local(), original.UpdatedAt.Local(), "updated_at should have been updated")

				require.Equal(t, "My Book", updated.Name)
				require.Equal(t, "EUR", updated.DefaultCurrencyIsoCode)

				require.Equal(t, updated.CreatedAt.Local(), resp.CreatedAt.Local())
				require.Equal(t, updated.UpdatedAt.Local(), resp.UpdatedAt.Local())
				require.Equal(t, "My Book", resp.Name)
				require.Equal(t, "EUR", resp.DefaultCurrencyIsoCode)
			},
		},
		"change name to another valid": {
			seed: seedMyBook,
			request: api.UpdateBookRequestObject{
				BookId: api.ID("1"),
				Body: &api.BookEdit{
					Name:                   "Another Name",
					DefaultCurrencyIsoCode: "EUR",
				},
			},
			expecting200: func(t *testing.T, resp api.UpdateBook200JSONResponse, original, updated *model.Book) {
				require.Greater(t, updated.UpdatedAt.Local(), original.UpdatedAt.Local())
				require.Equal(t, "Another Name", updated.Name)
				require.Equal(t, "EUR", updated.DefaultCurrencyIsoCode)
			},
		},
		"change to duplicate name": {
			seed: func(ctx context.Context, db *gorm.DB) *model.Book {
				repository.MustCreate0(ctx, &model.Book{
					Base:                   model.Base{ID: model.ID(999)},
					Name:                   "Tatu Tata",
					DefaultCurrencyIsoCode: "USD",
				})

				return seedMyBook(ctx, db)
			},
			request: api.UpdateBookRequestObject{
				BookId: api.ID("1"),
				Body: &api.BookEdit{
					Name:                   "Tatu Tata",
					DefaultCurrencyIsoCode: "EUR",
				},
			},
			expecting422: func(t *testing.T, resp api.ValidationError) {
				uniqueValidErr, ok := lo.Find(resp.ValidationErrors, func(fe api.FieldValidationFailure) bool { return fe.Validation == "unique" })
				require.True(t, ok, "expected unique validation error")
				require.Equal(t, "Book.Name", uniqueValidErr.Field)
				require.Equal(t, "book name must be unique", uniqueValidErr.Message)
			},
		},
		"change to wrong currency ISO code": {
			seed: seedMyBook,
			request: api.UpdateBookRequestObject{
				BookId: api.ID("1"),
				Body: &api.BookEdit{
					Name:                   "My Book",
					DefaultCurrencyIsoCode: "INVALID",
				},
			},
			expecting422: func(t *testing.T, resp api.ValidationError) {
				lenValidErr, ok := lo.Find(resp.ValidationErrors, func(fe api.FieldValidationFailure) bool { return fe.Validation == "len" })
				require.True(t, ok, "expected len validation error")
				require.Equal(t, "Book.DefaultCurrencyIsoCode", lenValidErr.Field)
				require.Equal(t, "currency ISO code must be 3 characters long", lenValidErr.Message)
				require.Equal(t, "3", lo.FromPtr(lenValidErr.Param))
			},
		},
		"change to another valid currency ISO code": {
			seed: seedMyBook,
			request: api.UpdateBookRequestObject{
				BookId: api.ID("1"),
				Body: &api.BookEdit{
					Name:                   "My Book",
					DefaultCurrencyIsoCode: "USD", // from "EUR"
				},
			},
			expecting200: func(t *testing.T, resp api.UpdateBook200JSONResponse, original, updated *model.Book) {
				require.Greater(t, updated.UpdatedAt.Local(), original.UpdatedAt.Local())
				require.Equal(t, "My Book", updated.Name)
				require.Equal(t, "USD", updated.DefaultCurrencyIsoCode)
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				ctx := app.CreateTestAppCtx(t)
				db := ctxval.MustResolve[*gorm.DB](ctx)

				var original *model.Book
				if tc.seed != nil {
					original = tc.seed(ctx, db)
					time.Sleep(5 * time.Second) // simulate time passing to see updated_at change
				}

				ctrl := &controller.ApplicationController{}
				resp, err := ctrl.UpdateBook(ctx, tc.request)

				require.NoError(t, err)

				if tc.expecting404 != nil {
					require.IsType(t, api.NotFoundError{}, resp)
					resp404 := resp.(api.NotFoundError)
					tc.expecting404(t, resp404)
				}

				if tc.expecting422 != nil {
					require.IsType(t, api.ValidationError{}, resp)
					resp422 := resp.(api.ValidationError)
					tc.expecting422(t, resp422)
				}

				if tc.expecting200 != nil {
					require.IsType(t, api.UpdateBook200JSONResponse{}, resp)
					resp201 := resp.(api.UpdateBook200JSONResponse)

					// Check if updated was really created in DB
					updated, err := gorm.G[model.Book](db).Where("id = ?", resp201.Id).First(ctx)
					require.NoError(t, err)

					tc.expecting200(t, resp201, original, &updated)
				}
			})
		})
	}
}

func TestDeleteBook(t *testing.T) {
	const existingBookId = model.ID(1)

	for name, tc := range map[string]struct {
		request      api.DeleteBookRequestObject
		expecting204 func(t *testing.T, resp api.DeleteBook204Response)
		expecting404 func(t *testing.T, resp api.NotFoundError)
	}{
		"deleting existing book returns 204": {
			request: api.DeleteBookRequestObject{
				BookId: api.ID(fmt.Sprintf("%d", existingBookId)),
			},
			expecting204: func(t *testing.T, resp api.DeleteBook204Response) {
				// nothing to check in 204 response
			},
		},
		"deleting non-existing book returns 404": {
			request: api.DeleteBookRequestObject{
				BookId: api.ID("9999"),
			},
			expecting404: func(t *testing.T, resp api.NotFoundError) {
				require.Equal(t, lo.ToPtr(`Book with ID "9999" not found`), resp.Detail)
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			ctx := app.CreateTestAppCtx(t)

			repository.MustCreate(ctx, &model.Book{
				Base:                   model.Base{ID: existingBookId},
				Name:                   "My Book",
				DefaultCurrencyIsoCode: "EUR",
			})

			ctrl := &controller.ApplicationController{}
			resp, err := ctrl.DeleteBook(ctx, tc.request)
			require.NoError(t, err)

			if tc.expecting204 != nil {
				require.IsType(t, resp, api.DeleteBook204Response{})
				resp204 := resp.(api.DeleteBook204Response)
				tc.expecting204(t, resp204)
			}

			if tc.expecting404 != nil {
				require.IsType(t, resp, api.NotFoundError{})
				resp404 := resp.(api.NotFoundError)
				tc.expecting404(t, resp404)
			}
		})
	}
}
