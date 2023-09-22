package gerror_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/drshriveer/gtools/gerror"
)

func TestCloneBase(t *testing.T) {
	t.Parallel()
	ErrNoMessage := gerror.FactoryOf(&gerror.GError{
		Name:    "ErrBase",
		Message: "",
		Source:  "",
	})

	tests := []struct {
		description string
		input       gerror.Error

		expectedName      string
		expectedMessage   string
		expectedSource    string
		expectedDetailTag string
	}{
		{
			description:     "add message does not include extra white space",
			input:           ErrNoMessage.Msg("hi! blah"),
			expectedName:    "ErrBase",
			expectedSource:  "gerror_test:TestCloneBase",
			expectedMessage: "hi! blah",
		},
		{
			description:     "add message is trimmed",
			input:           ErrNoMessage.Msg(" hi! blah "),
			expectedName:    "ErrBase",
			expectedSource:  "gerror_test:TestCloneBase",
			expectedMessage: "hi! blah",
		},
		{
			description:     "add message is trimmed + chaining",
			input:           ErrNoMessage.Msg(" hi! blah ").(*gerror.GError).Msg(". bye! "),
			expectedName:    "ErrBase",
			expectedSource:  "gerror_test:TestCloneBase",
			expectedMessage: "hi! blah. bye!",
		},
		{
			description:     "lots of whitespace is trimmed",
			input:           ErrNoMessage.Msg("     hi! blah  \n"),
			expectedName:    "ErrBase",
			expectedSource:  "gerror_test:TestCloneBase",
			expectedMessage: "hi! blah",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.expectedName, test.input.ErrName())
			assert.Equal(t, test.expectedMessage, test.input.ErrMessage())
			assert.Equal(t, test.expectedSource, test.input.ErrSource())
			assert.Equal(t, test.expectedDetailTag, test.input.ErrDetailTag())
		})
	}
}
