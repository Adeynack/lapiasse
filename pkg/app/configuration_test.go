package app

import (
	"fmt"
	"testing"

	"adeynack.net/lapiasse/pkg/env"
	"adeynack.net/lapiasse/pkg/platform/ctxval"
	"github.com/stretchr/testify/require"
)

func TestConfiguration(t *testing.T) {
	for _, tc := range []struct {
		env         env.Environment
		expectPanic bool
	}{
		{env: env.EnvDevelopment},
		{env: env.EnvTest, expectPanic: true},
		{env: env.EnvProduction},
	} {
		t.Run(fmt.Sprintf("in environment %s", tc.env.String()), func(t *testing.T) {
			ctx := ctxval.RegisterNamed(t.Context(), "run", tc.env)

			if tc.expectPanic {
				require.Panics(t, func() {
					_, _ = InitializeConfiguration(ctx, CliFlags{})
				})

				return
			}

			h, err := InitializeConfiguration(ctx, CliFlags{})
			require.NoError(t, err)

			data := h.Configuration.Data
			require.NotEmpty(t, data.BasePath)
		})
	}
}
