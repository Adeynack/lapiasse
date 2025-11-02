package ctxval

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewResolver(t *testing.T) {
	t.Run("creates resolver with provided context", func(t *testing.T) {
		ctx := context.Background()
		resolver := NewResolver(ctx)

		require.NotNil(t, resolver)
		require.NotNil(t, resolver.ctx)
		require.NotNil(t, resolver.dependenciesByKey)
		require.Equal(t, 0, len(resolver.dependenciesByKey))
	})

	t.Run("creates resolver with Background context when nil is provided", func(t *testing.T) {
		resolver := NewResolver(nil) //nolint:staticcheck // Testing nil context handling

		require.NotNil(t, resolver)
		require.NotNil(t, resolver.ctx)
		require.NotNil(t, resolver.dependenciesByKey)
	})

	t.Run("creates resolver with context that has a deadline", func(t *testing.T) {
		deadline := time.Now().Add(time.Hour)
		ctx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		resolver := NewResolver(ctx)

		require.NotNil(t, resolver)
		actualDeadline, ok := resolver.Deadline()
		require.True(t, ok)
		require.Equal(t, deadline, actualDeadline)
	})
}

func TestResolver_Deadline(t *testing.T) {
	t.Run("returns deadline from underlying context", func(t *testing.T) {
		deadline := time.Now().Add(time.Hour)
		ctx, cancel := context.WithDeadline(context.Background(), deadline)
		defer cancel()

		resolver := NewResolver(ctx)
		actualDeadline, ok := resolver.Deadline()

		require.True(t, ok)
		require.Equal(t, deadline, actualDeadline)
	})

	t.Run("returns false when context has no deadline", func(t *testing.T) {
		ctx := context.Background()
		resolver := NewResolver(ctx)

		_, ok := resolver.Deadline()
		require.False(t, ok)
	})
}

func TestResolver_Done(t *testing.T) {
	t.Run("returns nil channel for background context", func(t *testing.T) {
		ctx := context.Background()
		resolver := NewResolver(ctx)

		done := resolver.Done()
		require.Nil(t, done)
	})

	t.Run("returns channel that closes when context is canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		resolver := NewResolver(ctx)

		done := resolver.Done()
		require.NotNil(t, done)

		select {
		case <-done:
			t.Fatal("channel should not be closed yet")
		default:
			// Expected
		}

		cancel()

		select {
		case <-done:
			// Expected
		case <-time.After(100 * time.Millisecond):
			t.Fatal("channel should be closed after cancel")
		}
	})
}

func TestResolver_Err(t *testing.T) {
	t.Run("returns nil for non-canceled context", func(t *testing.T) {
		ctx := context.Background()
		resolver := NewResolver(ctx)

		err := resolver.Err()
		require.NoError(t, err)
	})

	t.Run("returns Canceled error when context is canceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		resolver := NewResolver(ctx)

		cancel()

		err := resolver.Err()
		require.ErrorIs(t, err, context.Canceled)
	})

	t.Run("returns DeadlineExceeded error when deadline expires", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()
		resolver := NewResolver(ctx)

		time.Sleep(10 * time.Millisecond)

		err := resolver.Err()
		require.ErrorIs(t, err, context.DeadlineExceeded)
	})
}

func TestResolver_Value(t *testing.T) {
	type TestStruct struct {
		Value string
	}

	t.Run("returns nil for non-existent key", func(t *testing.T) {
		resolver := NewResolver(context.Background())

		value := resolver.Value(keyFor[TestStruct](""))
		require.Nil(t, value)
	})

	t.Run("returns registered unnamed value", func(t *testing.T) {
		resolver := NewResolver(context.Background())
		expected := TestStruct{Value: "test-value"}

		RegisterInResolver(resolver, expected)

		value := resolver.Value(keyFor[TestStruct](""))
		require.NotNil(t, value)
		require.Equal(t, expected, value)
	})

	t.Run("returns registered named value", func(t *testing.T) {
		resolver := NewResolver(context.Background())
		expected := TestStruct{Value: "named-value"}

		RegisterNamedInResolver(resolver, "myName", expected)

		value := resolver.Value(keyFor[TestStruct]("myName"))
		require.NotNil(t, value)
		require.Equal(t, expected, value)
	})

	t.Run("distinguishes between unnamed and named values", func(t *testing.T) {
		resolver := NewResolver(context.Background())
		unnamedValue := TestStruct{Value: "unnamed"}
		namedValue := TestStruct{Value: "named"}

		RegisterInResolver(resolver, unnamedValue)
		RegisterNamedInResolver(resolver, "special", namedValue)

		// Get unnamed value
		value := resolver.Value(keyFor[TestStruct](""))
		require.NotNil(t, value)
		require.Equal(t, unnamedValue, value)

		// Get named value
		value = resolver.Value(keyFor[TestStruct]("special"))
		require.NotNil(t, value)
		require.Equal(t, namedValue, value)

		// Non-existent name returns nil
		value = resolver.Value(keyFor[TestStruct]("other"))
		require.Nil(t, value)
	})

	t.Run("falls back to underlying context for non-contextValueKey types", func(t *testing.T) {
		type customKey string
		key := customKey("custom")
		expectedValue := "context-value"

		ctx := context.WithValue(context.Background(), key, expectedValue)
		resolver := NewResolver(ctx)

		value := resolver.Value(key)
		require.Equal(t, expectedValue, value)
	})

	t.Run("resolver values take precedence over context values for contextValueKey", func(t *testing.T) {
		key := keyFor[TestStruct]("")
		contextValue := TestStruct{Value: "from-context"}
		resolverValue := TestStruct{Value: "from-resolver"}

		ctx := context.WithValue(context.Background(), key, contextValue)
		resolver := NewResolver(ctx)
		RegisterInResolver(resolver, resolverValue)

		value := resolver.Value(key)
		require.Equal(t, resolverValue, value)
	})
}

func TestRegisterInResolver(t *testing.T) {
	type CustomType struct {
		Data string
	}

	t.Run("registers value that can be retrieved", func(t *testing.T) {
		resolver := NewResolver(context.Background())
		expected := CustomType{Data: "test-data"}

		RegisterInResolver(resolver, expected)

		key := keyFor[CustomType]("")
		value := resolver.Value(key)
		require.NotNil(t, value)
		require.Equal(t, expected, value)
	})

	t.Run("overwrites previously registered value", func(t *testing.T) {
		resolver := NewResolver(context.Background())
		first := CustomType{Data: "first"}
		second := CustomType{Data: "second"}

		RegisterInResolver(resolver, first)
		RegisterInResolver(resolver, second)

		key := keyFor[CustomType]("")
		value := resolver.Value(key)
		require.NotNil(t, value)
		require.Equal(t, second, value)
	})

	t.Run("works with multiple different types", func(t *testing.T) {
		type TypeA struct{ Value string }
		type TypeB struct{ Value int }

		resolver := NewResolver(context.Background())
		valueA := TypeA{Value: "a-value"}
		valueB := TypeB{Value: 42}

		RegisterInResolver(resolver, valueA)
		RegisterInResolver(resolver, valueB)

		retrievedA := resolver.Value(keyFor[TypeA](""))
		require.Equal(t, valueA, retrievedA)

		retrievedB := resolver.Value(keyFor[TypeB](""))
		require.Equal(t, valueB, retrievedB)
	})
}

func TestRegisterNamedInResolver(t *testing.T) {
	type CustomType struct {
		Data string
	}

	t.Run("registers named value that can be retrieved", func(t *testing.T) {
		resolver := NewResolver(context.Background())
		expected := CustomType{Data: "named-data"}

		RegisterNamedInResolver(resolver, "myName", expected)

		key := keyFor[CustomType]("myName")
		value := resolver.Value(key)
		require.NotNil(t, value)
		require.Equal(t, expected, value)
	})

	t.Run("allows multiple named instances of same type", func(t *testing.T) {
		resolver := NewResolver(context.Background())
		primary := CustomType{Data: "primary"}
		secondary := CustomType{Data: "secondary"}

		RegisterNamedInResolver(resolver, "primary", primary)
		RegisterNamedInResolver(resolver, "secondary", secondary)

		primaryKey := keyFor[CustomType]("primary")
		retrievedPrimary := resolver.Value(primaryKey)
		require.Equal(t, primary, retrievedPrimary)

		secondaryKey := keyFor[CustomType]("secondary")
		retrievedSecondary := resolver.Value(secondaryKey)
		require.Equal(t, secondary, retrievedSecondary)
	})

	t.Run("named and unnamed registrations are independent", func(t *testing.T) {
		resolver := NewResolver(context.Background())
		unnamed := CustomType{Data: "unnamed"}
		named := CustomType{Data: "named"}

		RegisterInResolver(resolver, unnamed)
		RegisterNamedInResolver(resolver, "special", named)

		unnamedKey := keyFor[CustomType]("")
		retrievedUnnamed := resolver.Value(unnamedKey)
		require.Equal(t, unnamed, retrievedUnnamed)

		namedKey := keyFor[CustomType]("special")
		retrievedNamed := resolver.Value(namedKey)
		require.Equal(t, named, retrievedNamed)
	})

	t.Run("overwrites previously registered named value with same name", func(t *testing.T) {
		resolver := NewResolver(context.Background())
		first := CustomType{Data: "first"}
		second := CustomType{Data: "second"}

		RegisterNamedInResolver(resolver, "myName", first)
		RegisterNamedInResolver(resolver, "myName", second)

		key := keyFor[CustomType]("myName")
		value := resolver.Value(key)
		require.Equal(t, second, value)
	})
}

func TestResolver_ImplementsContextInterface(t *testing.T) {
	t.Run("Resolver implements context.Context interface", func(t *testing.T) {
		resolver := NewResolver(context.Background())

		// This should compile if Resolver implements context.Context
		var _ context.Context = resolver
		require.NotNil(t, resolver)
	})

	t.Run("can be used as context in functions expecting context.Context", func(t *testing.T) {
		resolver := NewResolver(context.Background())

		// Function that accepts context.Context
		acceptsContext := func(ctx context.Context) bool {
			return ctx != nil
		}

		require.True(t, acceptsContext(resolver))
	})
}
