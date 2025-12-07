package errorex

import "errors"

func CallJoinErr(err *error, fn func() error) {
	if fnErr := fn(); fnErr != nil {
		*err = errors.Join(*err, fnErr)
	}
}
