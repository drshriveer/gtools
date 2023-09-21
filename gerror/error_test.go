package gerror_test

import (
	"testing"

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

func (a AType) InlineError() error {
	return func() error {
		return func() error {
			return ErrMyError1.Stack()
		}()
	}()
}

func TestGError_WithStack(t *testing.T) {
	err := ErrMyError1.Stack().(gerror.Error)
	assert.NotSame(t, ErrMyError1, &err)
	assert.Equal(t, "gerror_test:TestGError_WithStack", err.ErrSource())
	assert.NotEmpty(t, err.ErrStack())
	assert.Same(t, gerror.ExtractFactoryReference(err), ErrMyError1)
	// ensure unchanged:
	assert.Empty(t, ErrMyError1.(gerror.Error).ErrStack())
	assert.Empty(t, ErrMyError1.(gerror.Error).ErrStack())

	strukt := &AType{}
	err = strukt.ReturnsError().(gerror.Error)
	assert.NotSame(t, ErrMyError1, &err)
	assert.Equal(t, "gerror_test:AType:ReturnsError", err.ErrSource())
	assert.NotEmpty(t, err.ErrStack())
	assert.Same(t, gerror.ExtractFactoryReference(err), ErrMyError1)
	// ensure unchanged:
	assert.Empty(t, ErrMyError1.(gerror.Error).ErrSource())
	assert.Empty(t, ErrMyError1.(gerror.Error).ErrStack())

	err = strukt.InlineError().(gerror.Error)
	assert.Equal(t, "gerror_test:AType:ReturnsError", err.ErrSource())
	for _, el := range err.ErrStack() {
		println(el.Metric())
	}

}

func TestGError_ExtMsgf(t *testing.T) {
	err := ErrMyError1.ExtMsgf("T-Shirts $%d", 5).(gerror.Error)
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
