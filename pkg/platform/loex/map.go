package loex

import "github.com/samber/lo"

// MapE is like Map but the iteratee function does not receive the index.
// Allows for passing simpler functions as arguments (e.g.: [strings.ToUpper])
// directly, without needing to wrap them in yet another function.
func MapE[T any, R any](collection []T, iteratee func(item T) R) []R {
	// NOTE: if it ever gets imported to `lo`, duplicate the implementation
	// instead of calling [lo.Map], to avoid unnecessary call stacks.
	return lo.Map(collection, func(item T, _ int) R {
		return iteratee(item)
	})
}
