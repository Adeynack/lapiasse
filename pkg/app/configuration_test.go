package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfiguration(t *testing.T) {
	h, err := InitializeConfiguration(CliFlags{})
	require.NoError(t, err)

	t.Run("the Data configuration is properly initialized", func(t *testing.T) {
		data := h.Configuration.Data

		require.NotEmpty(t, data.BasePath)
	})
}
