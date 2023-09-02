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
	Imports ImportDescs
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

	if err := g.calcInitialImports(fAST.Imports, pkg); err != nil {
		return err
	}

	g.PkgName = pkg.Name()
	g.Values = make([]Values, len(g.EnumTypeNames))
	g.Traits = make([]TraitDescs, len(g.EnumTypeNames))
	pkgScope := pkg.Scope()
	for i, enumType := range g.EnumTypeNames {
		values := make(Values, 0)
		traits := make(TraitDescs, 0)
		for _, decl := range fAST.Decls {
			switch d := decl.(type) {
			case *ast.GenDecl:
				for _, spec := range d.Specs {
					vSpec, ok := spec.(*ast.ValueSpec)
					if !ok || len(vSpec.Names) == 0 {
						continue
					}
					vName := vSpec.Names[0].Name
					v, ok := pkgScope.Lookup(vName).(*types.Const)
					if !ok || v.Type().String() != enumType {
						continue
					}
					value, isUint := constant.Uint64Val(v.Val())
					enumValue := Value{
						Name:         vName,
						Value:        value,
						Signed:       !isUint,
						IsDeprecated: isDeprecated(fAST, vName),
						Line:         fSet.Position(v.Pos()).Line,
					}
					values = append(values, enumValue)

					// Handle traits next:
					if len(vSpec.Values) > 1 {
						// if value == 0, create traits and their defaults
						if value == 0 {
							for j := 1; j < len(vSpec.Values); j++ {
								// FIXME: Gavin!! decide if we want to require `_` prefix for traits.
								name := vSpec.Names[j].Name
								v, ok := pkgScope.Lookup(name).(*types.Const)
								if !ok {
									continue
								}
								typeRef := g.Imports.extractTypeRef(v.Type())
								tDesc := TraitDesc{
									Name:    strings.TrimPrefix(name, "_"),
									TypeRef: typeRef,
									Traits: []TraitInstance{
										{
											OwningValue:  enumValue,
											variableName: name,
											value:        v.Val().ExactString(),
										},
									},
								}
								traits = append(traits, tDesc)
							}
						} else if len(traits) != len(vSpec.Values)-1 {
							// FIXME: Gavin!! improve this error message
							return fmt.Errorf("inconsistent number of traits")
						} else {
							for j := 1; j < len(vSpec.Values); j++ {
								// the code below attempts to evaluate the actual value
								// in the AST.as a _typed_ variable.
								xprStr := types.ExprString(vSpec.Values[j])
								tDesc := traits[j-1]
								tDesc.Traits = append(tDesc.Traits, TraitInstance{
									OwningValue:  enumValue,
									variableName: vSpec.Names[j].Name,
									value:        xprStr,
								})
								sort.Sort(tDesc.Traits)
								traits[j-1] = tDesc
							}
						}
					}
				}
			}
		}
		sort.Sort(values)
		warnDuplicates(values, enumType) // detect and warn duplicates
		g.Values[i] = values
		sort.Sort(traits)
		g.Traits[i] = traits
	}

	// FIXME: Gavin! maybe check import alias conflicts here and re-assign traits if needed.

	return nil
}

func (g *Generate) calcInitialImports(importSpecs []*ast.ImportSpec, pkg *types.Package) error {
	g.Imports = ImportDescs{
		currentPackage: pkg,
		imports:        make(map[string]*ImportDesc, len(importSpecs)),
	}
	for _, iSpec := range importSpecs {
		pkgPath := strings.Trim(iSpec.Path.Value, `"`)
		g.Imports.imports[pkgPath] = &ImportDesc{
			Alias:   iSpec.Name.Name,
			PkgPath: pkgPath,
			inUse:   false,
		}
	}
	// FIXME: consider adding json types.
	// This might be a bit brittle because it requires the template and code
	// to be in lock step. but it ensures we don't have odd duplicates.
	// if g.GenJSON {
	// 	g.Imports.imports["encoding/json"] = &ImportDesc{
	// 		Alias:   "",
	// 		PkgPath: "",
	// 		inUse:   true,
	// 	}
	// }
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
