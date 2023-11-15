package gencommon

import (
	"errors"
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
		Mode: packages.NeedTypesInfo | packages.NeedTypes,
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
	if !ok {
		return nil, errors.New("type " + typeName + " was not an interface!")
	}
	return iFace, nil
}
