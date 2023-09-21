package gerror_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"

	"github.com/drshriveer/gtools/gerror"
)

var ErrMyError1 = gerror.FactoryOf(&gerror.GError{
	Name:    "ErrMyError1",
	Message: "this is error 1",
})

type AType struct {
}

func (a AType) ReturnsError() error {
	return ErrMyError1.Stack()
}

func (a AType) InlineError2XL() error {
	return func() error {
		return func() error {
			return ErrMyError1.Stack()
		}()
	}()
}

func (a AType) InlineError3XL() error {
	return func() error {
		return func() error {
			return func() error {
				return ErrMyError1.Stack()
			}()
		}()
	}()
}

func TestGError_WithStack(t *testing.T) {
	err := ErrMyError1.Stack()
	assert.NotSame(t, ErrMyError1, &err)
	assert.Equal(t, "gerror_test:TestGError_WithStack", err.ErrSource())
	assert.NotEmpty(t, err.ErrStack())
	assert.Same(t, gerror.ExtractFactoryReference(err), ErrMyError1)
	// ensure unchanged:
	assert.Empty(t, ErrMyError1.(gerror.Error).ErrStack())
	assert.Empty(t, ErrMyError1.(gerror.Error).ErrStack())

	aType := &AType{}
	err = aType.ReturnsError().(gerror.Error)
	assert.NotSame(t, ErrMyError1, &err)
	assert.Equal(t, "gerror_test:AType:ReturnsError", err.ErrSource())
	assert.NotEmpty(t, err.ErrStack())
	assert.Same(t, gerror.ExtractFactoryReference(err), ErrMyError1)
	// ensure unchanged:
	assert.Empty(t, ErrMyError1.(gerror.Error).ErrSource())
	assert.Empty(t, ErrMyError1.(gerror.Error).ErrStack())
}

//nolint:errname // for testing
var pkgVarInlineFunc1 = func() error {
	return ErrMyError1.Stack()
}()

func pkgFuncInlineErr1() error {
	return func() error {
		return ErrMyError1.Stack()
	}()
}

func pkgFuncInlineErr2() error {
	return func() error {
		return func() error {
			return ErrMyError1.Stack()
		}()
	}()
}

func inlineAndOtherDeepLayer() error {
	aType := &AType{}
	return func() error {
		return aType.InlineError2XL()
	}()
}

func genericFuncResponse[T any](_ T) error {
	return ErrMyError1.Stack()
}

func TestStackSource(t *testing.T) {
	aType := &AType{}
	tests := []struct {
		description string
		err         error

		expected string
	}{
		{
			// Looks like: gerror_test:TestStackSource
			description: "local (instruct)",
			err:         ErrMyError1.Stack(),
			expected:    "gerror_test:TestStackSource",
		},
		{
			// Looks like: gerror_test:TestStackSource:func1
			description: "local (in instruct> anonymous",
			err: func() error {
				return ErrMyError1.Stack()
			}(),
			expected: "gerror_test:TestStackSource",
		},
		{
			// Looks Like: gerror_test:AType:ReturnsError
			description: "struct method returns error",
			err:         aType.ReturnsError(),
			expected:    "gerror_test:AType:ReturnsError",
		},
		{
			// Looks Like: github.com/drshriveer/gtools/gerror_test.AType.InlineError2XL.AType.InlineError2XL.func1.func2
			// Looks Like: <PKG>.AType.InlineError2XL.AType.InlineError2XL.func1.func2
			// Looks Like: <PKG>.<TYPE>.InlineError2XL.AType.InlineError2XL.func1.func2
			// Looks Like:
			description: "struct method > 2 layers of anonymous functions",
			err:         aType.InlineError2XL(),
			// is: gerror_test:AType:InlineError2XL:AType:InlineError2XL
			expected: "gerror_test:AType:InlineError2XL",
		},
		{
			// Looks Like: github.com/drshriveer/gtools/gerror_test.TestStackSource.AType.InlineError3XL.func3.1.1
			// Looks Like: <PKG>.TestStackSource.AType.InlineError3XL.func3.1.1
			// Looks Like: <PKG>.<nearest-function>.AType.InlineError3XL.func3.1.1
			// Looks Like: <PKG>.<nearest-function>.<type>.<type-function>.func3.1.1
			description: "struct method > 3 layers of anonymous functions",
			err:         aType.InlineError3XL(),
			expected:    "gerror_test:TestStackSource:AType:InlineError3XL",
		},
		{
			// Looks Like: github.com/drshriveer/gtools/gerror_test.glob..func1
			// Looks Like: <PKG>.glob..func1
			description: "pkg var > anonymous",
			err:         pkgVarInlineFunc1,
			expected:    "gerror_test:glob",
		},
		{
			// Looks Like: github.com/drshriveer/gtools/gerror_test.TestStackSource.pkgFuncInlineErr1.func4
			// Looks Like: <PKG>.TestStackSource.pkgFuncInlineErr1.func4
			// Looks Like: <PKG>.<nearest-function>.pkgFuncInlineErr1.func4
			// Looks Like: <PKG>.<nearest-function>.<named-func>.func4
			description: "pkg func> anonymous 1",
			err:         pkgFuncInlineErr1(),
			expected:    "gerror_test:TestStackSource:pkgFuncInlineErr1",
		},
		{
			// Looks Like: <PKG>.pkgFuncInlineErr2.pkgFuncInlineErr2.func1.func2
			description: "pkg func> anonymous 2",
			err:         pkgFuncInlineErr2(),
			expected:    "gerror_test:pkgFuncInlineErr2",
		},
		{
			// Looks Like: github.com/drshriveer/gtools/gerror_test.AType.InlineError2XL.AType.InlineError2XL.func1.func2
			// Looks Like: <PKG>.AType.InlineError2XL.AType.InlineError2XL.func1.func2
			// Looks Like: <PKG>.<type>.InlineError2XL.AType.InlineError2XL.func1.func2
			// Looks Like: <PKG>.<type>.<type-func>.AType.InlineError2XL.func1.func2
			description: "oddities at large",
			err:         inlineAndOtherDeepLayer(),
			expected:    "gerror_test:AType:InlineError2XL",
		},
		{
			description: "a generic what",
			err:         genericFuncResponse[string]("bhahahahhaha"),
			expected:    "gerror_test:genericFuncResponse",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			gerr, ok := test.err.(gerror.Error)
			require.True(t, ok)
			assert.Equal(t, test.expected, gerr.ErrSource())
		})
	}
}

func TestGError_ExtMsgf(t *testing.T) {
	err := ErrMyError1.ExtMsgf("T-Shirts $%d", 5)
	assert.NotSame(t, ErrMyError1, &err)
	assert.Equal(t, "gerror_test:TestGError_ExtMsgf", err.ErrSource())
	assert.NotEmpty(t, err.ErrStack())
	assert.Equal(t, ErrMyError1.(gerror.Error).ErrMessage()+" T-Shirts $5", err.ErrMessage())
	assert.Same(t, gerror.ExtractFactoryReference(err), ErrMyError1)
	// ensure unchanged:
	assert.Empty(t, ErrMyError1.(gerror.Error).ErrSource())
	assert.Empty(t, ErrMyError1.(gerror.Error).ErrStack())
	assert.Equal(t, "this is error 1", ErrMyError1.(gerror.Error).ErrMessage())
}

// BenchmarkGError_WithSource-10    	 2224312	       530.9 ns/op
// BenchmarkGError_WithSource-10    	 2161904	       557.3 ns/op <-- with pointer

func BenchmarkGError_WithSource(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ErrMyError1.Src()
	}
}

// BenchmarkGError_WithStack-10     	  978735	      1148 ns/op
// BenchmarkGError_WithStack-10     	 1000000	      1157 ns/op <-- with pointer

func BenchmarkGError_WithStack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ErrMyError1.Stack()
	}
}

// BenchmarkGError_Raw-10           	135645453	         8.868 ns/op
// BenchmarkGError_Raw-10           	41166675	        28.91 ns/op <-- with pointer.
func BenchmarkGError_Raw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ErrMyError1.Base()
	}
}
