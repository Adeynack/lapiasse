package model

import (
	"time"

	"adeynack.net/lapiasse/pkg/api"
)

type Register struct {
	Base

	Name               string           `gorm:"type:text(200);not null" json:"name" validate:"required,min=1,max=200"`
	Type               api.RegisterType `gorm:"type:integer;not null" json:"type" validate:"required,register_type"`
	BookID             ID               `gorm:"not null;index" json:"book_id" validate:"required"`
	ParentID           *ID              `gorm:"index" json:"parent_id" validate:"omitempty"`                                                     // A null parent means it is a root register.
	StartsAt           time.Time        `gorm:"type:date;not null" json:"starts_at" validate:"required,datetime=2006-01-02"`                     // Opening date of the register (eg: for accounts, but not for categories).
	ExpiresAt          *time.Time       `gorm:"type:date" json:"expires_at,omitempty" validate:"omitempty,datetime=2006-01-02,gtfield=StartsAt"` // Optional expiration date of the register (eg: for a credit card).
	CurrencyIsoCode    string           `gorm:"type:char(3);not null" json:"currency_iso_code" validate:"required,len=3,uppercase"`
	Notes              string           `gorm:"type:text(2000);default:null" json:"notes,omitzero" validate:"omitempty,max=2000"`
	InitialBalance     int64            `gorm:"type:bigint;not null;default:0" json:"initial_balance" validate:"required"`
	Active             bool             `gorm:"type:boolean;not null;default:true" json:"active" validate:"required"`
	DefaultCategoryID  ID               `gorm:"index" json:"default_category_id,omitempty" validate:"required"`                             // The category automatically selected when entering a new exchange from this register.
	InstitutionName    string           `gorm:"type:text(200);default:null" json:"institution_name,omitempty" validate:"omitempty,max=200"` // Name of the institution (eg: bank) managing the registry (eg: credit card).
	AccountNumber      string           `gorm:"type:text(100);default:null" json:"account_number,omitempty" validate:"omitempty,max=100"`   // Number by which the register is referred to (eg: bank account number).
	IBAN               string           `gorm:"type:text(34);default:null" json:"iban,omitempty" validate:"omitempty,max=34"`               // In the case the register is identified by an International Bank Account Number (IBAN).
	AnnualInterestRate float32          `gorm:"type:real" json:"annual_interest_rate" validate:"omitempty,gte=0"`                           // In the case the register is being charged interests, its rate per year (eg: credit card).
	CreditLimit        int64            `gorm:"type:bigint" json:"credit_limit" validate:"omitempty,gte=0"`                                 // In the case the register has a credit limit (eg: credit card, credit margin).
	CardNumber         string           `gorm:"type:text(50);default:null" json:"card_number,omitempty" validate:"omitempty,max=50"`        // In the case the register is linked to a card, its number (eg: a credit card).

	Book     *Book      `gorm:"foreignKey:BookID" json:"book,omitempty"`
	Parent   *Register  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Register `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}
