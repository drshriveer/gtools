package gencommon

import (
	"errors"
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/packages"
)

var (
	// ErrorInterface defines the error interface as a type for comparison.
	ErrorInterface *types.Interface

	// ContextInterface defines the error interface as a type for comparison.
	ContextInterface *types.Interface
)

func init() {
	var err error
	ContextInterface, err = FindIFaceDef("context", "Context")
	if err != nil {
		panic(err)
	}
	ErrorInterface, err = FindIFaceDef("builtin", "error")
	if err != nil {
		panic(err)
	}
}

// FindIFaceDef finds an interface definition of the given package and type.
// Which can be used in type matching.
func FindIFaceDef(pkgName, typeName string) (*types.Interface, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedDeps | packages.NeedTypesInfo | packages.NeedTypes |
			packages.NeedEmbedPatterns | packages.NeedSyntax,
	}
	pkgs, err := packages.Load(cfg, pkgName)
	if err != nil {
		return nil, err
	}
	if len(pkgs) != 1 {
		return nil, errors.New("did not find exactly one package for " + pkgName)
	}

	typeInfo := pkgs[0].Types.Scope().Lookup(typeName)
	iFace, ok := typeInfo.Type().Underlying().(*types.Interface)

	// this is really annoying. but basically, IF we had no parsing errors,
	// (most likely built-in package), the result returned is probably fine??
	if ok && len(pkgs[0].Errors) == 0 {
		iFace = iFace.Complete()
		return iFace, nil
	}

	for _, f := range pkgs[0].Syntax {
		typeInfo := f.Scope.Lookup(typeName)
		if typeInfo != nil {
			// This isn't very safe... but i guess we do it anyway?
			s := types.ExprString(typeInfo.Decl.(*ast.TypeSpec).Type)
			tv, err := types.Eval(pkgs[0].Fset, nil, typeInfo.Pos(), s)
			if err != nil {
				return nil, err
			}
			iFace, ok := tv.Type.(*types.Interface)
			if !ok {
				return nil, errors.New("type " + typeName + " was not an interface!")
			}
			iFace = iFace.Complete()
			return iFace, nil
		}
	}
	return nil, fmt.Errorf("did not find type %q in package %q", typeName, pkgName)
}
