package gencommon

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"
)

// Interface is a parsed interface.
type Interface struct {
	// IsInterface returns false if the actual underlying object is a struct rather than an interface.
	IsInterface bool

	// Comments related to the interface.
	Comments Comments

	// Name of the type (or interface).
	Name string

	// List of methods!
	Methods Methods
}

// utility helper for various things.
type hasMethods interface {
	NumMethods() int
	Method(i int) *types.Func
}

// ModErrorRefs modifies error references in method returns.
// Swapping out the inner type reference for the type supplied.
// This will return an empty string for ease of calling from inside templates.
// ...use with care.
func (i *Interface) ModErrorRefs(newRef string) string {
	for _, m := range i.Methods {
		if m.ReturnsError() {
			last := m.Output[len(m.Output)-1]
			last.TypeRef = newRef
		}
	}
	return ""
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

	var methodz hasMethods = t
	if methodz.NumMethods() == 0 {
		if iface, ok := t.Underlying().(*types.Interface); ok {
			methodz = iface
		}
	}

	result := &Interface{
		Name:        t.Obj().Name(),
		IsInterface: false,
		Comments:    CommentsFromObj(pkg, t.Obj().Name()),
		Methods:     make(Methods, 0, t.NumMethods()),
	}

	for i := 0; i < methodz.NumMethods(); i++ {
		mInfo := methodz.Method(i)
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

// TypeImplements seems to work where types.Implements does not.
func TypeImplements(aType types.Type, target *types.Interface) bool {
	a, ok := unwrapToHasMethods(aType)
	if !ok {
		// just kuz i guess?
		return types.Implements(aType, target)
	}

	targetMethods := make(map[string]*types.Func, target.NumMethods())
	for i := 0; i < target.NumMethods(); i++ {
		mInfo := target.Method(i)
		targetMethods[mInfo.Name()] = mInfo
	}

	for i := 0; i < a.NumMethods(); i++ {
		mA := a.Method(i)
		mB, ok := targetMethods[mA.Name()]
		if !ok {
			continue
		}
		sigA := mA.Type().(*types.Signature)
		sigB := mB.Type().(*types.Signature)
		if !IsSameSignature(sigA, sigB) {
			return false
		}
		delete(targetMethods, mA.Name())
	}
	return len(targetMethods) == 0
}

// IsSameSignature returns true if signature is essentially the same.
// It does this *with out* checking receivers.
func IsSameSignature(a, b *types.Signature) bool {
	// check inputs:
	if a.Variadic() != b.Variadic() {
		return false
	}
	if !IsTupleSame(a.Params(), b.Params()) {
		return false
	}
	if !IsTupleSame(a.Results(), b.Results()) {
		return false
	}
	return true
}

// IsTupleSame checks if two tuples share the same order of types.
func IsTupleSame(a, b *types.Tuple) bool {
	if a.Len() != b.Len() {
		return false
	}
	for i := 0; i < a.Len(); i++ {
		aParam, bParam := a.At(i), b.At(i)
		if aParam.Type().String() != bParam.Type().String() {
			return false
		}
	}
	return true
}

func unwrapToHasMethods(t types.Type) (hasMethods, bool) {
	if t == nil {
		return nil, false
	}
	v, ok := t.(hasMethods)
	if !ok || v.NumMethods() == 0 {
		return unwrapToHasMethods(t.Underlying())
	}
	return v, true
}
