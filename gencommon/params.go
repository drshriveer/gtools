package gencommon

import (
	"go/ast"
	"go/types"
	"strconv"
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
func ParamsFromSignatureTuple(ih *ImportHandler, tuple *types.Tuple, variadic bool) Params {
	result := make(Params, tuple.Len())
	for i := 0; i < tuple.Len(); i++ {
		v := tuple.At(i)
		p := &Param{
			ActualType: v.Type(),
			TypeRef:    ih.ExtractTypeRef(v.Type()),
			Name:       v.Name(),
			Variadic:   tuple.Len() == i+1 && variadic,
			// Comments: nil, // FIXME: need AST
		}

		// ...Variadic's seem to be very forced into the language
		// They exist at a signature level, but not lower.
		// Lower, the type just resolves to a slice, so we need to trim that out.
		if p.Variadic {
			p.TypeRef = strings.TrimPrefix(p.TypeRef, "[]")
		}

		result[i] = p
	}
	return result
}

// TypeNames returns a comma-separated list of the parameter types.
func (ps Params) TypeNames() string {
	result := strings.Builder{}
	for i, p := range ps {
		// probably not valid?
		if p.Variadic {
			result.WriteString("[]")
		}
		result.WriteString(p.TypeRef)
		if i+1 < len(ps) {
			result.WriteString(", ")
		}
	}
	return result.String()
}

// ParamNames returns a comma-separated list of the parameter names.
// e.g. arg1, arg2, arg3...
func (ps Params) ParamNames() string {
	result := strings.Builder{}
	for i, p := range ps {
		result.WriteString(p.Name)
		if p.Variadic {
			result.WriteString("...")
		}
		if i+1 < len(ps) {
			result.WriteString(", ")
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
		if p.Variadic {
			result.WriteString("...")
		}
		result.WriteString(" ")
		result.WriteString(p.TypeRef)
		if i+1 < len(ps) {
			result.WriteString(", ")
		}
	}
	return result.String()
}

func (ps Params) ensureNames(paramDeduper map[string]int, isOutput bool) {
	// paramName == the unnamed parameter paramName to use.
	prefix := "arg"
	if isOutput {
		prefix = "ret"
	}
	// To better preserve a customer's naming in case of them colliding with our own,
	// process the named variables first:
	for _, p := range ps {
		if p.Name != "" {
			p.Name = getSafeParamName(paramDeduper, prefix, false)
		}
	}

	for i, p := range ps {
		if p.Name == "" {
			if isOutput && len(ps)-1 == i && types.Implements(p.ActualType, ErrorInterface) {
				p.Name = getSafeParamName(paramDeduper, "err", false)
			} else if !isOutput && i == 0 && types.Implements(p.ActualType, ContextInterface) {
				p.Name = getSafeParamName(paramDeduper, "ctx", false)
			} else {
				p.Name = getSafeParamName(paramDeduper, prefix, true)
			}
		}
	}
}

// Param has information about a single parameter.
type Param struct {
	ActualType types.Type
	TypeRef    string
	Name       string
	Comments   Comments
	Variadic   bool
}

// Declaration returns a name and type.
func (p Param) Declaration() string {
	return p.Name + " " + p.TypeRef
}

// getSafeParamName returns a "safe" param name.
// note: I'm pretty sure this is technically only safe when the already defined params
// are processed first which is exactly what ensureNames does.
func getSafeParamName(paramDeduper map[string]int, paramName string, alwaysNumber bool) string {
	v, ok := paramDeduper[paramName]
	result := paramName
	if ok || alwaysNumber {
		result += strconv.FormatInt(int64(v), 10)
		v++
	}
	// else don't modify the intended paramName.
	// ensure the paramName is in the map:
	paramDeduper[paramName] = v
	return result
}
