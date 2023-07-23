package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/drshriveer/gcommon/pkg/enum/gen"
	"github.com/drshriveer/gcommon/pkg/error"
)

var (
	typeNames   = flag.String("types", "", "[Required] comma-separated names of types to generate enum code for")
	inFileArg   = flag.String("in", "", "path to input file (defaults to go:generate context)")
	outFileName = flag.String("out", "", "name of output file (defaults to go:generate context filename.enum.go)")
	genJSON     = flag.Bool("json", true, "generate json marshal methods (default true)")
	genYAML     = flag.Bool("yaml", true, "generate yaml marshal methods (default true)")
	genText     = flag.Bool("text", true, "generate text marshal methods (default true)")
)

type CMDError struct {
	error.GError
	exitCode int
}

func main() {
	flag.Parse()

	gofile := os.Getenv("GOFILE")
	pwd := os.Getenv("PWD")
	inFile := path.Join(pwd, gofile)
	if len(gofile) == 0 && len(*inFileArg) == 0 {
		println("this command should be run in a go:generate context or with -in file set")
		os.Exit(2)
	} else if len(gofile) == 0 {
		inFile = inFile
	}

	outFile := path.Join(pwd, strings.TrimSuffix(gofile, ".go")+".genum.go")
	if len(gofile) == 0 && len(*outFileName) == 0 {
		println("this command should be run in a go:generate context or with -out file set")
		os.Exit(2)
	} else if len(gofile) == 0 {
		outFile = path.Join(path.Dir(inFile), *outFileName)
	}
	if len(*typeNames) == 0 {
		println("type is required")
		os.Exit(2)
	}
	println(fmt.Sprintf("genum: %s::%s => %s", inFile, *typeNames, outFile))

	g := gen.Generate{
		InFile:        inFile,
		OutFile:       outFile,
		EnumTypeNames: strings.Split(*typeNames, ","),
		GenJSON:       *genJSON,
		GenYAML:       *genYAML,
		GenText:       *genText,
	}

	if err := g.Parse(); err != nil {
		panic(err)
	}

	if err := g.Write(); err != nil {
		panic(err)
	}
}
