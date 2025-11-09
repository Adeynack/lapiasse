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

type BooksController struct {
}

// (GET /books)
func (c *BooksController) GetBooks(ctx context.Context, request api.GetBooksRequestObject) (api.GetBooksResponseObject, error) {
	db := ctxval.MustResolve[*gorm.DB](ctx)

	books, err := gorm.G[model.Book](db).
		Scopes(scopePaginate(request.Params.Page, request.Params.PageSize)).
		Order("books.name").
		Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading books from database: %w", err)
	}

	return api.GetBooks200JSONResponse{Books: loex.MapE(books, toApiBook)}, nil
}

// (POST /books)
func (c *BooksController) CreateBook(ctx context.Context, request api.CreateBookRequestObject) (api.CreateBookResponseObject, error) {
	db := ctxval.MustResolve[*gorm.DB](ctx)

	p := request.Body.Book

	book := model.Book{
		Name:                   p.Name,
		DefaultCurrencyIsoCode: p.DefaultCurrencyIsoCode,
	}

	if ok, err := validate(ctx, book); !ok {
		return api.CreateBook422JSONResponse(err), nil
	}

	err := gorm.G[model.Book](db).Create(ctx, &book)
	if err != nil {
		return nil, fmt.Errorf("creating book in database: %w", err)
	}

	return api.CreateBook201JSONResponse{Book: toApiBook(book)}, nil
}

// (GET /books/{bookId})
func (c *BooksController) GetBook(ctx context.Context, request api.GetBookRequestObject) (api.GetBookResponseObject, error) {
	db := ctxval.MustResolve[*gorm.DB](ctx)

	book, err := gorm.G[model.Book](db).Where("id = ?", request.BookId).First(ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return api.GetBook404JSONResponse(api404ErrorFromId("Book", request.BookId)), nil
	} else if err != nil {
		return nil, fmt.Errorf("reading book from database: %w", err)
	}

	return api.GetBook200JSONResponse{Book: toApiBook(book)}, nil
}

func toApiBook(b model.Book) api.Book {
	return api.Book{
		CreatedAt:              b.CreatedAt,
		DefaultCurrencyIsoCode: b.DefaultCurrencyIsoCode,
		Id:                     fmt.Sprintf("%d", b.ID),
		Name:                   b.Name,
		UpdatedAt:              b.UpdatedAt,
	}
}
