package main

import (
	"flag"
	"log"
	"os"
	"path"
	"strings"

	"github.com/drshriveer/gcommon/pkg/genum/gen"
)

var (
	typeNames     = flag.String("types", "", "[Required] comma-separated names of types to generate enum code for")
	inFileArg     = flag.String("in", "", "path to input file (defaults to go:generate context)")
	outFileName   = flag.String("out", "", "name of output file (defaults to go:generate context filename.enum.go)")
	genJSON       = flag.Bool("json", true, "generate json marshal methods (default true)")
	genYAML       = flag.Bool("yaml", true, "generate yaml marshal methods (default true)")
	genText       = flag.Bool("text", true, "generate text marshal methods (default true)")
	disableTraits = flag.Bool("disableTraits", false, "disable trait syntax inspection (default false)")
)

func main() {
	flag.Parse()

	gofile := os.Getenv("GOFILE")
	pwd := os.Getenv("PWD")
	inFile := path.Join(pwd, gofile)
	if len(gofile) == 0 && len(*inFileArg) == 0 {
		log.Fatal("this command should be run in a go:generate context or with -in file set")
	} else if len(gofile) == 0 {
		inFile = inFile
	}

	outFile := path.Join(pwd, strings.TrimSuffix(gofile, ".go")+".genum.go")
	if len(gofile) == 0 && len(*outFileName) == 0 {
		log.Fatal("this command should be run in a go:generate context or with -out file set")
	} else if len(gofile) == 0 {
		outFile = path.Join(path.Dir(inFile), *outFileName)
	}
	if len(*typeNames) == 0 {
		log.Fatal("type is required")
	}
	log.Printf("genum: %s::%s => %s", inFile, *typeNames, outFile)

	g := gen.Generate{
		InFile:        inFile,
		OutFile:       outFile,
		EnumTypeNames: strings.Split(*typeNames, ","),
		GenJSON:       *genJSON,
		GenYAML:       *genYAML,
		GenText:       *genText,
		DisableTraits: *disableTraits,
	}

	if err := g.Parse(); err != nil {
		log.Fatalf("parsing failed: %+v", err)
	}

	if err := g.Write(); err != nil {
		log.Fatalf("writing failed: %+v", err)
	}
}
