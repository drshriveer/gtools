package gerrors_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/drshriveer/gtools/gerrors"
)

// What are the actual options here?!
// 1. move all the methods out of interfaces like pkg errors.
//    use reflect to properly update the values (example bleow)
//    this is dangrous for anyhing with a pointer / slice / ec
// 2. Is there an option for wrapping a type in a factory?
// 3. Annotate a type and code gen // but do we need tags for the same cloning problem reflect has?
// 4. just copy this shit over and do what we want.

// .. ultimately i think there are two problems here:
// 1. how to clone the parent error. A FactoryOf-type pattern should work for this, when the result is an Error type.
// 2. filed tags can be used to do the actual cloning as well.

type ExtendedError struct {
	gerrors.GError
	GRPCStatus  Status `gerror:"clone,print"`
	SomeMessage string `gerror:"clone,print"`
	DoNotPrint  string `gerror:"clone"`
}

var ErrExtendedExample = gerrors.FactoryOf(&ExtendedError{
	GError: gerrors.GError{
		Name:    "ErrExtendedExample",
		Message: "extended error example",
	},
	GRPCStatus:  InvalidArgument,
	SomeMessage: "Print this message",
	DoNotPrint:  "this is for internal issue only",
})

func TestExtendedError_Equality(t *testing.T) {
	err1 := L1()
	err2 := L1()
	assert.True(t, errors.Is(ErrExtendedExample, ErrExtendedExample))

	_, ok := err1.(interface{ Is(error) bool })
	require.True(t, ok)

	assert.True(t, err1.(gerrors.Error).Is(err2))
	assert.True(t, err2.(gerrors.Error).Is(err1))
	assert.True(t, errors.Is(err1, err2))
	assert.True(t, errors.Is(err2, err1))
	assert.True(t, errors.Is(err1, ErrExtendedExample))
	assert.True(t, errors.Is(ErrExtendedExample, err1))

	errToConvert := errors.New("random error")
	convertedErr := ErrExtendedExample.Convert(errToConvert)
	assert.True(t, errors.Is(convertedErr, ErrExtendedExample))
	assert.True(t, errors.Is(convertedErr, errToConvert))
	assert.True(t, errors.Is(convertedErr, errToConvert))

	switch errors.Unwrap(err2) {
	case ErrExtendedExample:
	default:
		assert.Fail(t, "was supposed to reach case above")
	}

}

func TestExtendedError_CorrectlyLogged(t *testing.T) {
	err, ok := L1().(gerrors.Error)
	require.Truef(t, ok, "error must implement the gerror.Error interface")
	assert.Contains(t, err.Error(), "GRPCStatus: InvalidArgument, ")
	assert.Contains(t, err.Error(), "SomeMessage: Print this message, ")
	assert.NotContains(t, err.Error(), "DoNotPrint")
	assert.NotContains(t, err.Error(), "this is for internal issue only")
	assert.Equal(t, "ErrExtendedExample", err.ErrName())
	assert.Equal(t, "extended error example", err.ErrMessage())
	assert.Equal(t, "gerrors_test:L3", err.ErrSource())
	assert.Equal(t, "", err.DTag())
}

func L1() error {
	return L2()
}
func L2() error {
	return L3()
}
func L3() error {
	return ErrExtendedExample.WithStack()
}

// Status is just a sample for testing.
type Status int

const (
	OK Status = iota
	Canceled
	Unknown
	InvalidArgument
	DeadlineExceeded
)

func (s Status) String() string {
	switch s {
	case OK:
		return "OK"
	case Canceled:
		return "Canceled"
	case Unknown:
		return "Unknown"
	case InvalidArgument:
		return "InvalidArgument"
	case DeadlineExceeded:
		return "DeadlineExceeded"
	}
	return "UNKNOWN"
}
