package controller

import (
	"context"
	"fmt"

	"adeynack.net/lapiasse/pkg/api"
	"adeynack.net/lapiasse/pkg/model"
	"adeynack.net/lapiasse/pkg/platform/ctxval"
	"adeynack.net/lapiasse/pkg/platform/loex"
	"gorm.io/gorm"
)

type RegistersController struct {
}

// Implements [api.StrictServerInterface.GetBookRegisters]
func (c *RegistersController) GetBookRegisters(ctx context.Context, request api.GetBookRegistersRequestObject) (api.GetBookRegistersResponseObject, error) {
	db := ctxval.MustResolve[*gorm.DB](ctx)

	registers, err := gorm.G[model.Register](db).
		Where("book_id = ?", request.BookId).
		Scopes(scopePaginate(request.Params.Page, request.Params.PageSize)).
		Order("registers.name").
		Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("reading registers from database: %w", err)
	}

	return api.GetBookRegisters200JSONResponse{Registers: loex.MapE(registers, toApiRegister)}, nil
}

func toApiRegister(r model.Register) api.Register {
	return api.Register{
		AccountType:     r.Type,
		CreatedAt:       r.CreatedAt,
		CurrencyIsoCode: r.CurrencyIsoCode,
		Id:              r.ID.String(),
		Name:            r.Name,
		UpdatedAt:       r.UpdatedAt,
	}
}

// (GET /books/{bookId}/registers)
func (c *RegistersController) CreateBookRegister(ctx context.Context, request api.CreateBookRegisterRequestObject) (api.CreateBookRegisterResponseObject, error) {
	return nil, fmt.Errorf("not implemented")
}
