//nolint:all // this is a test file
package nestedpkg

import (
	"github.com/drshriveer/gtools/gencommon/internal/duplicatename/a/duplicate"
	bAlias "github.com/drshriveer/gtools/gencommon/internal/duplicatename/b/duplicate"
)

type EmbeddedC struct{}

func (e *EmbeddedC) FooMethod()                                                           {}
func (e *EmbeddedC) BazMethod()                                                           {}
func (e *EmbeddedC) EmbeddedMethodTakesAlias(a duplicate.ADuplicate, b bAlias.BDuplicate) {}

// CMethod is a method with a comment, the rare and only method with a comment.
func (e *EmbeddedC) CMethod() {}

type EmbeddedD interface {
	DMethod()
}

type SomeID uint64
