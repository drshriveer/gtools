package gencommon

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/packages"
)

// FindInterface locates a given *ast.Interface in a package.
func FindInterface(pkgs []*packages.Package, pkgName, target string) (*ast.InterfaceType, error) {
	for _, pkg := range pkgs {
		if pkg.PkgPath == pkgName {
			return findIfaceByNameInPackage(pkg, target)
		}
		for pkgPath, pkg := range pkg.Imports {
			if pkgPath == pkgName {
				return findIfaceByNameInPackage(pkg, target)
			}
		}

	}
	return nil, fmt.Errorf("target %s in package %s not found", target, pkgName)
}

func findIfaceByNameInPackage(pkg *packages.Package, target string) (*ast.InterfaceType, error) {
	for _, file := range pkg.Syntax {
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
		return iface, nil
	}

	return nil, fmt.Errorf("target %s not found", target)
}
