package gen

import (
	_ "embed"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"text/template"

	"github.com/drshriveer/gtools/gencommon"
)

var (
	//go:embed gsort.gotmpl
	rawSortTemplate string
	sortTemplate    = template.Must(template.New("gsort").Parse(rawSortTemplate))
)

// Generate is the parser and writer of sorters
// It seems to double as its own 'options' holder.
type Generate struct {
	InFile  string            `alias:"in" env:"GOFILE" usage:"path to input file (defaults to go:generate context)"`
	OutFile string            `alias:"out" usage:"name of output file (defaults to go:generate context filename.gerror.go)"`
	Types   map[string]string `usage:"[required] mapping of type names to generate sorters for to name to use for the generated type"`
	// derived, (exposed for template use):
	Imports     gencommon.ImportDescs
	SorterDescs []*SorterDesc
	PkgName     string
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
	g.PkgName = pkg.Name()
	pkgScope := pkg.Scope()

	g.SorterDescs = make([]*SorterDesc, 0)
	for typeToSort, nameOfResultType := range g.Types {
		obj := pkgScope.Lookup(typeToSort)
		sortDesc, err := createSorterDesc(obj, typeToSort, nameOfResultType)
		if err != nil {
			return err
		}
		g.SorterDescs = append(g.SorterDescs, sortDesc)
	}

	return nil
}

// Write writes out the enum config file as configured.
func (g *Generate) Write() error {
	return gencommon.Write(sortTemplate, g, g.OutFile)
}
