package gen

import (
	"sort"
)

//go:generate gsort --types ErrorDesc=ErrorDescs,Field=Fields

// ErrorDesc describes an error.
type ErrorDesc struct {
	TypeName string `gsort:"1"`
	Fields   Fields
}

// FieldsToPrint returns fields that need to be included in Error() messages.
func (e *ErrorDesc) FieldsToPrint() Fields {
	r := filter(e.Fields, func(t *Field) bool { return t.Print })
	sort.Sort(Fields(r))
	return r
}

// FieldsToClone returns fields that need to be cloned.
func (e *ErrorDesc) FieldsToClone() Fields {
	r := filter(e.Fields, func(t *Field) bool { return t.Clone })
	sort.Sort(Fields(r))
	return r
}

// Field is an error field and its meaning.
type Field struct {
	Name    string `gsort:"1"`
	PrintAs string
	Clone   bool
	Print   bool
}

func filter[T any](in []T, include func(T) bool) []T {
	result := make([]T, 0, len(in))
	for _, v := range in {
		if include(v) {
			result = append(result, v)
		}
	}
	return result
}
