package gencommon

import (
	"go/ast"
	"go/types"
	"strings"
)

// Params is a handle on a list of params with useful methods.
type Params []*Param

// ParamsFromFieldList converts an *ast.FieldList to a list of params..
func ParamsFromFieldList(fields *ast.FieldList) Params {
	if fields.NumFields() == 0 {
		return nil
	}
	result := make(Params, fields.NumFields())
	for i, field := range fields.List {
		result[i] = &Param{
			Comments: docToString(field.Comment),
			TypeRef:  types.ExprString(field.Type),
			Name:     getName(field.Names...),
		}
	}

	return result
}

// TypeNames returns a comma-separated list of the parameter types.
func (p Params) TypeNames() string {
	result := strings.Builder{}
	for i, p_ := range p {
		result.WriteString(p_.TypeRef)
		if i+1 < len(p) {
			result.WriteString(",")
		}
	}
	return result.String()
}

// ParamsNames returns a comma-separated list of the parameter names.
// e.g. arg1, arg2, arg3...
func (p Params) ParamsNames() string {
	result := strings.Builder{}
	for i, p_ := range p {
		result.WriteString(p_.Name)
		if i+1 < len(p) {
			result.WriteString(",")
		}
	}
	return result.String()
}

// Declarations returns a comma-separated list of parameter name and type:
// e.g. arg1 Type1, arg2 Type2 ...,
func (p Params) Declarations() string {
	result := strings.Builder{}
	for i, p_ := range p {
		result.WriteString(p_.Name)
		result.WriteString(" ")
		result.WriteString(p_.TypeRef)
		if i+1 < len(p) {
			result.WriteString(",")
		}
	}
	return result.String()
}

// Param has information about a signle paramter.
type Param struct {
	TypeRef  string
	Name     string
	Comments []string
}

// Declaration returns a name and type.
func (p Param) Declaration() string {
	return p.Name + " " + p.TypeRef
}
