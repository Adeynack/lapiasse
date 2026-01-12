package model

import (
	"context"
	"errors"

	"adeynack.net/lapiasse/pkg/api"
	"adeynack.net/lapiasse/pkg/platform/ctxval"
	"gorm.io/gorm"
)

type Book struct {
	Base

	Name                   string `gorm:"type:text(255);not null" json:"name" validate:"required,min=1,max=255"`
	DefaultCurrencyIsoCode string `gorm:"type:char(3);not null" json:"default_currency_iso_code" validate:"required,currencyIsoCode"`

	// Registers []Register `gorm:"foreignKey:BookID" json:"registers,omitempty"`
}

func (b *Book) AssignAttributes(attr *api.BookEdit) {
	b.Name = attr.Name
	b.DefaultCurrencyIsoCode = attr.DefaultCurrencyIsoCode
}

func (b *Book) Validate(ctx context.Context, reason ValidationReason) error {
	return b.BaseValidate(ctx, reason, b, func(ctx context.Context, reason ValidationReason, validationErrors ValidationErrorBuilder) error {
		db := ctxval.MustResolve[*gorm.DB](ctx)

		// Name must be unique.
		// When users are introduced, should be "unique per user".
		if bookWithSameName, err := gorm.G[Book](db).Where("name = ?", b.Name).Select("id").First(ctx); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// No book with same name exists, all good.
			} else {
				return err
			}
		} else if reason == ValidationReasonCreate || bookWithSameName.ID != b.ID {
			validationErrors.Add("Book.Name", "book name must be unique", "unique", b.Name)
		}

		return nil
	})
}
