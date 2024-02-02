package errors

import "errors"

// UnwrapAll returns the first error that is not wrapped by another error.
func UnwrapAll(err error) error {
	for {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			return err
		}
		err = unwrapped
	}
}
