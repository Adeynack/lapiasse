package controller

import (
	"context"
	"errors"
	"fmt"

	"adeynack.net/lapiasse/pkg/api"
	"adeynack.net/lapiasse/pkg/model"
	"adeynack.net/lapiasse/pkg/platform/ctxval"
	"adeynack.net/lapiasse/pkg/platform/loex"
	"gorm.io/gorm"
)

// BooksIndex implements api.StrictServerInterface.
func (t *ApplicationController) BooksIndex(ctx context.Context, request api.BooksIndexRequestObject) (api.BooksIndexResponseObject, error) {
	db := ctxval.MustResolve[*gorm.DB](ctx)

	books, err := gorm.G[model.Book](db).
		Scopes(scopePaginate(request.Params.Page, request.Params.PageSize)).
		Order("books.name").
		Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading books from database: %w", err)
	}

	response := api.BooksIndex200JSONResponse{
		Data: loex.MapE(books, toApiBookShow),
	}

	return response, nil
}

// BooksShow implements api.StrictServerInterface.
func (t *ApplicationController) BooksShow(ctx context.Context, request api.BooksShowRequestObject) (api.BooksShowResponseObject, error) {
	db := ctxval.MustResolve[*gorm.DB](ctx)

	book, err := gorm.G[model.Book](db).Where("id = ?", request.BookId).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return api.BooksShow404JSONResponse(api404ErrorFromId("Book", request.BookId)), nil
	} else if err != nil {
		return nil, fmt.Errorf("reading book from database: %w", err)
	}

	response := api.BooksShow200JSONResponse(toApiBookShow(book))

	return response, nil
}

// BooksCreate implements api.StrictServerInterface.
func (t *ApplicationController) BooksCreate(ctx context.Context, request api.BooksCreateRequestObject) (api.BooksCreateResponseObject, error) {
	db := ctxval.MustResolve[*gorm.DB](ctx)

	book := model.Book{
		Name:                   request.Body.Name,
		DefaultCurrencyIsoCode: request.Body.DefaultCurrencyIsoCode,
	}

	if validErr, err := validate(ctx, book); err != nil {
		return nil, fmt.Errorf("validating book: %w", err)
	} else if validErr != nil {
		return api.BooksCreate422JSONResponse(*validErr), nil
	}

	err := gorm.G[model.Book](db).Create(ctx, &book)
	if err != nil {
		return nil, fmt.Errorf("creating book in database: %w", err)
	}

	response := api.BooksCreate201JSONResponse(toApiBookShow(book))

	return response, nil
}

// BooksUpdate implements api.StrictServerInterface.
func (t *ApplicationController) BooksUpdate(ctx context.Context, request api.BooksUpdateRequestObject) (api.BooksUpdateResponseObject, error) {
	panic("unimplemented")
}

// BooksDelete implements api.StrictServerInterface.
func (t *ApplicationController) BooksDelete(ctx context.Context, request api.BooksDeleteRequestObject) (api.BooksDeleteResponseObject, error) {
	panic("unimplemented")
}

func toApiBookShow(b model.Book) api.BookShow {
	return api.BookShow{
		CreatedAt:              b.CreatedAt,
		DefaultCurrencyIsoCode: b.DefaultCurrencyIsoCode,
		Id:                     fmt.Sprintf("%d", b.ID),
		Name:                   b.Name,
		UpdatedAt:              b.UpdatedAt,
	}
}
