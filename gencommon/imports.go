package gencommon

import (
	"fmt"
	"go/ast"
	"go/types"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"
)

// ImportDesc is a description of an import.
type ImportDesc struct {
	Alias   string
	PkgPath string
	inUse   bool
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

	for _, iSpec := range fAST.Imports {
		pkgPath := strings.Trim(iSpec.Path.Value, `"`)
		id := &ImportDesc{
			PkgPath: pkgPath,
			inUse:   false,
		}
		if iSpec.Name != nil {
			id.Alias = iSpec.Name.Name
		}
		result.imports[pkgPath] = id
	}

	return result
}

// ExtractTypeRef returns the way the type should be referenced in code.
func (id ImportHandler) ExtractTypeRef(t types.Type) string {
	// "named" means it is a type which may require importing.
	named, ok := t.(*types.Named)
	if !ok {
		// "*types.Basic"s e.g. string come out as "untyped string"; we need to drop
		//  that part... Not sure why this is how the type information is conveyed :-/.
		return strings.TrimPrefix(t.String(), "untyped ")
	}

	pkg := named.Obj().Pkg()
	typeName := named.Obj().Name()
	if pkg.Path() == id.PInfo.PkgPath {
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
func (id ImportHandler) GetActive() []ImportDesc {
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
func (id ImportHandler) HasActiveImports() bool {
	return len(id.GetActive()) > 0
}
