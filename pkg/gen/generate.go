package gen

import (
	"fmt"
	"github.com/drshriveer/gsenum/pkg/tmpl"
	"go/ast"
	"go/constant"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"sort"
)

type Generate struct {
	InFile       string
	OutFile      string
	EnumTypeName string
	GenJSON      bool
	GenYAML      bool
	GenText      bool

	// derived:
	Values  Values
	PkgName string

	//EnumIntType  string
	//File fs.File
}

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

	// TODO: extract the underlying enum type (if required)
	g.PkgName = pkg.Name()

	g.Values = make(Values, 0)

	pkgScope := pkg.Scope()
	for _, name := range pkgScope.Names() {
		// we only care about constants:
		v, ok := pkgScope.Lookup(name).(*types.Const)
		if !ok {
			continue
		}
		// we only care about the target enum Type.
		if v.Type().String() != g.EnumTypeName {
			continue
		}
		value, isUint := constant.Uint64Val(v.Val())
		toAdd := Value{
			Name:   name,
			Value:  value,
			Signed: !isUint,
		}
		g.Values = append(g.Values, toAdd)
	}
	sort.Sort(g.Values)
	g.warnDuplicates() // detect and warn duplicates
	return nil
}

func (g *Generate) Write() error {
	if len(g.Values) == 0 {
		return fmt.Errorf("no values to generate; was generate called?")
	}
	f, err := os.OpenFile(g.OutFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.EnumTemplate.Execute(f, g)
}

func (g *Generate) warnDuplicates() {
	if len(g.Values) == 0 {
		return
	}

	var lastVal = g.Values[0].Value
	var duplicates []string
	for _, v := range g.Values {
		if lastVal != v.Value {
			if len(duplicates) > 1 {
				println(
					fmt.Sprintf(
						"[WARN] - Definitions `%v` share the same value `%d`. "+
							"`%s` will be arbitarily chosen as the primary value when stringifying enums. "+
							"If this is undesireable, please mark a value as primary using <FIXME>",
						duplicates, lastVal, duplicates[0],
					),
				)
			}
			// reset
			duplicates = nil
			lastVal = v.Value
		}
		if lastVal == v.Value {
			duplicates = append(duplicates, v.Name)
		}
	}
}
