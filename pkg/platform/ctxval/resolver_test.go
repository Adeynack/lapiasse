package ctxval

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewResolver(t *testing.T) {
	t.Run("creates resolver", func(t *testing.T) {
		resolver := NewResolver(context.Background())

		require.NotNil(t, resolver)
		require.NotNil(t, resolver.Context)
		require.NotNil(t, resolver.dependenciesByKey)
	})
}

func TestResolverRegisterResolve(t *testing.T) {
	t.Run("when nothing is registered", func(t *testing.T) {
		resolver := NewResolver(t.Context())

		t.Run("Resolve fails", func(t *testing.T) {
			_, err := Resolve[Foo](resolver)
			require.ErrorIs(t, err, ErrUnregisteredDependency)
			require.Equal(t, `unable to resolve dependency: unregistered dependency "ctxval.Foo"`, err.Error())
		})

		t.Run("ResolveNamed fails", func(t *testing.T) {
			_, err := ResolveNamed[Foo](resolver, "asdf")
			require.ErrorIs(t, err, ErrUnregisteredDependency)
			require.Equal(t, `unable to resolve dependency: unregistered dependency "ctxval.Foo(asdf)"`, err.Error())
		})

		t.Run("MustResolve panics", func(t *testing.T) {
			require.PanicsWithError(
				t,
				`unable to resolve dependency: unregistered dependency "ctxval.Foo"`,
				func() { MustResolve[Foo](resolver) },
			)
		})

		t.Run("MustResolveNamed panics", func(t *testing.T) {
			require.PanicsWithError(
				t,
				`unable to resolve dependency: unregistered dependency "ctxval.Foo(asdf)"`,
				func() { MustResolveNamed[Foo](resolver, "asdf") },
			)
		})
	})

	t.Run("when an unnamed Foo value is registered", func(t *testing.T) {
		registeredFoo := Foo{Bar: "5160b303-f563-44c3-ac93-baebea18cbe7"}
		ctx := NewResolver(t.Context())
		RegisterInResolver(ctx, registeredFoo)

		t.Run("Resolve succeeds for unnamed Foo", func(t *testing.T) {
			result, err := Resolve[Foo](ctx)
			require.NoError(t, err)
			require.Equal(t, registeredFoo, result)
		})

		t.Run("Resolve fails for unnamed Bar", func(t *testing.T) {
			_, err := Resolve[Bar](ctx)
			require.ErrorIs(t, err, ErrUnregisteredDependency)
			require.Equal(t, `unable to resolve dependency: unregistered dependency "ctxval.Bar"`, err.Error())
		})

		t.Run("ResolveNamed fails for named Foo", func(t *testing.T) {
			_, err := ResolveNamed[Foo](ctx, "asdf")
			require.ErrorIs(t, err, ErrUnregisteredDependency)
			require.Equal(t, `unable to resolve dependency: unregistered dependency "ctxval.Foo(asdf)"`, err.Error())
		})
	})

	t.Run("when a named Foo value is registered", func(t *testing.T) {
		registeredFoo := Foo{Bar: "e1950227-441b-4238-804f-908110c0592a"}
		ctx := NewResolver(t.Context())
		RegisterNamedInResolver(ctx, "TheFuu", registeredFoo)

		t.Run("Resolve fails for unnamed Foo", func(t *testing.T) {
			_, err := Resolve[Foo](ctx)
			require.ErrorIs(t, err, ErrUnregisteredDependency)
			require.Equal(t, `unable to resolve dependency: unregistered dependency "ctxval.Foo"`, err.Error())
		})

		t.Run("Resolve fails for unnamed Bar", func(t *testing.T) {
			_, err := Resolve[Bar](ctx)
			require.ErrorIs(t, err, ErrUnregisteredDependency)
			require.Equal(t, `unable to resolve dependency: unregistered dependency "ctxval.Bar"`, err.Error())
		})

		t.Run("Resolve succeeds for Foo named 'TheFuu'", func(t *testing.T) {
			result, err := ResolveNamed[Foo](ctx, "TheFuu")
			require.NoError(t, err)
			require.Equal(t, registeredFoo, result)
		})

		t.Run("Resolve fails for Foo named 'SomethingElse'", func(t *testing.T) {
			_, err := ResolveNamed[Foo](ctx, "SomethingElse")
			require.ErrorIs(t, err, ErrUnregisteredDependency)
			require.Equal(t, `unable to resolve dependency: unregistered dependency "ctxval.Foo(SomethingElse)"`, err.Error())
		})
	})
}

// This benchmark compares the performance of dependency registration and resolution
// using the standard context-based approach versus the Resolver-based approach.
//
// It shows that the Resolver is slower to register values, but significantly faster to resolve them.
// It also, in both cases, uses less memory than the standard context approach.
//
// To run the benchmark, use the following command:
//
//	GOEXPERIMENT=jsonv2 go test -v -benchmem -run=^$ -bench ^BenchmarkResolver$ adeynack.net/lapiasse/pkg/platform/ctxval
func BenchmarkResolver(b *testing.B) {
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
			ctx := context.Background()
			b.ResetTimer() // Don't want to creation of the context itself to be included.

			registerValues(b.N, ctx, nil)
		})

		b.Run("with resolver", func(b *testing.B) {
			resolver := NewResolver(context.Background())
			b.ResetTimer() // Don't want to creation of the resolver itself to be included.

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
