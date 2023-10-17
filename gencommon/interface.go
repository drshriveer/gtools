package gencommon

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"
)

// ErrorInterface defines the error interface as a type for comparison.
var ErrorInterface = types.NewInterfaceType([]*types.Func{
	types.NewFunc(
		0,
		nil,
		"Error",
		types.NewSignatureType(nil, nil, nil, nil,
			types.NewTuple(types.NewVar(0, nil, "", types.Typ[types.String])),
			false)),
}, nil)

// Interface is a parsed interface.
type Interface struct {
	// IsInterface returns false if the actual underlying object is a struct rather than an interface.
	IsInterface bool

	// Name of the type (or interface).
	Name string

	// List of methods!
	Methods Methods
}

// FindInterface locates a given *ast.Interface in a package.
func FindInterface(
	ih *ImportHandler,
	pkgs []*packages.Package,
	pkgName, target string,
	includePrivate bool,
) (*Interface, error) {
	for _, pkg := range pkgs {
		if pkg.PkgPath == pkgName {
			return findIFaceByNameInPackage(ih, pkg, target, includePrivate)
		}
		// I don't really see why this should be necessary...
		for pkgPath, pkg := range pkg.Imports {
			if pkgPath == pkgName {
				return findIFaceByNameInPackage(ih, pkg, target, includePrivate)
			}
		}

	}
	return nil, fmt.Errorf("target %s in package %s not found", target, pkgName)
}

func findIFaceByNameInPackage(ih *ImportHandler, pkg *packages.Package, target string, includePrivate bool) (
	*Interface,
	error,
) {
	typ := pkg.Types.Scope().Lookup(target)
	if typ == nil {
		return nil, fmt.Errorf("target %s not found", target)
	}
	typLayer1, ok := typ.(*types.TypeName)
	if !ok {
		return nil, fmt.Errorf("target %s found but not a handled type (found %T)", target, typ)
	}
	typLayer2, ok := typLayer1.Type().(*types.Named)
	if !ok {
		return nil, fmt.Errorf("target %s found but not a handled nested type (found %T)", target, typLayer1)
	}

	return namedTypeToInterface(ih, pkg, typLayer2, includePrivate)
}

func namedTypeToInterface(ih *ImportHandler, pkg *packages.Package, t *types.Named, includePrivate bool) (
	*Interface,
	error,
) {
	result := &Interface{
		Name:        t.Obj().Name(),
		IsInterface: false,
		Methods:     make(Methods, 0, t.NumMethods()),
	}

	for i := 0; i < t.NumMethods(); i++ {
		mInfo := t.Method(i)
		if includePrivate || mInfo.Exported() {
			method := MethodFromSignature(ih, mInfo.Type().(*types.Signature))
			method.Name = mInfo.Name()
			method.IsExported = mInfo.Exported()
			method.Comments = CommentsFromMethod(pkg, t.Obj().Name(), mInfo.Name())
			result.Methods = append(result.Methods, method)
		}
	}
	return result, nil
}

func mapper[Tin any, Tout any](input []Tin, mapFn func(in Tin) Tout) []Tout {
	result := make([]Tout, len(input))
	for i, val := range input {
		result[i] = mapFn(val)
	}
	return result
}
