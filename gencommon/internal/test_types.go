//nolint:all // this is a test file
package internal

import (
	"github.com/drshriveer/gtools/gencommon/internal/nestedpkg"
)

// An alias to a type in a different package.
type AliasID = nestedpkg.SomeID

// Comment to throw us off the trail.

// TypeToGenerate has a comment.
// SecondLine of expected comment.
type TypeToGenerate struct {
	EmbeddedA
	EmbeddedB
}

func (t *TypeToGenerate) ParentMethod()              {}
func (t *TypeToGenerate) BazMethod()                 {}
func (t *TypeToGenerate) pooMethod()                 {}
func (t *TypeToGenerate) MethodTakesAlias(_ AliasID) {}

type EmbeddedA struct{ nestedpkg.EmbeddedD }

func (e *EmbeddedA) FooMethod() {}
func (e *EmbeddedA) BarMethod() {}
func (e *EmbeddedA) AMethod()   {}

type EmbeddedB struct{ *nestedpkg.EmbeddedC }

func (e *EmbeddedB) FooMethod() {}
func (e *EmbeddedB) BazMethod() {}
func (e *EmbeddedB) BMethod()   {}
func (e *EmbeddedB) bPrivate()  {}
