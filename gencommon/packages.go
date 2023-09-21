package gencommon

import (
	"errors"
	"go/ast"
	"path"
	"strings"

	"golang.org/x/tools/go/packages"
)

// LoadPackages is a utility for parsing packages etc of a given file.
func LoadPackages(fileName string) (*packages.Package, *ast.File, *ImportHandler, error) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles |
			packages.NeedImports | packages.NeedDeps | packages.NeedTypesInfo | packages.NeedTypes |
			packages.NeedEmbedPatterns | packages.NeedSyntax,
	}

	pkgs, err := packages.Load(cfg, path.Dir(fileName))
	if err != nil {
		return nil, nil, nil, err
	}

	if len(pkgs) < 1 {
		return nil, nil, nil, errors.New("package for file " + fileName + " NOT FOUND")
	}

	// XXX: this might be wrong; I don't like the idea of picking the first package randomly.
	// might need to search through to find one with the target file instead.
	fAST, err := FindFAST(pkgs[0], fileName)
	if err != nil {
		return nil, nil, nil, err
	}

	return pkgs[0], fAST, calcImports(pkgs[0], fAST), nil
}

// FindFAST finds a * in a package.
func FindFAST(p *packages.Package, fileName string) (*ast.File, error) {
	fileIndex := -1
	cleanFName := path.Clean(fileName)
	for i, fName := range p.GoFiles {
		// Depending on the run context we might have absolute or relative paths...
		// Thus we do matching both ways...
		// Simultaneously overkill and still dangerous, but whatever.
		if fName == fileName || strings.HasSuffix(fName, cleanFName) {
			fileIndex = i
			break
		}
	}
	if fileIndex == -1 {
		return nil, errors.New("Not found")
	}

	return p.Syntax[fileIndex], nil
}
