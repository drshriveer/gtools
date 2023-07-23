package internal_test

import (
	"testing"

	"github.com/drshriveer/gcommon/pkg/enum"
	"github.com/drshriveer/gcommon/pkg/enum/internal"
)

func TestMyEnum(t *testing.T) {
	_, ok := any(internal.UNSET).(enum.Enum)
	if !ok {
		t.FailNow()
	}
}
