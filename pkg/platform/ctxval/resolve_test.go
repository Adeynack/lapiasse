package ctxval

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Foo struct {
	Bar string
}

type Bar struct {
	Foo string
}

func TestRegisterAndResolve(t *testing.T) {
	t.Run("when nothing is registered", func(t *testing.T) {
		t.Run("Resolve fails", func(t *testing.T) {
			_, err := Resolve[Foo](t.Context())
			require.ErrorIs(t, err, ErrUnregisteredDependency)
			require.Equal(t, `unable to resolve dependency: unregistered dependency "ctxval.Foo"`, err.Error())
		})

		t.Run("ResolveNamed fails", func(t *testing.T) {
			_, err := ResolveNamed[Foo](t.Context(), "asdf")
			require.ErrorIs(t, err, ErrUnregisteredDependency)
			require.Equal(t, `unable to resolve dependency: unregistered dependency "ctxval.Foo(asdf)"`, err.Error())
		})

		t.Run("MustResolve panics", func(t *testing.T) {
			require.PanicsWithError(
				t,
				`unable to resolve dependency: unregistered dependency "ctxval.Foo"`,
				func() { MustResolve[Foo](t.Context()) },
			)
		})

		t.Run("MustResolveNamed panics", func(t *testing.T) {
			require.PanicsWithError(
				t,
				`unable to resolve dependency: unregistered dependency "ctxval.Foo(asdf)"`,
				func() { MustResolveNamed[Foo](t.Context(), "asdf") },
			)
		})
	})

	t.Run("when an unnamed Foo value is registered", func(t *testing.T) {
		registeredFoo := Foo{Bar: "5160b303-f563-44c3-ac93-baebea18cbe7"}
		ctx := Register(t.Context(), registeredFoo)

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
		ctx := RegisterNamed(t.Context(), "TheFuu", registeredFoo)

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

func TestResolve2(t *testing.T) {
	t.Run("when nothing is registered", func(t *testing.T) {
		t.Run("Resolve2 fails with the 1st type's error", func(t *testing.T) {
			_, _, err := Resolve2[Foo, Bar](t.Context())
			require.ErrorIs(t, err, ErrUnregisteredDependency)
			require.Equal(t, `unable to resolve dependency: unregistered dependency "ctxval.Foo"`, err.Error())
		})
	})

	t.Run("when the 1st type is registered", func(t *testing.T) {
		ctx := t.Context()
		registeredFoo := Foo{Bar: "f2074b00-8679-4a33-951b-167934dd707b"}
		ctx = Register(ctx, registeredFoo)

		t.Run("Resolve2 fails with the 2nd type's error", func(t *testing.T) {
			_, _, err := Resolve2[Foo, Bar](ctx)
			require.ErrorIs(t, err, ErrUnregisteredDependency)
			require.Equal(t, `unable to resolve dependency: unregistered dependency "ctxval.Bar"`, err.Error())
		})
	})

	t.Run("when the 2nd type is registered", func(t *testing.T) {
		ctx := t.Context()
		registeredFoo := Bar{Foo: "3288d2ac-077f-45e8-8981-674ee1853645"}
		ctx = Register(ctx, registeredFoo)

		t.Run("Resolve2 fails with the 1st type's error", func(t *testing.T) {
			_, _, err := Resolve2[Foo, Bar](ctx)
			require.ErrorIs(t, err, ErrUnregisteredDependency)
			require.Equal(t, `unable to resolve dependency: unregistered dependency "ctxval.Foo"`, err.Error())
		})
	})

	t.Run("when all types are registered", func(t *testing.T) {
		ctx := t.Context()
		registeredFoo := Foo{Bar: "50dcc406-6932-427e-af98-82bfc011dd9e"}
		registeredBar := Bar{Foo: "e8f22c3e-37a5-4ad2-8747-f5c7b4b288f7"}
		ctx = Register(ctx, registeredFoo)
		ctx = Register(ctx, registeredBar)

		t.Run("Resolve2 fails with the 1st type's error", func(t *testing.T) {
			foo, bar, err := Resolve2[Foo, Bar](ctx)
			require.NoError(t, err)
			require.Equal(t, registeredFoo.Bar, foo.Bar)
			require.Equal(t, registeredBar.Foo, bar.Foo)
		})
	})
}

func TestResolve3(t *testing.T) {
	t.Run("when all types are registered (smoke test)", func(t *testing.T) {
		ctx := t.Context()
		registeredString := "f3062b0d-97df-4436-a579-158f3175d62d"
		ctx = Register(ctx, registeredString)
		registeredFoo := Foo{Bar: "3e083275-7815-4b4c-bad7-ef59b578513a"}
		ctx = Register(ctx, registeredFoo)
		registeredBar := Bar{Foo: "80d9bcf2-746e-4026-aec5-f734f22fbe00"}
		ctx = Register(ctx, registeredBar)

		t.Run("Register3 succeeds with all types", func(t *testing.T) {
			s, foo, bar, err := Resolve3[string, Foo, Bar](ctx)
			require.NoError(t, err)
			require.Equal(t, registeredString, s)
			require.Equal(t, registeredFoo.Bar, foo.Bar)
			require.Equal(t, registeredBar.Foo, bar.Foo)
		})
	})
}

func TestResolve4(t *testing.T) {
	t.Run("when all types are registered (smoke test)", func(t *testing.T) {
		ctx := t.Context()
		registeredString := "f3062b0d-97df-4436-a579-158f3175d62d"
		ctx = Register(ctx, registeredString)
		registeredByte := byte(192)
		ctx = Register(ctx, registeredByte)
		registeredFoo := Foo{Bar: "3e083275-7815-4b4c-bad7-ef59b578513a"}
		ctx = Register(ctx, registeredFoo)
		registeredBar := Bar{Foo: "80d9bcf2-746e-4026-aec5-f734f22fbe00"}
		ctx = Register(ctx, registeredBar)

		t.Run("Register4 succeeds with all types", func(t *testing.T) {
			s, b, foo, bar, err := Resolve4[string, byte, Foo, Bar](ctx)
			require.NoError(t, err)
			require.Equal(t, registeredString, s)
			require.Equal(t, registeredByte, b)
			require.Equal(t, registeredFoo.Bar, foo.Bar)
			require.Equal(t, registeredBar.Foo, bar.Foo)
		})
	})
}
