package gencommon

import (
	"fmt"
	"go/types"

	"golang.org/x/tools/go/packages"

	"github.com/drshriveer/gtools/set"
)

// ParseIFaceOption are options for parsing an interface.
type ParseIFaceOption uint

const (
	// IncludePrivate indicates private methods should be included in the parsed interface.
	IncludePrivate ParseIFaceOption = 1 << iota

	// IncludeEmbedded indicates embedded methods should be included in the parsed interface.
	// Note 1: Overloaded methods will be dropped, giving priority to the parent.
	// e.g. given:
	//
	// type A struct { }
	// func (A) Foo() { }
	// func (A) Bar() { }
	//
	// type B struct {}
	// func (B) Foo() { }
	// func (B) Baz() { }
	//
	// generating an interface for C:
	// type C struct { A; B }
	// func (C) Blah() { }
	//
	// will result in:
	// type CIFace interface {
	// 	Bar()
	// 	Baz()
	// 	Blah()
	// }
	//
	// Note 2: This is recursive, so embedded methods of embedded methods will be included.
	IncludeEmbedded
)

// Interface is a parsed interface.
type Interface struct {
	// IsInterface returns false if the actual underlying object is a struct rather than an interface.
	IsInterface bool

	// Comments related to the interface.
	Comments Comments

	// Name of the type (or interface).
	Name string

	// TypeRef is how to reference this interface outside of the current package.
	TypeRef string

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
	options ...ParseIFaceOption,
) (*Interface, error) {
	opts := set.MakeBitSet(options...)
	allPkgs := allpkgs(pkgs)
	pkg, ok := allPkgs.findPKgByName(pkgName)
	if !ok {
		return nil, fmt.Errorf("target %s in package %s not found", target, pkgName)
	}
	return allPkgs.findIFaceByNameInPackage(ih, pkg, target, opts)
}

type allpkgs []*packages.Package

func (pkgs allpkgs) findPKgByName(pkgName string) (*packages.Package, bool) {
	for _, pkg := range pkgs {
		if pkg.PkgPath == pkgName {
			return pkg, true
		}
		// I don't really see why this should be necessary...
		for pkgPath, pkg := range pkg.Imports {
			if pkgPath == pkgName {
				return pkg, true
			}
		}
	}

	return nil, false
}

func (pkgs allpkgs) findIFaceByNameInPackage(
	ih *ImportHandler,
	pkg *packages.Package,
	target string,
	opts set.BitSet[ParseIFaceOption],
) (
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

	return pkgs.namedTypeToInterface(ih, typLayer2, opts), nil
}

func (pkgs allpkgs) namedTypeToInterface(
	ih *ImportHandler,
	t *types.Named,
	opts set.BitSet[ParseIFaceOption],
) *Interface {
	pkg, hasPkg := pkgs.findPKgByName(t.Obj().Pkg().Path())
	var methodz hasMethods = t
	if methodz.NumMethods() == 0 {
		if iface, ok := t.Underlying().(*types.Interface); ok {
			methodz = iface
		}
	}

	result := &Interface{
		Name:        t.Obj().Name(),
		IsInterface: false,
		TypeRef:     ih.ExtractTypeRef(t),
		Methods:     make(Methods, 0, t.NumMethods()),
	}
	if hasPkg {
		result.Comments = CommentsFromObj(pkg, t.Obj().Name())
	}

	for i := 0; i < methodz.NumMethods(); i++ {
		mInfo := methodz.Method(i)
		if opts.Has(IncludePrivate) || mInfo.Exported() {
			method := MethodFromSignature(ih, mInfo.Type().(*types.Signature))
			method.Name = mInfo.Name()
			method.IsExported = mInfo.Exported()
			method.Comments = CommentsFromMethod(pkg, t.Obj().Name(), mInfo.Name())
			result.Methods = append(result.Methods, method)
		}
	}

	if !opts.Has(IncludeEmbedded) {
		return result
	}

	// Nothing embedded, so we're done.
	s, ok := t.Underlying().(*types.Struct)
	if !ok {
		return result
	}

	// This is the way we look for overloaded embedded methods.
	// Since we cannot generate an interface with overloaded methods,
	// we will drop them.
	methodsToAdd := make(map[string]*Method)

	// We will use the methods of the _parent_ type over embedded methods,
	// so if we encounter a method with the same name, we will drop the embedded one.
	ignoreEmbeddedMethodsNamed := make(set.Set[string], len(result.Methods))
	for _, m := range result.Methods {
		ignoreEmbeddedMethodsNamed.Add(m.Name)
	}

	for i := 0; i < s.NumFields(); i++ {
		field := s.Field(i)
		if !field.Embedded() {
			continue
		}
		var embeddedIface *Interface
		switch v := field.Type().(type) {
		case *types.Pointer:
			if named, ok := v.Elem().(*types.Named); ok {
				embeddedIface = pkgs.namedTypeToInterface(ih, named, opts)
			}
		case *types.Named:
			embeddedIface = pkgs.namedTypeToInterface(ih, v, opts)
		default:
			continue
		}
		for _, m := range embeddedIface.Methods {
			if ignoreEmbeddedMethodsNamed.Has(m.Name) {
				continue
			}
			if _, ok := methodsToAdd[m.Name]; ok {
				ignoreEmbeddedMethodsNamed.Add(m.Name)
				delete(methodsToAdd, m.Name)
			} else {
				methodsToAdd[m.Name] = m
			}
		}
	}

	for _, m := range methodsToAdd {
		result.Methods = append(result.Methods, m)
	}

	return result
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

	// so.. i know it is absolutely necessary to do one unwrap sometimes
	// but why it looks infitiely other times I cannot say.
	// I reallly wish go's tooling was better around all these types.
	for i := 0; i < 5; i++ {
		v, ok := t.(hasMethods)
		if ok && v.NumMethods() > 0 {
			return v, true
		}
		t = t.Underlying()
	}
	return nil, false
}
