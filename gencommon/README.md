GenGen
======

A generator that generates generator boilerplate via struct annotations. Hopefully, over time, this will become a repository for tools related to generations. For example, common recipes or resolvers; import handlers.

### Getting Started

### Features:

### Usage:

```go
//go:generate gengen -type=Generator
type Generator struct {
	InFile  string   `gengen:"inFile,optional"`
	OutFile string   `gengen:"outFile,optional"`
	Types   []string `gengen:"inFile"`

	Option1 bool `gengen:"option1"`
}
```

Generates a main file with arguments

```go
flag
```
