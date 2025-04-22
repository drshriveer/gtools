//nolint:all // this is a test file
package nestedpkg

import (
	v4dupepkg "github.com/drshriveer/gtools/gencommon/internal/nestedpkg/v4"
	v5dupepkg "github.com/drshriveer/gtools/gencommon/internal/nestedpkg/v5"
)

type EmbeddedC struct{}

func (e *EmbeddedC) FooMethod() {}
func (e *EmbeddedC) BazMethod() {}
func (c *EmbeddedC) MethodWithParamsFromPackagesWithSameName(
	_ v4dupepkg.InputTypeAtV4,
	_ v5dupepkg.InputTypeAtV5,
) {
}

// CMethod is a method with a comment, the rare and only method with a comment.
func (e *EmbeddedC) CMethod() {}

type EmbeddedD interface {
	DMethod()
}

type SomeID uint64
