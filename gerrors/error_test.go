package gerrors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var ErrMyError1 = GError{
	Name:    "ErrMyError1",
	Message: "this is error 1",
}

type AType struct {
}

func (a AType) ReturnsError() error {
	return ErrMyError1.WithStack()
}

// FIXME! TEST INLINE FUNCTION

func TestGError_WithStack(t *testing.T) {
	err := ErrMyError1.WithStack().(*GError)
	assert.NotSame(t, ErrMyError1, &err)
	assert.Equal(t, "gerrors:TestGError_WithStack", err.Source)
	assert.NotEmpty(t, err.stack)
	assert.Same(t, err.srcFactory, &ErrMyError1)
	// ensure unchanged:
	assert.Empty(t, ErrMyError1.Source)
	assert.Empty(t, ErrMyError1.stack)

	strukt := &AType{}
	err = strukt.ReturnsError().(*GError)
	assert.NotSame(t, ErrMyError1, &err)
	assert.Equal(t, "gerrors:AType:ReturnsError", err.Source)
	assert.NotEmpty(t, err.stack)
	assert.Same(t, err.srcFactory, &ErrMyError1)
	// ensure unchanged:
	assert.Empty(t, ErrMyError1.Source)
	assert.Empty(t, ErrMyError1.stack)
}

func TestGError_ExtMsgf(t *testing.T) {
	err := ErrMyError1.ExtMsgf("T-Shirts $%d", 5).(*GError)
	assert.NotSame(t, ErrMyError1, &err)
	assert.Equal(t, "gerrors:TestGError_ExtMsgf", err.Source)
	assert.NotEmpty(t, err.stack)
	assert.Equal(t, ErrMyError1.Message+" T-Shirts $5", err.Message)
	assert.Same(t, err.srcFactory, &ErrMyError1)
	// ensure unchanged:
	assert.Empty(t, ErrMyError1.Source)
	assert.Empty(t, ErrMyError1.stack)
	assert.Equal(t, "this is error 1", ErrMyError1.Message)

	switch Unwrap(err) {
	case &ErrMyError1:
	default:
		assert.Fail(t, "darn")
	}
}

//
// func L1() error {
// 	return L2()
// }
// func L2() error {
// 	return L3()
// }
// func L3() error {
// 	return ErrMyError1.
// }

// BenchmarkGError_WithSource-10    	 2224312	       530.9 ns/op
// BenchmarkGError_WithSource-10    	 2161904	       557.3 ns/op <-- with pointer

func BenchmarkGError_WithSource(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ErrMyError1.WithSource()
	}
}

// BenchmarkGError_WithStack-10     	  978735	      1148 ns/op
// BenchmarkGError_WithStack-10     	 1000000	      1157 ns/op <-- with pointer

func BenchmarkGError_WithStack(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ErrMyError1.WithStack()
	}
}

// BenchmarkGError_Raw-10           	135645453	         8.868 ns/op
// BenchmarkGError_Raw-10           	41166675	        28.91 ns/op <-- with pointer.
func BenchmarkGError_Raw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ErrMyError1.Base()
	}
}
