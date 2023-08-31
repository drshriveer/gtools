package errors

import (
	"fmt"
)

// Include adds more details to an error if the input error is a GError.
func Include(err error, format string, elems ...any) error {
	// FIXME! this needs to skip another for the source.
	if gerr, ok := err.(GError); ok {
		return gerr.Merge(GError{ExtMessage: fmt.Sprintf(format, elems...)})
	}
	return err
}
