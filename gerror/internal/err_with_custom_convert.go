package internal

import (
	"fmt"

	"github.com/drshriveer/gtools/gerror"
)

//go:generate gerror --types ErrWithCustomConvert --skipConvertGen

// ErrWithCustomConvert is just a test.
//
//nolint:errname // dumb
type ErrWithCustomConvert struct {
	gerror.GError
	Property Status `gerror:"_,print,clone"`
}

// Convert will attempt to convert the supplied error into a gError.Error of the
// Factory's type, including the source errors details in the result's error message.
// The original error's equality can be checked with errors.Is().
func (e *ErrWithCustomConvert) Convert(err error) gerror.Error {
	if gerr, ok := err.(gerror.Error); ok {
		return gerr
	}
	clone := gerror.CloneBase(e, gerror.SourceStack, "", "", fmt.Sprintf("originalError: %+v", err), err)
	return e.toPrimaryType(clone)
}

// ConvertS is the same as Convert but includes a full StackTrace.
func (e *ErrWithCustomConvert) ConvertS(err error) gerror.Error {
	if gerr, ok := err.(gerror.Error); ok {
		return gerr
	}
	clone := gerror.CloneBase(e, gerror.DefaultStack, "", "", fmt.Sprintf("originalError: %+v", err), err)
	return e.toPrimaryType(clone)
}
