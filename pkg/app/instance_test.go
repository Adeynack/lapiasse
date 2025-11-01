package app

import (
	"io"
	"testing"
)

func TestInstanceImplementations(t *testing.T) {
	t.Run("Instance implements io.Closer", func(t *testing.T) {
		var _ io.Closer = (*Instance)(nil)
	})
}
