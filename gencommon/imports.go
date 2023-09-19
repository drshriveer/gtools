package gencommon

import (
	"fmt"
	"go/ast"
	"go/types"
	"sort"
	"strings"
)

// ImportDesc is a description of an import.
type ImportDesc struct {
	Alias   string
	PkgPath string
	inUse   bool
}

// ImportDescs a collection of import descriptions indexed by package path.
type ImportDescs struct {
	currentPackage *types.Package
	imports        map[string]*ImportDesc
}

// CalcImports the imports relevant to a specific package and ImportSpec.
func CalcImports(importSpecs []*ast.ImportSpec, pkg *types.Package) ImportDescs {
	result := ImportDescs{
		currentPackage: pkg,
		imports:        make(map[string]*ImportDesc, len(importSpecs)),
	}

	for _, iSpec := range importSpecs {
		pkgPath := strings.Trim(iSpec.Path.Value, `"`)
		result.imports[pkgPath] = &ImportDesc{
			Alias:   iSpec.Name.Name,
			PkgPath: pkgPath,
			inUse:   false,
		}
	}

	return result
}

// ExtractTypeRef returns the way the type should be referenced in code.
func (id ImportDescs) ExtractTypeRef(t types.Type) string {
	// "named" means it is a type which may require importing.
	named, ok := t.(*types.Named)
	if !ok {
		// "*types.Basic"s e.g. string come out as "untyped string"; we need to drop
		//  that part... Not sure why this is how the type information is conveyed :-/.
		return strings.TrimPrefix(t.String(), "untyped ")
	}

	pkg := named.Obj().Pkg()
	typeName := named.Obj().Name()
	if pkg == id.currentPackage {
		return typeName
	}

	// first check if we have a mapping for the package:
	i, ok := id.imports[pkg.Path()]
	if ok {
		i.inUse = true
	} else {
		i = &ImportDesc{
			Alias:   pkg.Name(),
			PkgPath: pkg.Path(),
			inUse:   true,
		}
		id.imports[i.PkgPath] = i
	}

	return fmt.Sprintf("%s.%s", i.Alias, typeName)
}

// GetActive returns ordered, active imports.
// Used by templates.
func (id ImportDescs) GetActive() []ImportDesc {
	result := make([]ImportDesc, 0, len(id.imports))
	for _, i := range id.imports {
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
func (id ImportDescs) HasActiveImports() bool {
	return len(id.GetActive()) > 0
}
