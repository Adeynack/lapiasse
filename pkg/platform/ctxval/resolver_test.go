package ctxval

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewResolver(t *testing.T) {
	t.Run("creates resolver with Background context when nil is provided", func(t *testing.T) {
		resolver := NewResolver(nil) //nolint:staticcheck // Testing nil context handling

		require.NotNil(t, resolver)
		require.NotNil(t, resolver.ctx)
		require.NotNil(t, resolver.dependenciesByKey)
	})
}

func BenchmarkResolver(b *testing.B) {
	// GOEXPERIMENT=jsonv2 go test -v -benchmem -run=^$ -bench ^BenchmarkResolver$ adeynack.net/lapiasse/pkg/platform/ctxval
	type structWithFloat struct{ C float64 }

	registerValues := func(bn int, ctx context.Context, resolver *Resolver) context.Context {
		for i := range bn + 1 {
			name := fmt.Sprintf("value-%d", i)

			switch i % 3 {
			case 0:
				if resolver == nil {
					ctx = RegisterNamed(ctx, name, i)
				} else {
					RegisterNamedInResolver(resolver, name, i)
				}
			case 1:
				if resolver == nil {
					ctx = RegisterNamed(ctx, name, strconv.Itoa(i))
				} else {
					RegisterNamedInResolver(resolver, name, strconv.Itoa(i))
				}
			case 2:
				if resolver == nil {
					ctx = RegisterNamed(ctx, name, structWithFloat{C: float64(i)})
				} else {
					RegisterNamedInResolver(resolver, name, structWithFloat{C: float64(i)})
				}
			}
		}

		if resolver == nil {
			return ctx
		} else {
			return resolver
		}
	}

	b.Run("register values", func(b *testing.B) {
		b.Run("with basic context approach", func(b *testing.B) {
			b.Logf("b.N = %d\n", b.N)
			registerValues(b.N, context.Background(), nil)
		})

		b.Run("with resolver", func(b *testing.B) {
			b.Logf("b.N = %d\n", b.N)
			resolver := NewResolver(context.Background())
			registerValues(b.N, nil, resolver)
		})
	})

	b.Run("resolve values", func(b *testing.B) {
		resolveValues := func(bn int, ctx context.Context) {
			for i := range bn + 1 {
				name := fmt.Sprintf("value-%d", i)
				var err error

				switch i % 3 {
				case 0:
					_, err = ResolveNamed[int](ctx, name)
				case 1:
					_, err = ResolveNamed[string](ctx, name)
				case 2:
					_, err = ResolveNamed[structWithFloat](ctx, name)
				}

				require.NoError(b, err)
			}
		}

		b.Run("with basic context approach", func(b *testing.B) {
			basicCtx := registerValues(b.N, context.Background(), nil)
			b.ResetTimer()
			resolveValues(b.N, basicCtx)
		})

		b.Run("with resolver", func(b *testing.B) {
			resolverCtx := registerValues(b.N, nil, NewResolver(context.Background()))
			b.ResetTimer()
			resolveValues(b.N, resolverCtx)
		})
	})
}
