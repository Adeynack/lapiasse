// Contains implementation of the API interface, grouped by resource.
package controller

import "adeynack.net/lapiasse/pkg/api"

type Service struct {
	ApplicationController
	BooksController
	RegistersController
	ExchangesController
}

func New() api.StrictServerInterface {
	return &Service{
		ApplicationController: ApplicationController{},
		BooksController:       BooksController{},
		RegistersController:   RegistersController{},
		ExchangesController:   ExchangesController{},
	}
}
