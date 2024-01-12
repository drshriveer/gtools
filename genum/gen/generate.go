package gen

import (
	_ "embed"
	"fmt"
	"go/ast"
	"go/constant"
	"go/types"
	"log"
	"slices"
	"sort"
	"strings"
	"text/template"

	"github.com/drshriveer/gtools/gencommon"
	"github.com/drshriveer/gtools/set"
)

//go:embed enumTemplate.gotmpl
var tmpl string

var intType = types.Universe.Lookup("int").Type()
var stringType = types.Universe.Lookup("string").Type()

// enumTemplate is the base template for an enum.
var enumTemplate = template.Must(template.New("genum").Parse(tmpl))

// Generate is the parser and writer of enums and their generated code.
// It seems to double as its own 'options' holder.
type Generate struct {
	InFile           string   `aliases:"in" env:"GOFILE" usage:"path to input file (defaults to go:generate context)"`
	OutFile          string   `aliases:"out" usage:"name of output file (defaults to go:generate context filename.enum.go)"`
	Types            []string `usage:"[required] comma-separated names of types to generate enum code for"`
	GenJSON          bool     `aliases:"json" default:"true" usage:"generate json marshal methods"`
	GenYAML          bool     `aliases:"yaml" default:"true" usage:"generate yaml marshal methods"`
	GenText          bool     `aliases:"text" default:"true" usage:"generate text marshal methods"`
	DisableTraits    bool     `aliases:"disableTraits" default:"false" usage:"disable trait syntax inspection"`
	CaseInsensitive  bool     `aliases:"caseInsensitive" default:"false" usage:"parsing will be case insensitive"`
	ParsableByTraits []string `aliases:"parsableByTraits" usage:"Comma separated list of trait names which will generate their own parser. This will throw an error if the values of that trait are not unique or the trait does not exist."`

	// derived, (exposed for template use):
	Values  []Values                 `flag:""` // ignore these fields
	Traits  []TraitDescs             `flag:""` // ignore these fields
	Imports *gencommon.ImportHandler `flag:""` // ignore these fields
	PkgName string                   `flag:""` // ignore these fields
}

// Parse the input file and drives the attributes above.
func (g *Generate) Parse() error {
	_, pkg, fAST, importInfo, err := gencommon.LoadPackages(g.InFile)
	if err != nil {
		return err
	}

	g.Imports = importInfo
	g.PkgName = pkg.Name
	pkgScope := pkg.Types.Scope()
	g.Values = make([]Values, len(g.Types))
	g.Traits = make([]TraitDescs, len(g.Types))
	for i, enumType := range g.Types {
		values := make(Values, 0)
		for _, decl := range fAST.Decls {
			if d, ok := decl.(*ast.GenDecl); ok {
				for _, spec := range d.Specs {
					vSpec, ok := spec.(*ast.ValueSpec)
					if !ok || len(vSpec.Names) == 0 {
						continue
					}
					vName := vSpec.Names[0].Name
					v, ok := pkgScope.Lookup(vName).(*types.Const)
					if !ok || (v.Type().String() != enumType && !strings.HasSuffix(v.Type().String(), "."+enumType)) {
						continue
					}
					value, isUint := constant.Uint64Val(v.Val())
					enumValue := Value{
						Name:         vName,
						Value:        value,
						Signed:       !isUint,
						IsDeprecated: isDeprecated(fAST, vName),
						Line:         pkg.Fset.Position(v.Pos()).Line,
						astLine:      vSpec,
					}
					values = append(values, enumValue)
				}
			}
		}
		sort.Sort(values)
		g.Values[i] = values

		if g.DisableTraits || len(values) == 0 {
			continue
		}

		// handle traits next!
		traits, err := g.extractTraitDescs(enumType, pkgScope, values)
		if err != nil {
			return err
		} else if len(traits) == 0 {
			processDuplicates(values, traits, enumType)
			continue
		}

		for i := 1; i < len(values); i++ {
			v := values[i]
			for j := 1; j < len(v.astLine.Values); j++ {
				// the code below attempts to evaluate the actual value
				// in the AST.as a _typed_ variable.
				xprStr := types.ExprString(v.astLine.Values[j])
				tDesc := traits[j-1]
				tDesc.Traits = append(tDesc.Traits, TraitInstance{
					OwningValue:  v,
					variableName: v.astLine.Names[j].Name,
					value:        xprStr,
				})
				sort.Sort(tDesc.Traits)
				traits[j-1] = tDesc
			}
		}

		processDuplicates(values, traits, enumType) // detect and warn duplicates
		err = validateParsableTraits(enumType, traits)
		if err != nil {
			return err
		}

		sort.Sort(traits)
		g.Traits[i] = traits
	}

	return nil
}

// validateParsableTraits returns an error if two instances of a value of a parsable trait map to
// different enums.
// eg:
//
//	type EnumerableWithTraits int
//	const (
//		E1, _Trait1 = EnumerableWithTraits(iota), "val"
//		E2, _ = EnumerableWithTraits(iota), "val"
//	)
//
// This will throw an error because "val" matches E1 and E2
func validateParsableTraits(enumType string, traits TraitDescs) error {
	parsableTraitResults := make(map[string]string)
	for _, trait := range traits {
		if trait.Parsable {
			for _, instance := range trait.Traits {
				if parseTo, ok := parsableTraitResults[instance.value]; ok {
					if parseTo != instance.OwningValue.Name {
						return fmt.Errorf(
							"Enum: %s cannot have parsableTrait %s because trait value %s is "+
								"found in %s and %s. parsableByTrait values must be unique within the enum.",
							enumType, trait.Name, instance.value, parseTo, instance.OwningValue.Name)
					}

				}
				parsableTraitResults[instance.value] = instance.OwningValue.Name
			}

		}
	}
	return nil
}

// extractTraitDescs attempts to extract trait descriptions, and does some (minor) validation in the process.
// TraitDescs come from the first type value of an enum. Generally this is 0, but on occasion it can be
// a negative value...
// Note: this expects values to have been sorted.
func (g *Generate) extractTraitDescs(tName string, pkgScope *types.Scope, values Values) (TraitDescs, error) {
	if len(values) == 0 {
		return nil, nil
	}
	firstV := values[0] //nolint:gosec // just wrong!
	traits := make(TraitDescs, 0, max(len(firstV.astLine.Values)-1, 0))
	for j := 1; j < len(firstV.astLine.Values); j++ {
		name := firstV.astLine.Names[j].Name
		v, ok := pkgScope.Lookup(name).(*types.Const)
		if !ok {
			continue // XXX: consider throwing here.. this is probably an invalid condition.
		}

		typeRef := g.Imports.ExtractTypeRef(v.Type())
		// Trait names can be constants or they can be prefixed with `_`
		// which makes them private to the package.
		traitName := strings.TrimPrefix(name, "_")
		if traitName == "" || traitName == "_" {
			return nil, fmt.Errorf(
				"Enum: %s, value: %s (%d) trait %d has no name that can be converted "+
					"into a trait function; this is a violation of the genum contract for "+
					"traits which expects the first enum (by number) to define trait names. "+
					"If this is unexpected, consider setting the DisableTraits flag.",
				tName, firstV.Name, firstV.Value, j,
			)
		}
		tDesc := TraitDesc{
			Name:     traitName,
			Type:     v.Type(),
			TypeRef:  typeRef,
			Parsable: slices.Contains(g.ParsableByTraits, traitName),
			Traits: []TraitInstance{
				{
					OwningValue:  firstV,
					variableName: name,
					value:        v.Val().ExactString(),
				},
			},
		}
		traits = append(traits, tDesc)
	}

	// now some validation...
	foundWithValidValues := make(set.Set[uint64], len(values))
	for _, v := range values {
		// note: we could do more here by ensuring consistent types
		// however inconsistent types will fail a compiler after generation anyway
		// soo... who cares.
		if len(v.astLine.Values)-1 == len(traits) {
			foundWithValidValues.Add(v.Value)
		}
	}

	for _, v := range values {
		if !foundWithValidValues.Has(v.Value) && len(v.astLine.Values) > 1 {
			if len(traits) == 0 {
				return nil, fmt.Errorf(
					"Enum: %s. value: %s (%d) has invalid trait defintions; were trait names defined?. "+
						"Expected %d traits, found %d without well-defined duplicated value "+
						"with expected number of traits.",
					tName, v.Name, v.Value, len(traits), len(v.astLine.Values)-1)
			}
			return nil, fmt.Errorf(
				"Enum: %s. value: %s (%d) has inconsistent trait defintions. "+
					"Expected %d traits, found %d without well-defined duplicated value "+
					"witth expected number of traits.",
				tName, v.Name, v.Value, len(traits), len(v.astLine.Values)-1)
		}
	}

	return traits, nil
}

// Write writes out the enum config file as configured.
func (g *Generate) Write() error {
	if len(g.Values) == 0 {
		return fmt.Errorf("no values to generate; was generate called?")
	}

	return gencommon.Write(enumTemplate, g, g.OutFile)
}

// processDuplicates prints duplicate warnings and selects the "primary" value(s) of traits.
func processDuplicates(values Values, traits TraitDescs, enumTypeName string) {
	if len(values) == 0 {
		return
	}

	data := make(map[uint64]Values, len(values))
	for _, v := range values {
		duplicates, ok := data[v.Value]
		if ok {
			duplicates = append(duplicates, v)
		} else {
			duplicates = Values{v}
		}
		data[v.Value] = duplicates
	}

	for _, duplicates := range data {
		primary, safe := duplicates.getPrimary()
		if safe {
			continue
		}
		// warn about potentially unsafe duplicates.
		log.Printf("[WARN] - Definitions `%v` of `%s` share the same value `%d`. "+
			"`%s` will be arbitrarily chosen as the primary value when stringifying enums. "+
			"If this is undesirable, please mark values other than the intended primary "+
			"as Deprecated.",
			duplicates.stringList(), enumTypeName, primary.Value, primary.Name)

		// correct any traits.
		for i, td := range traits {
			traits[i].Traits = slices.DeleteFunc(td.Traits, func(t TraitInstance) bool {
				return t.OwningValue.Value == primary.Value && t.OwningValue.Name != primary.Name
			})
		}
	}
	sort.Sort(traits)
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

func (g *Generate) ConvertibleFromString(inputType types.Type) bool {
	return types.ConvertibleTo(stringType, inputType)
}

func (g *Generate) ConvertibleFromInt(inputType types.Type) bool {
	// you can get an int from an untyped int and an untyped string
	// but we don't want them to return true because
	if basicType, ok := inputType.(*types.Basic); ok {
		if basicType.Kind() == types.UntypedInt || basicType.Kind() == types.UntypedString {
			return false
		}
	}
	return types.ConvertibleTo(intType, inputType)
}
