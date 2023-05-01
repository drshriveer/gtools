package tester_test

import (
	"github.com/drshriveer/gsenum/pkg/enum"
	"github.com/drshriveer/gsenum/pkg/tester"
	"testing"
)

func TestMyEnum(t *testing.T) {
	_, ok := any(tester.UNSET).(enum.Enum[tester.MyEnum])
	if !ok {
		t.FailNow()
	}
}
