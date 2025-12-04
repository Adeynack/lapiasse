//go:build test

package model_test

import (
	"testing"

	"adeynack.net/lapiasse/pkg/appvalidator"
	"adeynack.net/lapiasse/pkg/model"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/require"
)

func TestRegisterCreate(t *testing.T) {
	// ctx := app.CreateTestAppCtx(t)
	// db := ctxval.MustResolve[*gorm.DB](ctx)

	register := &model.Register{
		Name: "", // Name is required
	}

	err := appvalidator.Default().Struct(register)
	requireValidationError(t, err, "Name", "required")
	requireValidationError(t, err, "Name", "foo")
}

func requireValidationError(t *testing.T, err error, field string, tag string) {
	var validationErr validator.ValidationErrors
	require.ErrorAs(t, err, &validationErr)

	for _, fe := range validationErr {
		if fe.Field() == field && fe.Tag() == tag {
			return
		}
	}

	require.Failf(t, "Expected validation error not found", "field: %s, tag: %s\nerrors:\n%s", field, tag, validationErr)
}
