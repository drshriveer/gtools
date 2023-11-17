package gencommon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getSafeParamName(t *testing.T) {
	prefixCounter := map[string]int{}
	tests := []struct {
		prefix       string
		alwaysNumber bool
		expected     string
	}{
		{prefix: "ret", alwaysNumber: true, expected: "ret0"},
		{prefix: "ret", alwaysNumber: true, expected: "ret1"},
		{prefix: "ret", alwaysNumber: true, expected: "ret2"},
		{prefix: "ret", alwaysNumber: false, expected: "ret3"},
		{prefix: "ret1", alwaysNumber: false, expected: "ret1"},

		{prefix: "arg", alwaysNumber: false, expected: "arg"},
		{prefix: "arg", alwaysNumber: true, expected: "arg0"},
		{prefix: "arg", alwaysNumber: true, expected: "arg1"},
		{prefix: "arg", alwaysNumber: true, expected: "arg2"},

		{prefix: "someArg", alwaysNumber: false, expected: "someArg"},
		{prefix: "somethingelse", alwaysNumber: false, expected: "somethingelse"},
	}

	for _, test := range tests {
		t.Run("expected "+test.expected, func(t *testing.T) {
			result := getSafeParamName(prefixCounter, test.prefix, test.alwaysNumber)
			assert.Equal(t, test.expected, result)
		})
	}
}
