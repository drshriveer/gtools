package gerrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var ErrMyError1 Factory = &GError{
	Name:    "ErrMyError1",
	Message: "this is error 1",
}

type AType struct {
}

func (a AType) ReturnsError() error {
	return ErrMyError1.Stack()
}

// FIXME! TEST INLINE FUNCTION

func TestGError_WithStack(t *testing.T) {
	err := ErrMyError1.Stack().(*GError)
	assert.NotSame(t, ErrMyError1, &err)
	assert.Equal(t, "gerrors:TestGError_WithStack", err.Source)
	assert.NotEmpty(t, err.stack)
	assert.Same(t, err.srcFactory, &ErrMyError1)
	// ensure unchanged:
	assert.Empty(t, ErrMyError1.(*GError).Source)
	assert.Empty(t, ErrMyError1.(*GError).stack)

	strukt := &AType{}
	err = strukt.ReturnsError().(*GError)
	assert.NotSame(t, ErrMyError1, &err)
	assert.Equal(t, "gerrors:AType:ReturnsError", err.Source)
	assert.NotEmpty(t, err.stack)
	assert.Same(t, err.srcFactory, &ErrMyError1)
	// ensure unchanged:
	assert.Empty(t, ErrMyError1.(*GError).Source)
	assert.Empty(t, ErrMyError1.(*GError).stack)
}

func TestGError_ExtMsgf(t *testing.T) {
	err := ErrMyError1.ExtMsgf("T-Shirts $%d", 5).(*GError)
	assert.NotSame(t, ErrMyError1, &err)
	assert.Equal(t, "gerrors:TestGError_ExtMsgf", err.Source)
	assert.NotEmpty(t, err.stack)
	assert.Equal(t, ErrMyError1.(*GError).Message+" T-Shirts $5", err.Message)
	assert.Same(t, err.srcFactory, &ErrMyError1)
	// ensure unchanged:
	assert.Empty(t, ErrMyError1.(*GError).Source)
	assert.Empty(t, ErrMyError1.(*GError).stack)
	assert.Equal(t, "this is error 1", ErrMyError1.(*GError).Message)
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
