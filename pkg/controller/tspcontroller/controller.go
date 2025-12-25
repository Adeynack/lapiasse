package tspcontroller

import (
	"context"

	"adeynack.net/lapiasse/pkg/api/tspapi"
)

type TspController struct{}

var _ tspapi.StrictServerInterface = (*TspController)(nil)

// BooksCreate implements tspapi.StrictServerInterface.
func (t *TspController) BooksCreate(ctx context.Context, request tspapi.BooksCreateRequestObject) (tspapi.BooksCreateResponseObject, error) {
	panic("unimplemented")
}

// BooksDelete implements tspapi.StrictServerInterface.
func (t *TspController) BooksDelete(ctx context.Context, request tspapi.BooksDeleteRequestObject) (tspapi.BooksDeleteResponseObject, error) {
	panic("unimplemented")
}

// BooksIndex implements tspapi.StrictServerInterface.
func (t *TspController) BooksIndex(ctx context.Context, request tspapi.BooksIndexRequestObject) (tspapi.BooksIndexResponseObject, error) {
	panic("unimplemented")
}

// BooksShow implements tspapi.StrictServerInterface.
func (t *TspController) BooksShow(ctx context.Context, request tspapi.BooksShowRequestObject) (tspapi.BooksShowResponseObject, error) {
	panic("unimplemented")
}

// BooksUpdate implements tspapi.StrictServerInterface.
func (t *TspController) BooksUpdate(ctx context.Context, request tspapi.BooksUpdateRequestObject) (tspapi.BooksUpdateResponseObject, error) {
	panic("unimplemented")
}

// Health implements [tspapi.StrictServerInterface.RootServiceHealth].
func (t *TspController) Health(ctx context.Context, request tspapi.HealthRequestObject) (tspapi.HealthResponseObject, error) {
	return tspapi.Health200JSONResponse(tspapi.Health{
		Status: tspapi.ServerHealthStatusHealthy,
	}), nil
}
