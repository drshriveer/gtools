//nolint:all // this is a test file
package internal

type TypeToGenerate struct {
	EmbeddedA
	EmbeddedB
}

func (t *TypeToGenerate) ParentMethod() {}
func (t *TypeToGenerate) BazMethod()    {}
func (t *TypeToGenerate) pooMethod()    {}

type EmbeddedA struct{}

func (e *EmbeddedA) FooMethod() {}
func (e *EmbeddedA) BarMethod() {}
func (e *EmbeddedA) AMethod()   {}

type EmbeddedB struct{ EmbeddedC }

func (e *EmbeddedB) FooMethod() {}
func (e *EmbeddedB) BazMethod() {}
func (e *EmbeddedB) BMethod()   {}
func (e *EmbeddedB) bPrivate()  {}

type EmbeddedC struct{}

func (e *EmbeddedC) FooMethod() {}
func (e *EmbeddedC) BazMethod() {}
func (e *EmbeddedC) CMethod()   {}
