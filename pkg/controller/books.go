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

// GetBook implements [api.StrictServerInterface.GetBook].
func (t *ApplicationController) GetBook(ctx context.Context, request api.GetBookRequestObject) (api.GetBookResponseObject, error) {
	db := ctxval.MustResolve[*gorm.DB](ctx)

	book, err := gorm.G[model.Book](db).Where("id = ?", request.BookId).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return api404ErrorFromId("Book", request.BookId), nil
	} else if err != nil {
		return nil, fmt.Errorf("reading book from database: %w", err)
	}

	return api.GetBook200JSONResponse(toApiBookShow(book)), nil
}

// CreateBook implements [api.StrictServerInterface.CreateBook].
func (t *ApplicationController) CreateBook(ctx context.Context, request api.CreateBookRequestObject) (api.CreateBookResponseObject, error) {
	db := ctxval.MustResolve[*gorm.DB](ctx)

	var book model.Book
	book.AssignAttributes(request.Body)

	if done, res, err := validate(ctx, &book, model.ValidationReasonCreate); done {
		return res, err
	}

	err := gorm.G[model.Book](db).Create(ctx, &book)
	if err != nil {
		return nil, fmt.Errorf("creating book in database: %w", err)
	}

	return api.CreateBook201JSONResponse(toApiBookShow(book)), nil
}

// UpdateBook implements [api.StrictServerInterface.UpdateBook].
func (t *ApplicationController) UpdateBook(ctx context.Context, request api.UpdateBookRequestObject) (api.UpdateBookResponseObject, error) {
	db := ctxval.MustResolve[*gorm.DB](ctx)

	book, err := gorm.G[model.Book](db).Where("id = ?", request.BookId).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return api404ErrorFromId("Book", request.BookId), nil
	} else if err != nil {
		return nil, fmt.Errorf("reading book from database: %w", err)
	}

	book.AssignAttributes(request.Body)

	if done, res, err := validate(ctx, &book, model.ValidationReasonUpdate); done {
		return res, err
	}

	_, err = gorm.G[*model.Book](db).Updates(ctx, &book)
	if err != nil {
		return nil, fmt.Errorf("updating book in database: %w", err)
	}

	return api.UpdateBook200JSONResponse(toApiBookShow(book)), nil
}

// DeleteBook implements [api.StrictServerInterface.DeleteBook].
func (t *ApplicationController) DeleteBook(ctx context.Context, request api.DeleteBookRequestObject) (api.DeleteBookResponseObject, error) {
	db := ctxval.MustResolve[*gorm.DB](ctx)

	rowsAffected, err := gorm.G[model.Book](db).Where("id = ?", request.BookId).Delete(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading book from database: %w", err)
	}
	if rowsAffected == 0 {
		return api404ErrorFromId("Book", request.BookId), nil
	}

	return api.DeleteBook204Response{}, nil
}

func toApiBookShow(b model.Book) api.BookShow {
	return api.BookShow{
		Id:                     apiId(b.ID),
		CreatedAt:              b.CreatedAt,
		UpdatedAt:              b.UpdatedAt,
		Name:                   b.Name,
		DefaultCurrencyIsoCode: b.DefaultCurrencyIsoCode,
	}
}
