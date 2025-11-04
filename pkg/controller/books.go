package controller

import (
	"context"

	"adeynack.net/lapiasse/pkg/api"
)

type BooksController struct {
}

// (GET /books)
func (c *BooksController) GetBooks(ctx context.Context, request api.GetBooksRequestObject) (api.GetBooksResponseObject, error) {
	books := []api.Book{}

	return api.GetBooks200JSONResponse{Books: books}, nil
}

// (POST /books)
func (c *BooksController) CreateBook(ctx context.Context, request api.CreateBookRequestObject) (api.CreateBookResponseObject, error) {
	book := api.Book{}

	return api.CreateBook201JSONResponse{Book: book}, nil
}

// (GET /books/{bookId})
func (c *BooksController) GetBook(ctx context.Context, request api.GetBookRequestObject) (api.GetBookResponseObject, error) {
	book := api.Book{}

	return api.GetBook200JSONResponse{Book: book}, nil
}
