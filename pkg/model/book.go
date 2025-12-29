package model

import (
	"context"
	"errors"

	"adeynack.net/lapiasse/pkg/platform/ctxval"
	"gorm.io/gorm"
)

type Book struct {
	Base

	Name                   string `gorm:"type:text(255);not null" json:"name" validate:"required,min=1,max=255"`
	DefaultCurrencyIsoCode string `gorm:"type:char(3);not null" json:"default_currency_iso_code" validate:"required,currencyIsoCode"`

	// Registers []Register `gorm:"foreignKey:BookID" json:"registers,omitempty"`
}

func (b *Book) Validate(ctx context.Context, reason ValidationReason) error {
	db := ctxval.MustResolve[*gorm.DB](ctx)

	val, err := b.ValidateBase(ctx, db, reason, b)
	if err != nil {
		return err
	}

	// Name must be unique.
	// When users are introduced, should be "unique per user".
	_, err = gorm.G[Book](db).Where("name = ? AND id <> ?", b.Name, b.ID).Select("id").First(ctx)
	if err == nil {
		val.Add("Book.Name", "book name must be unique", "unique", b.Name)
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return val.ToError()
}
