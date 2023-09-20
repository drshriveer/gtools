package gencommon

import (
	"fmt"
	"go/types"
	"sort"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/drshriveer/gtools/set"
)

// ImportDesc is a description of an import.
type ImportDesc struct {
	Alias   string
	PkgPath string
}

// ImportHandler a collection of import descriptions indexed by package path.
type ImportHandler struct {
	PInfo        *packages.Package
	importsInUse set.Set[string]
}

// CalcImports the imports relevant to a specific package and ImportSpec.
func CalcImports(pkg *packages.Package) ImportHandler {
	result := ImportHandler{
		PInfo:        pkg,
		importsInUse: make(set.Set[string]),
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

	id.importsInUse.Add(pkg.Path())
	importInfo := id.PInfo.Imports[pkg.Path()]
	return fmt.Sprintf("%s.%s", importInfo.Name, typeName)
}

// GetActive returns ordered, active imports.
// Used by templates.
func (id ImportHandler) GetActive() []ImportDesc {
	result := make([]ImportDesc, 0, len(id.PInfo.Imports))
	for pkgPath, pkg := range id.PInfo.Imports {
		if id.importsInUse.Has(pkgPath) {
			result = append(result, ImportDesc{
				Alias:   pkg.Name,
				PkgPath: pkgPath,
			})
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
