//nolint:all // this is a test file
package nestedpkg

type EmbeddedC struct{}

func (e *EmbeddedC) FooMethod() {}
func (e *EmbeddedC) BazMethod() {}

// CMethod is a method with a comment, the rare and only method with a comment.
func (e *EmbeddedC) CMethod() {}

type EmbeddedD interface {
	DMethod()
}

type SomeID uint64
