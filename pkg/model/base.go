package model

import (
	"encoding/json/jsontext"
	"fmt"
	"strconv"
	"time"

	"adeynack.net/lapiasse/pkg/api"
	"adeynack.net/lapiasse/pkg/appvalidator"
	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

// ID represents the unique identifier type of a model.
type ID uint64

// Implements [fmt.Stringer.String].
func (i ID) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

// Implements [json.MarshalerTo.MarshalJSONTo].
func (i ID) MarshalJSONTo(e *jsontext.Encoder) error {
	return e.WriteToken(jsontext.String(i.String()))
}

// Implements [json.ValueUnmarshaler.UnmarshalJSON].
func (i *ID) UnmarshalJSONFrom(d *jsontext.Decoder) error {
	val, err := d.ReadValue()
	if err != nil {
		return fmt.Errorf("reading ID value from JSON decoder: %w", err)
	}

	if val.Kind() != '"' {
		return fmt.Errorf("invalid kind for ID, expected string but got: %v", val.Kind())
	}

	strVal := val.String()
	strVal = strVal[1 : len(strVal)-1] // Remove surrounding quotes

	parsedID, err := strconv.ParseUint(strVal, 10, 64)
	if err != nil {
		return fmt.Errorf("parsing ID value from string: %w", err)
	}

	*i = ID(parsedID)

	return err
}

// Base model including common fields.
type Base struct {
	ID        ID             `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitzero"`
}

func init() {
	// Register custom validations
	v := appvalidator.Default()
	lo.Must0(v.RegisterValidation("register_type", validateRegisterType))
}

func validateRegisterType(fl validator.FieldLevel) bool {
	registerType, ok := fl.Field().Interface().(api.RegisterType)
	if !ok {
		return false
	}

	_, err := registerType.Value()

	return err == nil
}
