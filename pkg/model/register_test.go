//go:build test

package model_test

// import (
// 	"testing"
// 	"time"

// 	"adeynack.net/lapiasse/pkg/api"
// 	"adeynack.net/lapiasse/pkg/appvalidator"
// 	"adeynack.net/lapiasse/pkg/model"
// 	"adeynack.net/lapiasse/pkg/platform/requireex"
// 	"adeynack.net/lapiasse/pkg/platform/validatorex"
// 	"github.com/samber/lo"
// )

// func TestRegisterValidate(t *testing.T) {
// 	for name, tc := range map[string]struct {
// 		register                 model.Register
// 		expectedValidationErrors map[string][]string
// 	}{
// 		"zero value": {
// 			register: model.Register{},
// 			expectedValidationErrors: map[string][]string{
// 				"Name":              {"required"},
// 				"CurrencyIsoCode":   {"required"},
// 				"BookID":            {"required"},
// 				"StartsAt":          {"required"},
// 				"DefaultCategoryID": {"required"},
// 				"Type":              {"required"},
// 				"InitialBalance":    {"required"},
// 			},
// 		}} {
// 		t.Run(name, func(t *testing.T) {
// 			err := appvalidator.Default().Struct(tc.register)
// 			validatorex.RequireValidationErrorTags(t, err, tc.expectedValidationErrors)
// 		})
// 	}
// }

// func TestRegisterJSONMarshall(t *testing.T) {
// 	register := model.Register{
// 		Base: model.Base{
// 			ID: 1, // should appear as string in JSON
// 		},
// 		Name:            "My Register",
// 		Type:            api.RegisterTypeBank,
// 		BookID:          9, // should appear as string in JSON
// 		StartsAt:        lo.Must1(time.Parse(time.DateOnly, "2025-05-24")),
// 		ExpiresAt:       nil, // should not appear in JSON
// 		CurrencyIsoCode: "EUR",
// 		// Notes:          "", // should not appear in JSON
// 		// InitialBalance: 0, // should default to 0 and appear in JSON
// 		Active:            true,
// 		DefaultCategoryID: 1234, // should appear as string in JSON
// 	}

// 	requireex.NoJsonDiffFromStruct(t, register, `{
// 		"id": "1",
// 		"created_at": "0001-01-01T00:00:00Z",
// 		"updated_at": "0001-01-01T00:00:00Z",
// 		"name": "My Register",
// 		"type": "bank",
// 		"book_id": "9",
// 		"starts_at": "2025-05-24",
// 		"parent_id": null,
// 		"currency_iso_code": "EUR",
// 		"initial_balance": 0,
// 		"active": true,
// 		"default_category_id": "1234"
// 	}`)
// }
