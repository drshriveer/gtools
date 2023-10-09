package internal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/drshriveer/gtools/gerror"
	"github.com/drshriveer/gtools/gerror/internal"
)

func TestErrWithCustomConvert(t *testing.T) {
	ErrTest := internal.ErrWithCustomConvert{
		GError:   gerror.GError{Name: "ErrTest"},
		Property: 0,
	}
	assert.Implements(t, (*gerror.Error)(nil), ErrTest.Msg("blah"))
}
