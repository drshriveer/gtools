package gen

import (
	_ "embed"
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
	InFile     string            `alias:"in" env:"GOFILE" usage:"path to input file (defaults to go:generate context)"`
	OutFile    string            `alias:"out" usage:"name of output file (defaults to go:generate context filename.gerror.go)"`
	Types      map[string]string `usage:"[required] mapping of type names to generate sorters for to name to use for the generated type"`
	UsePointer bool              `default:"true" usage:"use pointer to value in slice"`

	// derived, (exposed for template use):
	Imports     *gencommon.ImportHandler `flag:""` // ignore these fields
	SorterDescs []*SorterDesc            `flag:""` // ignore these fields
	PkgName     string                   `flag:""` // ignore these fields
}

// Parse the input file and drives the attributes above.
func (g *Generate) Parse() error {
	println("infil!", g.InFile)
	_, pkg, _, imports, err := gencommon.LoadPackages(g.InFile)
	if err != nil {
		return err
	}

	g.Imports = imports
	g.PkgName = pkg.Name
	pkgScope := pkg.Types.Scope()

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
