package errortrace

import "errors"

func Is(err error, target error) bool {
	var errTrace *Error
	if errors.As(err, &errTrace) {
		return errors.Is(errTrace.err, target)
	}

	return errors.Is(err, target)
}
