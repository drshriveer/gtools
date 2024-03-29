package gencommon

import (
	"errors"
	"fmt"
	"go/ast"
	"path"
	"strings"

	"golang.org/x/tools/go/packages"
)

// LoadPackages is a utility for parsing packages etc of a given file.
func LoadPackages(fileName string, additional ...string) (
	[]*packages.Package, // all
	*packages.Package, // primary
	*ast.File, // primary
	*ImportHandler,
	error,
) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedDeps | packages.NeedTypesInfo | packages.NeedTypes |
			packages.NeedEmbedPatterns | packages.NeedSyntax,
	}

	paths := []string{path.Dir(fileName)}
	paths = append(paths, additional...)

	pkgs, err := packages.Load(cfg, paths...)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	if len(pkgs) < 1 {
		return nil, nil, nil, nil, errors.New("package for file " + fileName + " NOT FOUND")
	}

	pkg, err := FindPackageWithFile(pkgs, fileName)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	fAST, err := FindFAST(pkg, fileName)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return pkgs, pkg, fAST, calcImports(pkg, fAST), nil
}

// FindFAST finds an *ast.File in a package.
func FindFAST(pkg *packages.Package, fileName string) (*ast.File, error) {
	cleanFName := path.Clean(fileName)
	for i, fName := range pkg.GoFiles {
		// Depending on the run context we might have absolute or relative paths...
		// Thus we do matching both ways...
		// Simultaneously overkill and still dangerous, but whatever.
		if fName == fileName || strings.HasSuffix(fName, cleanFName) {
			return pkg.Syntax[i], nil
		}
	}
	return nil, errors.New("fAST for " + fileName + " Not found")
}

// FindPackageWithFile finds a package with a file.
func FindPackageWithFile(pkgs []*packages.Package, fileName string) (*packages.Package, error) {
	cleanFName := path.Clean(fileName)
	for _, pkg := range pkgs {
		for _, fName := range pkg.GoFiles {
			// Depending on the run context we might have absolute or relative paths...
			// Thus we do matching both ways...
			// Simultaneously overkill and still dangerous, but whatever.
			if fName == fileName || strings.HasSuffix(fName, cleanFName) {
				return pkg, nil
			}
		}
	}
	return nil, errors.New("package for " + fileName + " Not found")
}

// PackageNameFromPath returns a fully-qualified package path from a given filename.
// TODO: cache this per directory.
func PackageNameFromPath(fileName string) (string, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles,
	}

	pkgs, err := packages.Load(cfg, fileName)
	if err != nil {
		return "", err
	}
	for _, pkg := range pkgs {
		return pkg.PkgPath, nil
	}
	return "", fmt.Errorf("package for path %s not found", fileName)
}
