package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
)

// TODO: make this an input flag
const file = "./pkg/tester/TestEnum.go"
const typeName = "MyEnum"

func main() {
	fset := token.NewFileSet()
	fAST, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	// XX: i think I want fAST.Decals
	// I should see one type definition, *ast.GenDecal
	// Then a group of elements ( also *ast.GenDecal ) -> items under ::Specs
	// and the *first* should be the expected type.

	// this should be tested with various blocks of enum declarations, as well
	// as enums declared outside of a block.

	for _, v := range extractFromDeclarations(typeName, fAST.Decls) {
		Println("ExtracteD: ", v)
	}
}

type enumValue struct {
	name  string
	value int
}

func extractFromDeclarations(targetType string, decals []ast.Decl) []enumValue {
	result := make([]enumValue, 0)
	for _, d := range decals {
		if v, ok := d.(*ast.GenDecl); ok {
			result = append(result, extractFromSpec(targetType, v.Specs)...)
		}
	}
	return result
}

func extractFromSpec(targetType string, specs []ast.Spec) []enumValue {
	// cases:
	// 1- purely type spec
	// 2- iota
	// 3- distractions
	inIOTA := false
	result := make([]enumValue, 0)
	for _, s := range specs {
		switch v := s.(type) {
		case *ast.ValueSpec:
			if isIOTAStart(targetType, v) {
				inIOTA = true
			} // maybe this should be an else-if
			if inIOTA {
				result = append(result, extractIOTAValues(v)...)
			}
			result = append(result, extractSoloValues(targetType, v)...)
		}
	}
	return result
}

func isIOTAStart(targetType string, s *ast.ValueSpec) bool {
	if vt, ok := s.Type.(*ast.Ident); !ok || vt.Name != targetType {
		return false
	}
	if len(s.Values) != 1 {
		return false
	}
	if vv, ok := s.Values[0].(*ast.Ident); ok && vv.Name == "iota" {
		return true
	}
	return false
}

func extractIOTAValues(s *ast.ValueSpec) []enumValue {
	if s.Type != nil {
		return nil
	}
	// I think logic here has to check that type == nil && that values == nil
	// if type == nil, && values == Kind != (int) 100% not my enum

	// however! for our "ValueSeven"
	// type == nil  values (int, "7" ((AS STRING))) can exist, WITH names data = 3 (from iota)
	// i'm not actually sure which is correct to follow here.
	// I would expect 7 to be the actual answer here, but data == 3. That reflects what the iota
	// position would be, but is rather confusing.
	// There's gotta ve another way
	return extractResultsFromIdents(s.Names)
}

func extractSoloValues(targetType string, s *ast.ValueSpec) []enumValue {
	if vt, ok := s.Type.(*ast.Ident); !ok || vt.Name != targetType {
		return nil
	}
	return extractResultsFromIdents(s.Names)
}

func extractResultsFromIdents(idents []*ast.Ident) []enumValue {
	result := make([]enumValue, 0)
	for _, ident := range idents {
		result = append(
			result, enumValue{
				name:  ident.Name,
				value: ident.Obj.Data.(int),
			},
		)
	}
	return result
}

func Println(msg string, args ...any) {
	println(fmt.Sprintf(msg, args...))
}
