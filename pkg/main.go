package main

import (
	"flag"
	"github.com/drshriveer/gsenum/pkg/gen"
	"path"
)

var (
	inFile      = flag.String("in", "", "path to input file")
	outFileName = flag.String("out", "", "name of output file")
	typeName    = flag.String("type", "", "name of type to generate enum code for")
	genJSON     = flag.Bool("json", true, "generate json marshal methods (default true)")
	genYAML     = flag.Bool("yaml", true, "generate yaml marshal methods (default true)")
	genText     = flag.Bool("text", true, "generate text marshal methods (default true)")
)

// TODO: make this an input flag
const _file = "./pkg/tester/TestEnum.go"
const _outFileName = "test_enum.enum.go"
const _typeName = "MyEnum"

func main() {
	flag.Parse()

	if len(*inFile) == 0 {
		*inFile = _file
		//panic("input file is required")
	}
	if len(*outFileName) == 0 {
		*outFileName = _outFileName
		//panic("output file name is required")
	}
	if len(*typeName) == 0 {
		*typeName = _typeName
		//panic("typename is required")
	}

	g := gen.Generate{
		InFile:       *inFile,
		OutFile:      path.Join(path.Dir(*inFile), *outFileName),
		EnumTypeName: *typeName,
		GenJSON:      *genJSON,
		GenYAML:      *genYAML,
		GenText:      *genText,
	}

	if err := g.Parse(); err != nil {
		panic(err)
	}

	if err := g.Write(); err != nil {
		panic(err)
	}
}
