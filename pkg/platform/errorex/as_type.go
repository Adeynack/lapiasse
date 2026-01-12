package errorex

import "errors"

// AsType tries to cast the given error to the target type T using errors.As.
// It returns the casted value and true if successful, or the zero value of T
// and false otherwise.
//
// To be replaced by `errors.AsType` in Go 1.26+.
func AsType[T error](err error) (T, bool) {
	var target T
	if errors.As(err, &target) {
		return target, true
	}

	var zero T
	return zero, false
}
