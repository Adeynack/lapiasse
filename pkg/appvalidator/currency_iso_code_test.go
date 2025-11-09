package appvalidator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCurrencyIsoCode(t *testing.T) {
	for name, tc := range map[string]struct {
		code      string
		expectErr bool
	}{
		"valid code": {
			code:      "CAD",
			expectErr: false,
		},
		"lowercase code": {
			code:      "usd",
			expectErr: true,
		},
		"too short code": {
			code:      "US",
			expectErr: true,
		},
		"too long code": {
			code:      "USDA",
			expectErr: true,
		},
		"numeric code": {
			code:      "123",
			expectErr: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			err := Default().Var(tc.code, "currencyIsoCode")
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
