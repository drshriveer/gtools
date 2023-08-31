package gerrors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/drshriveer/gcommon/pkg/gerrors"
)

var ErrMyError1 = gerrors.GError{
	Name:    "ErrMyError1",
	Message: "this is error 1",
}

type AType struct {
}

func (a AType) ReturnsError() error {
	return ErrMyError1.WithStack()
}

func TestGError_WithStack(t *testing.T) {
	err := ErrMyError1.WithStack()
	assert.NotSame(t, ErrMyError1, &err)
	assert.Equal(t, "gerrors_test:TestGError_WithStack", err.Source)
	assert.NotEmpty(t, err.Stack)
	assert.Same(t, err.SrcFactory, &ErrMyError1)
	// ensure unchanged:
	assert.Empty(t, ErrMyError1.Source)
	assert.Empty(t, ErrMyError1.Stack)

	strukt := &AType{}
	err = strukt.ReturnsError().(gerrors.GError)
	assert.NotSame(t, ErrMyError1, &err)
	assert.Equal(t, "gerrors_test:AType:ReturnsError", err.Source)
	assert.NotEmpty(t, err.Stack)
	assert.Same(t, err.SrcFactory, &ErrMyError1)
	// ensure unchanged:
	assert.Empty(t, ErrMyError1.Source)
	assert.Empty(t, ErrMyError1.Stack)
}

func TestGError_Include(t *testing.T) {
	err := ErrMyError1.Include("T-Shirts $%d", 5)
	assert.NotSame(t, ErrMyError1, &err)
	assert.Equal(t, "gerrors_test:TestGError_Include", err.Source)
	assert.NotEmpty(t, err.Stack)
	assert.Equal(t, "T-Shirts $5", err.ExtMessage)
	assert.Same(t, err.SrcFactory, &ErrMyError1)
	// ensure unchanged:
	assert.Empty(t, ErrMyError1.Source)
	assert.Empty(t, ErrMyError1.Stack)
	assert.Empty(t, ErrMyError1.ExtMessage)

	switch gerrors.Unwrap(err) {
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
// BenchmarkGError_Raw-10           	41166675	        28.91 ns/op <-- with pointer
func BenchmarkGError_Raw(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ErrMyError1.Raw()
	}
}
