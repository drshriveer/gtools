package main

import (
	"flag"
	"log"

	"github.com/itzg/go-flagsfiller"

	"github.com/drshriveer/gtools/gencommon"
	"github.com/drshriveer/gtools/genum/gen"
)

func main() {
	g := gen.Generate{}
	filler := flagsfiller.New()
	if err := filler.Fill(flag.CommandLine, &g); err != nil {
		log.Fatal(err)
	}
	flag.Parse()

	g.InFile = gencommon.SanitizeSourceFile(g.InFile)
	g.OutFile = gencommon.SanitizeOutFile(g.OutFile, g.InFile, "genum")

	if len(g.Types) == 0 {
		log.Fatal("type is required")
	}
	log.Printf("genum: %s::%s => %s", g.InFile, g.Types, g.OutFile)

	if err := g.Parse(); err != nil {
		log.Fatalf("parsing failed: %+v", err)
	}

	if err := g.Write(); err != nil {
		log.Fatalf("writing failed: %+v", err)
	}
}
