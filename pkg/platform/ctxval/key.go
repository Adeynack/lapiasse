package ctxval

import (
	"fmt"
	"reflect"
)

type contextValueKey struct {
	typ  reflect.Type
	name string
}

func keyFor[T any](name string) contextValueKey {
	return contextValueKey{
		typ:  reflect.TypeFor[T](),
		name: name,
	}
}

func (k contextValueKey) String() string {
	typePart := k.typ.String()

	if k.name == "" {
		return typePart
	}

	return fmt.Sprintf("%s(%s)", typePart, k.name)
}
