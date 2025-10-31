package model

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Name                   string `gorm:"type:text(255);not null" json:"name" validate:"required,min=1,max=255"`
	DefaultCurrencyIsoCode string `gorm:"type:char(3);not null" json:"default_currency_iso_code" validate:"required,len=3,uppercase"`
}
