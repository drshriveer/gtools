# gConfig

gConfig is a simple, light-weight static configuration library that leverages generics and yaml. 

### Getting started

```bash
go get -u github.com/drshriveer/gtools/pkg/gconfig
```

#### Requirements:

- golang 1.19+
- [genum](github.com/drshriveer/gtools/pkg/genum) - for dimensions

### Docs:

https://pkg.go.dev/github.com/drshriveer/gtools/pkg/gconfig

### Features
- __Generics:__ This library uses generics to fetch configuration values from a yaml file. This works with primitives, slices, maps<sup>†</sup>, and structs (supporting yaml). The same key can be resolved into multiple types. _<sup>†</sup> - note: maps with dimensional keys do not currently work_
- __Internal Type Caching:__ After a setting has been parsed into a type it is cached along with that type information for future resolution.   
- __Dimensions:__ A single configuration file may multiple "dimensions" that are resolved at runtime based on program flags to determine the variation of a setting to vend. Differentiating setting variables by environment/stage (e.g. Development, Beta, Prod) is a great example of how this can be leveraged.
  - __Auto-flagging:__ The configuration library will automatically turn dimensions into flags! (unless otherwise specified)
- __Environmental Overrides:__ (TODO) In some cases it is useful to override a single static configuration variable in a specific environment. This can be done though the use of environmental variables.

### Usage

**Define your Dimensions with genum:**

```go
package environment
//go:generate genum -types=Stage

type Stage int

const (
	Development Stage = iota
	Beta
	Prod
)
```
- don't forget to run the generate script ;).

**Define a yaml configuration file:**

```yaml
runtime: 
  max-goroutines: 1_100_000
  request-timeout: 15s

clients:
  redis:
    address:
      Prod: path.to.prod.redis:4090
      Beta: path.to.beta.redis:4090
      default: localhost:4090
    requestTimeout: 5s
    maxTries: 3
  serviceA:
    address: servceA.com:443
    requestTimeout:
        Prod: 2s
        default: 3s
    maxTries: 
      Development: 1
      default: 5
      
```

**Initialize the Config object:**

```go
package main

import (
	_ "embed"
	"github.com/drshriveer/gtools/gconfig"
)

//go:embed path/to/config.yaml
var rawConfig []byte

func main() {
	cfg, err := gconfig.NewBuilder().
		WithDimension("stage", environment.Development).
		FromBytes(rawConfig)
	if err != nil {
		panic(err)
	}
	// ...
}
```

**Extract your configurations!**

```go
func DoSomething(cfg *gconfig.Config) {
	// Fetch individual values:
	maxRoutines := gconfig.MustGet[uint64](cfg, "runtime.max-goroutines")
	reqTimeout := gconfig.MustGet[time.Duration](cfg, "runtime.request-timeout")
	redisAddress := gconfig.MustGet[string](cfg, "clients.redis.address")
	
	// Fetch entire structs:
	type ClientSettings struct { 
		Address        string        `yaml:"address"`
		RequestTimeout time.Duration `yaml:"requestTimeout"`
		MaxTries       int           `yaml:"maxTries"`
    }
	redisCfg := gconfig.MustGet[ClientSettings](cfg, "clients.redis")
	servieACfg := gconfig.MustGet[ClientSettings](cfg, "clients.serviceA")
}
```

### TODO / Missing: 
- Environmental overrides. 
- Support for Maps with keys of dimension types.
