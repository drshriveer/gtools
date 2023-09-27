package main

import (
	"flag"
	"log"
	"os"

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

	for _, s := range os.Environ() {
		println("env: " + s)
	}

	if err := g.Run(); err != nil {
		log.Fatalf("run failed: %+v", err)
	}
}
