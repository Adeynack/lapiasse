package mapex

import (
	"encoding/json"
	"fmt"
)

type BiDirectionalLookup[Left comparable, Right comparable] struct {
	LeftToRight map[Left]Right
	RightToLeft map[Right]Left
}

// NewBiDirectionalLookupOrPanic creates a new BiDirectionalLookup from the given source map.
// It panics if there are duplicate right values in the source map.
func NewBiDirectionalLookupOrPanic[Left comparable, Right comparable](
	source map[Left]Right,
) *BiDirectionalLookup[Left, Right] {
	b := &BiDirectionalLookup[Left, Right]{
		LeftToRight: make(map[Left]Right),
		RightToLeft: make(map[Right]Left),
	}

	for left, right := range source {
		if _, exists := b.RightToLeft[right]; exists {
			panic(fmt.Sprintf(`duplicate right value "%v" found when creating BiDirectionalLookup`, right))
		}

		b.LeftToRight[left] = right
		b.RightToLeft[right] = left
	}

	return b
}

// FromLeft returns the right value for the given left value.
func (b *BiDirectionalLookup[Left, Right]) FromLeft(left Left) (Right, bool) {
	right, ok := b.LeftToRight[left]

	return right, ok
}

// FromRight returns the left value for the given right value.
func (b *BiDirectionalLookup[Left, Right]) FromRight(right Right) (Left, bool) {
	left, ok := b.RightToLeft[right]

	return left, ok
}

// MustFromLeft returns the right value for the given left value.
// It panics if the left value is not found.
func (b *BiDirectionalLookup[Left, Right]) MustFromLeft(left Left) Right {
	right, ok := b.LeftToRight[left]
	if !ok {
		panic(fmt.Sprintf(`no right value found for left value "%v"`, left))
	}
	return right
}

// MustFromRight returns the left value for the given right value.
// It panics if the right value is not found.
func (b *BiDirectionalLookup[Left, Right]) MustFromRight(right Right) Left {
	left, ok := b.RightToLeft[right]
	if !ok {
		panic(fmt.Sprintf(`no left value found for right value "%v"`, right))
	}

	return left
}

// Len returns the number of entries in the lookup.
func (b *BiDirectionalLookup[Left, Right]) Len() int {
	return len(b.LeftToRight)
}

// MarshalJSON implements [json.Marshaler.MarshalJSON].
func (b *BiDirectionalLookup[Left, Right]) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.LeftToRight)
}
