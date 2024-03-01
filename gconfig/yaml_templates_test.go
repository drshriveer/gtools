package gconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvVarTmpl_MatchAndResolve(t *testing.T) {
	envVarKey := "MY_ENV_VAR"
	tests := []struct {
		description string
		input       string
		envVal      string

		expectedOutput string
		expectedUsed   bool
		expectedError  error
	}{
		{
			description:    "not a template",
			input:          "4m3s2ms",
			expectedOutput: "4m3s2ms",
		},
		{
			description:    "template found, env var not found",
			input:          "${{env:MY_ENV_VAR}}",
			expectedOutput: "",
			expectedError:  ErrFailedParsing,
		},
		{
			description:    "template found, env var used",
			input:          "${{env:MY_ENV_VAR}}",
			envVal:         "aws:secret",
			expectedUsed:   true,
			expectedOutput: "aws:secret",
		},
		{
			description:    "template found, env var used with spacing 1",
			input:          "${{ env:MY_ENV_VAR }}",
			envVal:         "aws:secret",
			expectedUsed:   true,
			expectedOutput: "aws:secret",
		},
		{
			description:    "template found, env var used with spacing 2",
			input:          "${{env: MY_ENV_VAR}}",
			envVal:         "aws:secret",
			expectedUsed:   true,
			expectedOutput: "aws:secret",
		},
		{
			description:    "template found, env var used with spacing 3",
			input:          "${{   env:   MY_ENV_VAR   }}",
			envVal:         "aws:secret",
			expectedUsed:   true,
			expectedOutput: "aws:secret",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			if test.envVal != "" {
				require.NoError(t, os.Setenv(envVarKey, test.envVal))
				t.Cleanup(func() {
					require.NoError(t, os.Unsetenv(envVarKey))
				})
			}
			result, used, err := envVarTmpl{}.MatchAndResolve(test.input)
			assert.ErrorIs(t, err, test.expectedError)
			assert.Equal(t, test.expectedUsed, used)
			assert.Equal(t, test.expectedOutput, result)

		})
	}
}
