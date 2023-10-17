package gencommon

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/packages"
	"sort"
)

// Interface is a parsed interface
type Interface struct {
	// IsInterface returns false if the actual underlying object is an struct rather than an interface.
	IsInterface bool
	Name        string
	Methods     Methods
}

// FindInterface locates a given *ast.Interface in a package.
func FindInterface(pkgs []*packages.Package, pkgName, target string) (*Interface, error) {
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

func findIfaceByNameInPackage(pkg *packages.Package, target string) (*Interface, error) {
	for _, file := range pkg.Syntax {
		tt := file.Scope.Lookup(target)
		if tt == nil {
			continue
		}
		ts, ok := tt.Decl.(*ast.TypeSpec)
		if !ok {
			continue
		}

		switch v := ts.Type.(type) {
		case *ast.InterfaceType:
			return interfaceAsInterface(v, target)
		case *ast.StructType:
			return structAsInterface(pkg, target)
		}
	}

	return nil, fmt.Errorf("target %s not found", target)
}

func structAsInterface(pkg *packages.Package, iFaceName string) (*Interface, error) {
	result := &Interface{
		Name:        iFaceName,
		IsInterface: false,
		Methods:     make([]*Method, 0),
	}

	for _, stax := range pkg.Syntax {
		for _, decl := range stax.Decls {
			if v, ok := decl.(*ast.FuncDecl); ok {
				if m, ok := MethodFrom(v, iFaceName); ok {
					result.Methods = append(result.Methods, m)
				}
			}
		}
	}

	sort.Sort(result.Methods)

	return result, nil
}

func interfaceAsInterface(v *ast.InterfaceType, iFaceName string) (*Interface, error) {
	result := &Interface{
		Name:        iFaceName,
		IsInterface: true,
		Methods:     make([]*Method, v.Methods.NumFields()),
	}
	for i, m := range v.Methods.List {
		funcType := m.Type.(*ast.FuncType)
		result.Methods[i] = &Method{
			Name:       m.Names[0].Name,
			IsExported: m.Names[0].IsExported(),
			Comments:   docToString(m.Comment),
			Input:      ParamsFromFieldList(funcType.Params),
			Output:     ParamsFromFieldList(funcType.Results),
		}
	}
	return nil, nil
}

func docToString(group *ast.CommentGroup) []string {
	if group == nil {
		return nil
	}
	return mapper(group.List, func(in *ast.Comment) string { return in.Text })
}

func mapper[Tin any, Tout any](input []Tin, mapFn func(in Tin) Tout) []Tout {
	result := make([]Tout, len(input))
	for i, val := range input {
		result[i] = mapFn(val)
	}
	return result
}
