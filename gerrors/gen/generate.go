package gen

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"

	"github.com/drshriveer/gtools/gencommon"
	"github.com/drshriveer/gtools/genum/tmpl"
)

// TODO: Gavin!!
// -eat
// - make common package for generation. Include ImportDescs calculation, handling, and template.

// Generate is the parser and writer of gerrors
// It seems to double as its own 'options' holder.
type Generate struct {
	InFile  string   `alias:"in" env:"GOFILE" usage:"path to input file (defaults to go:generate context)"`
	OutFile string   `alias:"out" usage:"name of output file (defaults to go:generate context filename.gerror.go)"`
	Types   []string `usage:"[required] names of types to generate gerrors for"`

	// derived, (exposed for template use):
	Imports gencommon.ImportDescs
	PkgName string
}

// Parse the input file and drives the attributes above.
func (g *Generate) Parse() error {
	fSet := token.NewFileSet()
	fAST, err := parser.ParseFile(fSet, g.InFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}
	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check("", fSet, []*ast.File{fAST}, nil)
	if err != nil {
		return err
	}
	g.Imports = gencommon.CalcImports(fAST.Imports, pkg)

	return nil
}

// Write writes out the enum config file as configured.
func (g *Generate) Write() error {
	return gencommon.Write(tmpl.EnumTemplate, g, g.OutFile)
}
