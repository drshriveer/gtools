package gen

import (
	_ "embed"
	"sort"
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
	InFile  string   `alias:"in" env:"GOFILE" usage:"path to input file (defaults to go:generate context)"`
	OutFile string   `alias:"out" usage:"name of output file (defaults to go:generate context filename.gerror.go)"`
	Types   []string `usage:"list of type names to generate sorters for"`

	// TODO: remove this argument in favor of an indicator in the tags themselves.
	UsePointer bool `default:"true" usage:"use pointer to value in slice"`

	// derived, (exposed for template use):
	Imports     *gencommon.ImportHandler `flag:""` // ignore these fields
	SorterDescs SorterDescs              `flag:""` // ignore these fields
	PkgName     string                   `flag:""` // ignore these fields
}

// Parse the input file and drives the attributes above.
func (g *Generate) Parse() error {
	_, pkg, _, imports, err := gencommon.LoadPackages(g.InFile)
	if err != nil {
		return err
	}

	g.Imports = imports
	g.PkgName = pkg.Name
	pkgScope := pkg.Types.Scope()

	g.SorterDescs = make(SorterDescs, 0)
	for _, typeToSort := range g.Types {
		obj := pkgScope.Lookup(typeToSort)
		sortDescs, err := createSorterDesc(obj, typeToSort)
		if err != nil {
			return err
		}
		g.SorterDescs = append(g.SorterDescs, sortDescs...)
	}

	sort.Sort(g.SorterDescs)

	return nil
}

// Write writes out the enum config file as configured.
func (g *Generate) Write() error {
	if len(g.SorterDescs) == 0 {
		return nil
	}
	return gencommon.Write(sortTemplate, g, g.OutFile)
}
