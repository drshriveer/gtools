package gen

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"sort"
	"strings"

	"github.com/drshriveer/gcommon/pkg/genum/tmpl"
	"github.com/drshriveer/gcommon/pkg/set"
)

type Generate struct {
	InFile        string
	OutFile       string
	EnumTypeNames []string
	GenJSON       bool
	GenYAML       bool
	GenText       bool

	// derived:
	Values  []Values
	Traits  []TraitDescs
	Imports []string
	PkgName string
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

	g.PkgName = pkg.Name()
	g.Values = make([]Values, len(g.EnumTypeNames))
	g.Traits = make([]TraitDescs, len(g.EnumTypeNames))
	pkgScope := pkg.Scope()

	for i, enumType := range g.EnumTypeNames {
		values := make(Values, 0)
		traits := make(map[string]TraitInstances)
		var lastValue Value
		for _, name := range pkgScope.Names() {
			// we only care about constants:
			v, ok := pkgScope.Lookup(name).(*types.Const)
			if !ok {
				continue
			}
			// we only care about the target enum Type.
			if v.Type().String() == enumType {
				value, isUint := constant.Uint64Val(v.Val())
				toAdd := Value{
					Name:         name,
					Value:        value,
					Signed:       !isUint,
					IsDeprecated: isDeprecated(fAST, name),
					Line:         fSet.Position(v.Pos()).Line,
				}
				values = append(values, toAdd)
				lastValue = toAdd
			} else if traitName, ok := lastValue.HasTrait(name, fSet.Position(v.Pos()).Line); ok {
				temp, ok := traits[traitName]
				if !ok {
					temp = make(TraitInstances, 0, 1)
				}
				temp = append(temp, TraitInstance{
					OwningValue:  lastValue,
					VariableName: name,
					Type:         v.Type(),
				})
				traits[traitName] = temp
			}
		}
		sort.Sort(values)
		warnDuplicates(values, enumType) // detect and warn duplicates
		g.Values[i] = values
		g.Traits[i] = traitsFromMap(traits)
	}

	return g.calcImports()
}

func (g *Generate) calcImports() error {
	imports := make(set.Set[string])
	for _, traitGroup := range g.Traits {
		for _, traitDesc := range traitGroup {
			if len(traitDesc.PackageImport) > 0 {
				imports.Add(traitDesc.PackageImport)
			}
		}
	}

	// TODO: We can to package validation here i think...
	// look for conflicting import simple package names....
	// need to think it through

	g.Imports = imports.Slice()
	sort.Strings(g.Imports)
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

func warnDuplicates(values Values, enumTypeName string) {
	if len(values) == 0 {
		return
	}

	var lastVal = values[0].Value
	var duplicates []string
	for _, v := range values {
		if lastVal != v.Value {
			if len(duplicates) > 1 {
				println(
					fmt.Sprintf(
						"[WARN] - Definitions `%v` of `%s` share the same value `%d`. "+
							"`%s` will be arbitarily chosen as the primary value when stringifying enums. "+
							"If this is undesireable, please mark values other than the primary Deprecated.",
						duplicates, enumTypeName, lastVal, duplicates[0],
					),
				)
			}
			// reset
			duplicates = nil
			lastVal = v.Value
		}
		if lastVal == v.Value && !v.IsDeprecated {
			duplicates = append(duplicates, v.Name)
		}
	}
}

func isDeprecated(fAST *ast.File, name string) bool {
	obj := fAST.Scope.Lookup(name)
	spec, ok := obj.Decl.(*ast.ValueSpec)
	if !ok {
		return false
	}
	if spec.Doc == nil {
		return false
	}

	for _, comment := range spec.Doc.List {
		trimmed := strings.TrimPrefix(comment.Text, "//")
		trimmed = strings.TrimSpace(trimmed)
		if strings.HasPrefix(trimmed, "Deprecated:") {
			return true
		}
	}
	return false
}
