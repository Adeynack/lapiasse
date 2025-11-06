package loex

import "errors"

type Erroreable = func() error

func GetAllOrErr(f ...Erroreable) error {
	for _, oneF := range f {
		if err := oneF(); err != nil {
			return err
		}
	}

	return nil
}

func GetAllOrErr2[T1, T2 any](
	f1 func() (T1, error),
	f2 func() (T2, error),
) (T1, T2, error) {
	v1, err := f1()
	if err != nil {
		var zeroT2 T2
		return v1, zeroT2, err
	}

	v2, err := f2()
	if err != nil {
		return v1, v2, err
	}

	return v1, v2, nil
}

// OnErrJoin calls the provided function and, if it returns an error,
// joins it to the error referenced by errRef.
//
// Useful for, e.g., deferring cleanup functions without ignoring errors.
//
// Example:
//
//	func someFunction() (err error) {
//	    resource, err := acquireResource()
//	    if err != nil {
//	        return err
//	    }
//	    defer loex.OnErrJoin(&err, resource.Close)
//
//	    // do stuff with resource
//
//	    return nil
//	}
func OnErrJoin(errRef *error, fn func() error) {
	if err := fn(); err != nil {
		*errRef = errors.Join(*errRef, err)
	}
}
