package controller

import (
	"testing"

	"adeynack.net/lapiasse/pkg/api"
)

func TestService(t *testing.T) {
	t.Run("Service implements the API interface", func(t *testing.T) {
		var _ api.StrictServerInterface = (*Service)(nil)
	})
}
