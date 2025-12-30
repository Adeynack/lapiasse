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

// ListBooks implements [api.StrictServerInterface.ListBooks].
func (t *ApplicationController) ListBooks(ctx context.Context, request api.ListBooksRequestObject) (api.ListBooksResponseObject, error) {
	db := ctxval.MustResolve[*gorm.DB](ctx)

	books, err := gorm.G[model.Book](db).
		Scopes(scopePaginate(request.Params.Page, request.Params.PageSize)).
		Order("books.name").
		Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading books from database: %w", err)
	}

	response := api.ListBooks200JSONResponse{
		Books: loex.MapE(books, toApiBookShow),
	}

	return response, nil
}

// BooksShow implements [api.StrictServerInterface.BooksShow].
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

// CreateBook implements [api.StrictServerInterface.CreateBook].
func (t *ApplicationController) CreateBook(ctx context.Context, request api.CreateBookRequestObject) (api.CreateBookResponseObject, error) {
	db := ctxval.MustResolve[*gorm.DB](ctx)

	book := model.Book{
		Name:                   request.Body.Name,
		DefaultCurrencyIsoCode: request.Body.DefaultCurrencyIsoCode,
	}

	if done, res, err := validate(ctx, &book, model.ValidationReasonCreate); done {
		return res, err
	}

	err := gorm.G[model.Book](db).Create(ctx, &book)
	if err != nil {
		return nil, fmt.Errorf("creating book in database: %w", err)
	}

	response := api.CreateBook201JSONResponse(toApiBookShow(book))

	return response, nil
}

// BooksUpdate implements [api.StrictServerInterface.BooksUpdate].
func (t *ApplicationController) BooksUpdate(ctx context.Context, request api.BooksUpdateRequestObject) (api.BooksUpdateResponseObject, error) {
	panic("unimplemented")
}

// BooksDelete implements [api.StrictServerInterface.BooksDelete].
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
