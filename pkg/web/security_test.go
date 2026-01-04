package web

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateSessionToken(t *testing.T) {
	t.Run("it provides a non empty string", func(t *testing.T) {
		token := generateSessionToken()
		require.NotEmpty(t, token)
	})

	t.Run("it never provides twice the same", func(t *testing.T) {
		const iterations = 100
		tokens := make([]string, iterations)
		for i := range iterations {
			tokens[i] = generateSessionToken()
			for j := range i - 1 {
				require.NotEqual(t, tokens[j], tokens[i], "generated token is not unique")
			}
		}
	})
}
