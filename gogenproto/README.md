gogenproto
==========

Is a generator that will generate go files from proto in a directory via a directive.

### Getting Started

**Dependencies:**

This repo does not include any proto dependencies intentionally, because you should be in control of your own versions. Including any here would be a rather cumbersome bottleneck and a very strong reason not to use this repo. Therefore, you must bring your own versions of (i.e. they must be discoverable in PATH):

-	`protoc`
	-	hint mac: `brew install protobuf`
-	`protoc-gen-go`
	-	we recommend using a tools.go file to track and install versions of this.
-	`protoc-gen-go-grpc` (if required)
-	`protoc-gen-go-vtproto` (if required)

**Install gogenproto:**

```bash
go install github.com/drshriveer/gtool/gogenproto/cmd/gogenproto@latest
```

##### Usage:

In your working directory...

**Define a proto file:**

`./base/pkg/models/message.proto`

```protobuf
syntax = "proto3";

package base.pkg.models;

option go_package = "models/";

message Message {
  uint64 id = 1;
  string content = 2;
}
```

**Add the generate directive:**

`./base/pkg/models/generate.go`

```go
package models

//go:generate gogenproto -vt-proto
```

##### Options

```bash
→ ./path/to/bin/gogenproto --help
Usage of ./bin/gogenproto:
  -grpc
    	also generate grpc service definitions (experimental)
  -include value
    	comma-separated paths to additional directories to add to the proto include path. You can set an optional Go package mapping by appending a = and the package path, e.g. foo=github.com/foo/bar
  -input-dir string
    	path to root directory for proto generation (env PWD)
  -inputDir string
    	path to root directory for proto generation (env PWD)
  -recurse
    	generate protos recursively
  -vt-proto
    	also generate vtproto
```

### TODO:

-	add (native!) support for proto validation
-	add flags for other languages
