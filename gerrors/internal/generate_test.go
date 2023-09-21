package internal_test

import (
	"errors"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/drshriveer/gtools/gerrors"
	"github.com/drshriveer/gtools/gerrors/gen"
	"github.com/drshriveer/gtools/gerrors/internal"
)

func TestGenerator(t *testing.T) {
	generator := gen.Generate{
		InFile:  path.Join(os.Getenv("PWD"), "/custom_error.go"),
		OutFile: path.Join(os.Getenv("PWD"), "/custom_error.gerror.go"),
		Types:   []string{"GRPCError"},
	}

	require.NoError(t, generator.Parse())
	require.NoError(t, generator.Write())
}

func TestExtendedError_Equality(t *testing.T) {
	err1 := internal.L1()
	err2 := internal.L1()
	assert.True(t, errors.Is(internal.ErrExtendedExample, internal.ErrExtendedExample))

	_, ok := err1.(interface{ Is(error) bool })
	require.True(t, ok)

	assert.True(t, err1.(gerrors.Error).Is(err2))
	assert.True(t, err2.(gerrors.Error).Is(err1))
	assert.True(t, errors.Is(err1, err2))
	assert.True(t, errors.Is(err2, err1))
	assert.True(t, errors.Is(err1, internal.ErrExtendedExample))
	assert.True(t, errors.Is(internal.ErrExtendedExample, err1))

	errToConvert := errors.New("random error")
	convertedErr := internal.ErrExtendedExample.Convert(errToConvert)
	assert.True(t, errors.Is(convertedErr, internal.ErrExtendedExample))
	assert.True(t, errors.Is(convertedErr, errToConvert))
	// XXX: not really sure i like this.
	assert.False(t, errors.Is(errToConvert, convertedErr))

	// Switches _still_ fucking fail.
	// FFF Gavin!!
	// should we just accept error wrapping?
	switch errors.Unwrap(err2) {
	case internal.ErrExtendedExample:
	default:
		assert.Fail(t, "was supposed to reach case above")
	}

}

func TestExtendedError_CorrectlyLogged(t *testing.T) {
	err, ok := internal.L1().(gerrors.Error)
	require.Truef(t, ok, "error must implement the gerror.Error interface")
	assert.Contains(t, err.Error(), "GRPCStatus: InvalidArgument, ")
	assert.Contains(t, err.Error(), "SomeMessage: Print this message, ")
	assert.NotContains(t, err.Error(), "DoNotPrint")
	assert.NotContains(t, err.Error(), "this is for internal issue only")
	assert.Equal(t, "ErrExtendedExample", err.ErrName())
	assert.Equal(t, "extended error example", err.ErrMessage())
	assert.Equal(t, "internal:L3", err.ErrSource())
	assert.Equal(t, "", err.ErrDetailTag())
}
