package main

import (
	"flag"
	"github.com/itzg/go-flagsfiller"
	"log"

	"github.com/drshriveer/gtools/gsort/gen"

	"github.com/drshriveer/gtools/gencommon"
)

const generator = "gsort"

func main() {
	g := gen.Generate{}
	filler := flagsfiller.New()
	if err := filler.Fill(flag.CommandLine, &g); err != nil {
		log.Fatal(err)
	}
	flag.Parse()

	g.InFile = gencommon.SanitizeSourceFile(g.InFile)
	g.OutFile = gencommon.SanitizeOutFile(g.OutFile, g.InFile, generator)

	if len(g.Types) == 0 {
		log.Fatal("type is required")
	}
	log.Printf("%s: %s::%s => %s", generator, g.InFile, g.Types, g.OutFile)

	if err := g.Parse(); err != nil {
		log.Fatalf("parsing failed: %+v", err)
	}

	if err := g.Write(); err != nil {
		log.Fatalf("writing failed: %+v", err)
	}
}
