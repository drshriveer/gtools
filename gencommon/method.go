package gencommon

import (
	"fmt"
	"go/ast"
	"strings"
)

// Method is a type representing a method.
type Method struct {
	Name       string `gsort:"2"`
	Comments   []string
	Input      Params
	Output     Params
	IsExported bool `gsort:"1"`
}

// MethodFrom attempts to create a new Method type from a function declaration.
// It returns false if the function declaration is not a method of the target type.
func MethodFrom(decl *ast.FuncDecl, targetType string) (*Method, bool) {
	// for now there can only be one receiver right???
	if decl.Recv == nil || len(decl.Recv.List) != 1 {
		return nil, false
	}

	switch t := decl.Recv.List[0].Type.(type) {
	case *ast.StarExpr:
		ident, ok := t.X.(*ast.Ident)
		if !ok {
			return nil, false
		}
		if ident.Name != targetType {
			return nil, false
		}

		m := &Method{
			Comments:   docToString(decl.Doc),
			Name:       getName(decl.Name),
			IsExported: decl.Name.IsExported(),
			Input:      ParamsFromFieldList(decl.Type.Params),
			Output:     ParamsFromFieldList(decl.Type.Results),
		}
		return m, true
	case *ast.Ident:
		if t.Name != targetType {
			return nil, false
		}
		m := &Method{
			Comments:   docToString(decl.Doc),
			Name:       getName(decl.Name),
			IsExported: decl.Name.IsExported(),
			Input:      ParamsFromFieldList(decl.Type.Params),
			Output:     ParamsFromFieldList(decl.Type.Results),
		}
		return m, true
	default:
	}
	return nil, false
}

// Signature returns the full signature of the method.
// e.g. MethodName(arg1 Type1, arg2 Type2, arg3 Type3) (*Thing, string, error)
func (m *Method) Signature() string {
	if len(m.Output) > 1 {
		return fmt.Sprintf("%s(%s) (%s)", m.Name, m.Input.Declarations(), m.Output.TypeNames())
	}
	return fmt.Sprintf("%s(%s) %s", m.Name, m.Input.Declarations(), m.Output.TypeNames())
}

// Call returns a method call of the form:
// e.g. MethodName(arg1, arg2, arg3).
func (m *Method) Call() string {
	args := mapper(m.Input, func(in *Param) string { return in.Name })
	return fmt.Sprintf("%s(%s) ", m.Name, strings.Join(args, ", "))
}

// ReturnsError indicates whether this method returns an error.
// This method does not use real type information -- it just looks for the type name
// to end in `Error` or `error`.
func (m *Method) ReturnsError() bool {
	if len(m.Output) == 0 {
		return false
	}
	last := m.Output[len(m.Output)-1]
	return strings.HasSuffix(last.TypeRef, "Error") ||
		strings.HasSuffix(last.TypeRef, "error")
}

// CommentBlock processes comments into a single string for convenience.
func (m *Method) CommentBlock() string {
	return strings.Join(m.Comments, "\n")
}

func getName(names ...*ast.Ident) string {
	for _, name := range names {
		return name.Name
	}

	return ""
}
