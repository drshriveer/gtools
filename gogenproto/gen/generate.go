package gen

import (
	"bufio"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/drshriveer/gtools/gencommon"
)

// Logger is the name-spaced logger for this script.
// It is exposed for use in the main package.
var Logger = log.New(log.Writer(), "[gogenproto] ", log.LstdFlags)

// Generate is a simple script for generating proto files with a go:generate directive
// relative to the input directory.
type Generate struct {
	InputDir   string   `aliases:"inputDir" env:"PWD" usage:"path to root directory for proto generation"`
	ProtocPath string   `usage:"path to protoc executable (e.g. /path/to/bin/protoc) defaults to 'protoc' in PATH if not set"`
	Recurse    bool     `default:"false" usage:"generate protos recursively"`
	VTProto    bool     `default:"false" usage:"also generate vtproto"`
	GRPC       bool     `default:"false" usage:"also generate grpc service definitions (experimental)"`
	Include    []string `usage:"comma-separated paths to additional directories to add to the proto include path. You can set an optional Go package mapping by appending a = and the package path, e.g. foo=github.com/foo/bar"`

	// TODO: add flags for other languages, TS, etc.
	// TODO: add NATIVE validation support.
}

// Run runs the generate command.
func (g Generate) Run() error {
	paths, err := g.findProtos(g.InputDir, g.Recurse)
	if err != nil {
		return err
	}
	args := []string{
		"--go_out=.",
		"--go_opt=paths=source_relative",
		"--fatal_warnings",
	}
	if g.VTProto {
		args = append(args,
			"--go-vtproto_out=.",
			"--go-vtproto_opt=paths=source_relative,features=marshal+unmarshal+size+equal+clone+pool",
		)
	}
	if g.GRPC {
		args = append(args,
			"--go-grpc_out=.",
			"--go-grpc_opt=paths=source_relative",
		)
	}
	includePaths := append([]string{g.InputDir}, g.Include...)
	for _, pathAndMaybePkg := range includePaths {
		path, pkgPrefix, hasPkgPrefix := strings.Cut(pathAndMaybePkg, "=")
		includePath, err := filepath.Abs(path)
		if err != nil {
			return err
		}
		args = append(args, "-I="+includePath)
		protoImportPaths, err := g.findProtos(includePath, true)
		if err != nil {
			return err
		}
		for _, path := range protoImportPaths {
			hasGoPackage, err := protoFileHasGoPackage(path)
			if err != nil {
				return err
			}
			if hasGoPackage {
				// Don't make a mapping for this, it already specifies its output package name
				continue
			}

			relPath, err := filepath.Rel(includePath, path)
			if err != nil {
				return err
			}

			var pkg string
			if hasPkgPrefix {
				// Use the explicit package mapping
				pkg = filepath.Join(pkgPrefix, filepath.Dir(relPath))
			} else {
				// Auto detect package name from directory
				pkg, err = gencommon.PackageNameFromPath(filepath.Dir(path))
				if err != nil {
					return err
				}
			}

			mapping := relPath + "=" + pkg

			args = append(args,
				"--go_opt=M"+mapping,
			)
			if g.VTProto {
				args = append(args,
					"--go-vtproto_opt=M"+mapping,
				)
			}
			if g.GRPC {
				args = append(args,
					"--go-grpc_opt=M"+mapping,
				)
			}
		}
	}
	args = append(args, paths...)
	path := "protoc"
	if g.ProtocPath != "" {
		path = g.ProtocPath
	}
	cmd := exec.Command(path, args...)
	cmd.Stdout = logPipe{}
	cmd.Stderr = logPipe{}
	return cmd.Run()
}

func (g Generate) findProtos(dir string, recurse bool) ([]string, error) {
	protoList := []string{}
	err := filepath.WalkDir(dir,
		func(pathname string, d fs.DirEntry, err error) error {
			if err != nil || pathname == "." || pathname == g.InputDir {
				return err
			} else if d.IsDir() && !recurse {
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

// A basic check for whether a proto file has a go_package option declared.
func protoFileHasGoPackage(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// Read line-by-line instead of loading the whole file into memory
	scanner := bufio.NewScanner(f)
	if err != nil {
		return false, err
	}
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "option go_package =") {
			return true, nil
		}
	}
	return false, nil
}

type logPipe struct{}

func (lw logPipe) Write(p []byte) (n int, err error) {
	toLog := strings.TrimSuffix(string(p), "\n")
	Logger.Println(toLog)
	return len(p), nil
}
