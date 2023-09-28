package main

import (
	"flag"

	"github.com/itzg/go-flagsfiller"

	"github.com/drshriveer/gtools/gogenproto/gen"
)

func main() {
	g := gen.Generate{}
	filler := flagsfiller.New()
	if err := filler.Fill(flag.CommandLine, &g); err != nil {
		gen.Logger.Fatal(err)
	}
	flag.Parse()

	if err := g.Run(); err != nil {
		gen.Logger.Fatalf("run failed: %+v", err)
	}
}
