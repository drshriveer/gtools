package gencommon

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

// FindInterface locates a given *ast.Interface in a package.
func FindInterface(p *packages.Package, target string) *ast.InterfaceType {
	for _, file := range p.Syntax {
		tt := file.Scope.Lookup(target)
		if tt == nil {
			continue
		}
		ts, ok := tt.Decl.(*ast.TypeSpec)
		if !ok {
			continue
		}
		iface, ok := ts.Type.(*ast.InterfaceType)
		if !ok {
			continue
		}
		return iface
	}
	return nil
}
