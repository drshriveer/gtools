package gencommon

import (
	"fmt"
	"go/ast"
	"go/types"
	"path"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

// ImportDesc is a description of an import.
type ImportDesc struct {
	Alias              string
	PkgPath            string
	aliasIsPackageName bool
	inUse              bool
}

// ImportString returns a formatted import string with alias (if required).
func (id *ImportDesc) ImportString() string {
	if id.aliasIsPackageName {
		return "\"" + id.PkgPath + "\""
	}
	return id.Alias + " \"" + id.PkgPath + "\""
}

// ImportHandler a collection of import descriptions indexed by package path.
type ImportHandler struct {
	PInfo   *packages.Package
	imports map[string]*ImportDesc
}

// calcImports the imports relevant to a specific package and ImportSpec.
func calcImports(pkg *packages.Package, fAST *ast.File) *ImportHandler {
	result := &ImportHandler{
		PInfo:   pkg,
		imports: make(map[string]*ImportDesc),
	}

	// Note: we do this loop here because we understand import aliases in this path.
	for _, iSpec := range fAST.Imports {
		pkgPath := strings.Trim(iSpec.Path.Value, `"`)

		// using the base path as the alias is not always correct, the package directory
		// is not necessarily the same as the package name within that directory.
		// e.g. something/v2 where the underlying package is actually 'something',
		// but it serves as a default when there is no other info.
		alias := path.Base(pkgPath)
		aliasIsPackageName := true
		if iSpec.Name != nil {
			alias = iSpec.Name.Name
			aliasIsPackageName = false
		} else if importPkg, ok := pkg.Imports[pkgPath]; ok {
			alias = importPkg.Name
			aliasIsPackageName = true
		}

		result.imports[pkgPath] = &ImportDesc{
			PkgPath:            pkgPath,
			Alias:              alias,
			inUse:              false,
			aliasIsPackageName: aliasIsPackageName,
		}
	}

	return result
}

// ExtractTypeRef returns the way the type should be referenced in code.
func (ih *ImportHandler) ExtractTypeRef(typ types.Type) string {
	// "named" means it is a type which may require importing.
	switch t := typ.(type) {
	case *types.Pointer:
		return "*" + ih.ExtractTypeRef(t.Elem())
	case *types.Slice:
		return "[]" + ih.ExtractTypeRef(t.Elem())
	case *types.Array:
		return fmt.Sprintf("[%d]", t.Len()) + ih.ExtractTypeRef(t.Elem())
	case *types.Signature:
		// recurse to register relevant method imports-> then we only need the signature.
		m := MethodFromSignature(ih, t)
		return m.Signature()
	case *types.Map:
		// recurse to register types.
		key := ih.ExtractTypeRef(t.Key())
		value := ih.ExtractTypeRef(t.Elem())
		return "map[" + key + "]" + value
	case *types.Alias:
		return ih.addNamed(t)
	case *types.Named:
		return ih.addNamed(t)
	default:
		// *types.Interface is usually handled here too.
		// "*types.Basic"s e.g. string come out as "untyped string"; we need to drop
		//  that part... Not sure why this is how the type information is conveyed :-/.
		return strings.TrimPrefix(t.String(), "untyped ")
	}
}

type named interface {
	Obj() *types.TypeName
	TypeArgs() *types.TypeList
}

func (ih *ImportHandler) addNamed(t named) string {
	pkg := t.Obj().Pkg()
	typeName := t.Obj().Name()
	alias := ""

	// If we need an import, find it and use the proper alias
	if pkg != nil && pkg.Path() != ih.PInfo.PkgPath {
		// first check if we have a mapping for the package:
		i, ok := ih.imports[pkg.Path()]
		if ok {
			i.inUse = true
		} else {
			pkgAlias := pkg.Name()
			var aliasIsPackageName bool
			if importPkg, ok := ih.PInfo.Imports[pkg.Path()]; ok && pkgAlias == "" {
				pkgAlias = importPkg.Name
				aliasIsPackageName = true
			} // else we don't know, so we just want to say aliasIsPackageName = false

			i = &ImportDesc{
				Alias:              pkgAlias,
				PkgPath:            pkg.Path(),
				inUse:              true,
				aliasIsPackageName: aliasIsPackageName,
			}
			ih.imports[i.PkgPath] = i
		}
		alias = i.Alias + "."
	}

	// Recurse into type arguments for generic types
	targs := t.TypeArgs()
	if targs != nil {
		typeArgNames := make([]string, targs.Len())
		for i := 0; i < targs.Len(); i++ {
			typeArg := targs.At(i)
			typeArgNames[i] = ih.ExtractTypeRef(typeArg)
		}
		return fmt.Sprintf("%s%s[%s]", alias, typeName, strings.Join(typeArgNames, ", "))
	}

	return fmt.Sprintf("%s%s", alias, typeName)
}

// GetActive returns ordered, active imports.
// Used by templates.
func (ih *ImportHandler) GetActive() []ImportDesc {
	result := make([]ImportDesc, 0, len(ih.imports))
	for _, i := range ih.imports {
		if i.inUse {
			result = append(result, *i)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].PkgPath < result[j].PkgPath
	})
	return result
}

// HasActiveImports returns true if there are any active imports.
func (ih *ImportHandler) HasActiveImports() bool {
	return len(ih.GetActive()) > 0
}
