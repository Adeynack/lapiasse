package controller

import (
	"context"

	"adeynack.net/lapiasse/pkg/api"
)

type ApplicationController struct {
}

// (GET /health)
func (c *ApplicationController) GetHealth(ctx context.Context, request api.GetHealthRequestObject) (api.GetHealthResponseObject, error) {
	return api.GetHealth200JSONResponse{Status: api.ServerHealthStatusHealthy}, nil
}
