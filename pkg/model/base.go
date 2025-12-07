package model

import (
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

// Base model including common fields.
type Base struct {
	ID        ID `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
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
