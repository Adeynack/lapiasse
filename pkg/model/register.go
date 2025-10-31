package model

import (
	"time"

	"gorm.io/gorm"
)

type Register struct {
	gorm.Model
	Name            string     `gorm:"type:text(255);not null" json:"name" validate:"required,min=1,max=255"`
	ParentID        *uint      `gorm:"index" json:"-" validate:"omitempty,gt=0"`
	Parent          *Register  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	StartsAt        time.Time  `gorm:"type:date;not null" json:"starts_at" validate:"required,datetime=2006-01-02"`
	ExpiresAt       *time.Time `gorm:"type:date" json:"expires_at,omitempty" validate:"omitempty,datetime=2006-01-02,gtfield=StartsAt"`
	CurrencyIsoCode string     `gorm:"type:char(3);not null" json:"currency_iso_code" validate:"required,len=3,uppercase"`
	Notes           string     `gorm:"type:text;default:null" json:"notes,omitzero" validate:"omitempty"`
}
