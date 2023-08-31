package internal_test

import (
	"testing"

	"github.com/drshriveer/gcommon/pkg/genum"
	"github.com/drshriveer/gcommon/pkg/genum/internal"
)

func TestMyEnum(t *testing.T) {
	_, ok := any(internal.UNSET).(genum.Enum)
	if !ok {
		t.FailNow()
	}
}
