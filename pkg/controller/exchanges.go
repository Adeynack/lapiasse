package controller

import (
	"context"

	"adeynack.net/lapiasse/pkg/api"
)

type ExchangesController struct {
}

// (GET /exchanges)
func (c *ExchangesController) GetExchanges(ctx context.Context, request api.GetExchangesRequestObject) (api.GetExchangesResponseObject, error) {
	exchanges := []api.ExchangeWithSplits{}

	return api.GetExchanges200JSONResponse{Exchanges: exchanges}, nil
}
