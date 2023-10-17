package gencommon

import (
	"go/ast"
	"go/types"
	"strings"
)

// Params is a handle on a list of params with useful methods.
type Params []*Param

// ParamsFromFieldList converts an *ast.FieldList to a list of params.
func ParamsFromFieldList(fields *ast.FieldList) Params {
	if fields.NumFields() == 0 {
		return nil
	}
	result := make(Params, fields.NumFields())
	for i, field := range fields.List {
		result[i] = &Param{
			Comments: FromCommentGroup(field.Comment),
			TypeRef:  types.ExprString(field.Type),
			Name:     getName(field.Names...),
		}
	}

	return result
}

// ParamsFromSignatureTuple converts an *types.Tuple to params.
func ParamsFromSignatureTuple(ih *ImportHandler, tuple *types.Tuple) Params {
	result := make(Params, tuple.Len())
	for i := 0; i < tuple.Len(); i++ {
		v := tuple.At(i)
		result[i] = &Param{
			actualType: v.Type(),
			TypeRef:    ih.ExtractTypeRef(v.Type()),
			Name:       v.Name(),
			// Comments: nil, // FIXME: need AST
		}
	}
	return result
}

// TypeNames returns a comma-separated list of the parameter types.
func (ps Params) TypeNames() string {
	result := strings.Builder{}
	for i, p := range ps {
		result.WriteString(p.TypeRef)
		if i+1 < len(ps) {
			result.WriteString(",")
		}
	}
	return result.String()
}

// ParamsNames returns a comma-separated list of the parameter names.
// e.g. arg1, arg2, arg3...
func (ps Params) ParamsNames() string {
	result := strings.Builder{}
	for i, p := range ps {
		result.WriteString(p.Name)
		if i+1 < len(ps) {
			result.WriteString(",")
		}
	}
	return result.String()
}

// Declarations returns a comma-separated list of parameter name and type:
// e.g. arg1 Type1, arg2 Type2 ...,.
func (ps Params) Declarations() string {
	result := strings.Builder{}
	for i, p := range ps {
		result.WriteString(p.Name)
		result.WriteString(" ")
		result.WriteString(p.TypeRef)
		if i+1 < len(ps) {
			result.WriteString(",")
		}
	}
	return result.String()
}

// Param has information about a single parameter.
type Param struct {
	actualType types.Type
	TypeRef    string
	Name       string
	Comments   Comments
}

// Declaration returns a name and type.
func (p Param) Declaration() string {
	return p.Name + " " + p.TypeRef
}
