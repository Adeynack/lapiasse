package tspcontroller

import (
	"context"

	"adeynack.net/lapiasse/pkg/api/tspapi"
)

type TspController struct{}

var _ tspapi.StrictServerInterface = (*TspController)(nil)

// (GET /books)
func (c *TspController) BooksList(ctx context.Context, request tspapi.BooksListRequestObject) (tspapi.BooksListResponseObject, error) {
	return tspapi.BooksList200JSONResponse{
		Items: []tspapi.Book{},
	}, nil
}

// (POST /books)
func (c *TspController) BooksCreate(ctx context.Context, request tspapi.BooksCreateRequestObject) (tspapi.BooksCreateResponseObject, error) {
	return tspapi.BooksCreate201JSONResponse{
		Item: tspapi.Book{},
	}, nil
}
