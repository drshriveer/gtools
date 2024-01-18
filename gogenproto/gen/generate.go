package gen

import (
	"io/fs"
	"log"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

// Logger is the name-spaced logger for this script.
// It is exposed for use in the main package.
var Logger = log.New(log.Writer(), "[gogenproto] ", log.LstdFlags)

// Generate is a simple script for generating proto files with a go:generate directive
// relative to the input directory.
type Generate struct {
	InputDir  string `aliases:"inputDir" env:"PWD" usage:"path to root directory for proto generation"`
	OutputDir string `aliases:"outputDir" default:"../" usage:"relative output path for generated files"`

	Recurse bool     `default:"false" usage:"generate protos recursively"`
	VTProto bool     `default:"false" usage:"also generate vtproto"`
	GRPC    bool     `default:"false" usage:"also generate grpc service definitions (experimental)"`
	Include []string `usage:"comma-separated paths to additional packages to include"`

	// TODO: add flags for other languages, TS, etc,
}

// Run runs the generate command.
func (g Generate) Run() error {
	paths, err := g.findProtos()
	if err != nil {
		return err
	}
	args := []string{
		"--proto_path=" + path.Dir(g.InputDir),
		"--go_out=" + g.OutputDir,
		"--fatal_warnings",
	}
	if g.VTProto {
		args = append(args,
			"--go-vtproto_out="+g.OutputDir,
			"--go-vtproto_opt=paths=source_relative,features=marshal+unmarshal+size+equal+clone+pool",
		)
	}
	if g.GRPC {
		args = append(args,
			"--go-grpc_out="+g.OutputDir,
			"--go-grpc_opt=paths=source_relative",
		)
	}
	for _, include := range g.Include {
		args = append(args, "-I="+filepath.Join(g.InputDir, include))
	}
	args = append(args, paths...)
	cmd := exec.Command("protoc", args...)
	cmd.Stdout = logPipe{}
	cmd.Stderr = logPipe{}
	return cmd.Run()
}

func (g Generate) findProtos() ([]string, error) {
	protoList := []string{}
	err := filepath.WalkDir(g.InputDir,
		func(pathname string, d fs.DirEntry, err error) error {
			if err != nil || pathname == "." || pathname == g.InputDir {
				return err
			} else if d.IsDir() && !g.Recurse {
				return fs.SkipDir
			}
			if d.Type().IsRegular() {
				if filepath.Ext(d.Name()) == ".proto" {
					protoList = append(protoList, pathname)
				}
			}
			return nil
		},
	)
	return protoList, err
}

type logPipe struct{}

func (lw logPipe) Write(p []byte) (n int, err error) {
	toLog := strings.TrimSuffix(string(p), "\n")
	Logger.Println(toLog)
	return len(p), nil
}
