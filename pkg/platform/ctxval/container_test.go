package ctxval

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewContainer(t *testing.T) {
	t.Run("creates container", func(t *testing.T) {
		container := NewContainer(context.Background())

		require.NotNil(t, container)
		require.NotNil(t, container.Context)
		require.NotNil(t, container.dependenciesByKey)
	})
}

func TestContainerRegisterAndResolve(t *testing.T) {
	t.Run("when nothing is registered", func(t *testing.T) {
		container := NewContainer(t.Context())

		t.Run("Resolve fails", func(t *testing.T) {
			_, err := Resolve[Foo](container)
			require.ErrorIs(t, err, ErrUnregisteredDependency)
			require.Equal(t, `unable to resolve dependency: unregistered dependency "ctxval.Foo"`, err.Error())
		})

		t.Run("ResolveNamed fails", func(t *testing.T) {
			_, err := ResolveNamed[Foo](container, "asdf")
			require.ErrorIs(t, err, ErrUnregisteredDependency)
			require.Equal(t, `unable to resolve dependency: unregistered dependency "ctxval.Foo(asdf)"`, err.Error())
		})

		t.Run("MustResolve panics", func(t *testing.T) {
			require.PanicsWithError(
				t,
				`unable to resolve dependency: unregistered dependency "ctxval.Foo"`,
				func() { MustResolve[Foo](container) },
			)
		})

		t.Run("MustResolveNamed panics", func(t *testing.T) {
			require.PanicsWithError(
				t,
				`unable to resolve dependency: unregistered dependency "ctxval.Foo(asdf)"`,
				func() { MustResolveNamed[Foo](container, "asdf") },
			)
		})
	})

	t.Run("when an unnamed Foo value is registered", func(t *testing.T) {
		registeredFoo := Foo{Bar: "5160b303-f563-44c3-ac93-baebea18cbe7"}
		ctx := NewContainer(t.Context())
		RegisterInContainer(ctx, registeredFoo)

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
		ctx := NewContainer(t.Context())
		RegisterNamedInContainer(ctx, "TheFuu", registeredFoo)

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

func TestContainerFallback(t *testing.T) {
	t.Run("when the parent context has a value", func(t *testing.T) {
		parentCtx := t.Context()
		parentCtx = Register(parentCtx, "value from parent context")
		container := NewContainer(parentCtx)
		var ctx context.Context = container

		t.Run("when the container does not have the same value registered", func(t *testing.T) {
			t.Run("Value falls back to parent context", func(t *testing.T) {
				value, err := Resolve[string](ctx)
				require.NoError(t, err)
				require.Equal(t, "value from parent context", value)
			})
		})

		t.Run("when the container has the same value registered", func(t *testing.T) {
			RegisterInContainer(container, "value from container")

			t.Run("Value returns the container's value", func(t *testing.T) {
				value, err := Resolve[string](ctx)
				require.NoError(t, err)
				require.Equal(t, "value from container", value)
			})

			t.Run("when a normal register is happening on top of the container", func(t *testing.T) {
				ctx := Register(container, "value from child context")

				t.Run("Value returns the overriding value", func(t *testing.T) {
					value, err := Resolve[string](ctx)
					require.NoError(t, err)
					require.Equal(t, "value from child context", value)
				})
			})
		})
	})
}

// This benchmark compares the performance of dependency registration and resolution
// using the standard context-based approach versus the Container-based approach.
//
// It shows that the Container is slower to register values, but significantly faster to resolve them.
// It also, in both cases, uses less memory than the standard context approach.
//
// To run the benchmark, use the following command:
//
//	GOEXPERIMENT=jsonv2,synctest go test -v -benchmem -run=^$ -bench ^BenchmarkContainer$ adeynack.net/lapiasse/pkg/platform/ctxval
func BenchmarkContainer(b *testing.B) {
	b.ReportAllocs()
	type structWithFloat struct{ C float64 }

	registerValues := func(bn int, ctx context.Context, container *Container) context.Context {
		for i := range bn + 1 {
			name := fmt.Sprintf("value-%d", i)

			switch i % 3 {
			case 0:
				if container == nil {
					ctx = RegisterNamed(ctx, name, i)
				} else {
					RegisterNamedInContainer(container, name, i)
				}
			case 1:
				if container == nil {
					ctx = RegisterNamed(ctx, name, strconv.Itoa(i))
				} else {
					RegisterNamedInContainer(container, name, strconv.Itoa(i))
				}
			case 2:
				if container == nil {
					ctx = RegisterNamed(ctx, name, structWithFloat{C: float64(i)})
				} else {
					RegisterNamedInContainer(container, name, structWithFloat{C: float64(i)})
				}
			}
		}

		if container == nil {
			return ctx
		} else {
			return container
		}
	}

	b.Run("register values", func(b *testing.B) {
		b.Run("with basic context approach", func(b *testing.B) {
			ctx := context.Background()
			b.ResetTimer() // Don't want to creation of the context itself to be included.

			registerValues(b.N, ctx, nil)
		})

		b.Run("with container", func(b *testing.B) {
			container := NewContainer(context.Background())
			b.ResetTimer() // Don't want to creation of the container itself to be included.

			registerValues(b.N, nil, container)
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

		b.Run("with container", func(b *testing.B) {
			containerCtx := registerValues(b.N, nil, NewContainer(context.Background()))
			b.ResetTimer()

			resolveValues(b.N, containerCtx)
		})
	})
}
