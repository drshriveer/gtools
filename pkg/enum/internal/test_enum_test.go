package internal_test

import (
	"github.com/drshriveer/gcommon/pkg/enum"
	"github.com/drshriveer/gcommon/pkg/enum/internal"
	"testing"
)

func TestMyEnum(t *testing.T) {
	_, ok := any(internal.UNSET).(enum.Enum[internal.MyEnum])
	if !ok {
		t.FailNow()
	}
}
