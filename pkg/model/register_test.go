//go:build test

package model_test

import (
	"testing"

	"adeynack.net/lapiasse/pkg/appvalidator"
	"adeynack.net/lapiasse/pkg/model"
	"adeynack.net/lapiasse/pkg/platform/validatorex"
)

func TestValidateCreate(t *testing.T) {
	for name, tc := range map[string]struct {
		register                 model.Register
		expectedValidationErrors map[string][]string
	}{
		"zero value": {
			register: model.Register{},
			expectedValidationErrors: map[string][]string{
				"Name":              {"required"},
				"CurrencyIsoCode":   {"required"},
				"BookID":            {"required"},
				"StartsAt":          {"required"},
				"DefaultCategoryID": {"required"},
				"Type":              {"required"},
				"InitialBalance":    {"required"},
			},
		}} {
		t.Run(name, func(t *testing.T) {
			err := appvalidator.Default().Struct(tc.register)
			validatorex.RequireValidationErrorTags(t, err, tc.expectedValidationErrors)
		})
	}
}
