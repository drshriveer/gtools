package gencommon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getSafeParamName(t *testing.T) {
	prefixCounter := map[string]int{}
	tests := []struct {
		paramName    string
		alwaysNumber bool
		expected     string
	}{
		{paramName: "ret", alwaysNumber: true, expected: "ret0"},
		{paramName: "ret", alwaysNumber: true, expected: "ret1"},
		{paramName: "ret", alwaysNumber: true, expected: "ret2"},
		{paramName: "ret", alwaysNumber: false, expected: "ret3"},
		{paramName: "ret1", alwaysNumber: false, expected: "ret1"},

		{paramName: "arg", alwaysNumber: false, expected: "arg"},
		{paramName: "arg", alwaysNumber: true, expected: "arg0"},
		{paramName: "arg", alwaysNumber: true, expected: "arg1"},
		{paramName: "arg", alwaysNumber: true, expected: "arg2"},

		{paramName: "someArg", alwaysNumber: false, expected: "someArg"},
		{paramName: "somethingelse", alwaysNumber: false, expected: "somethingelse"},
	}

	for _, test := range tests {
		t.Run("expected "+test.expected, func(t *testing.T) {
			result := getSafeParamName(prefixCounter, test.paramName, test.alwaysNumber)
			assert.Equal(t, test.expected, result)
		})
	}
}
