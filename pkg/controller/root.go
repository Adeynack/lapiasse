package controller

import (
	"context"

	"adeynack.net/lapiasse/pkg/api"
)

type ApplicationController struct{}

func New() api.StrictServerInterface {
	return &ApplicationController{}
}

var _ api.StrictServerInterface = (*ApplicationController)(nil)

func (t *ApplicationController) Health(ctx context.Context, request api.HealthRequestObject) (api.HealthResponseObject, error) {
	return api.Health200JSONResponse(api.Health{
		Status: api.ServerHealthStatusHealthy,
	}), nil
}
