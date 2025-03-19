gConfig
=======

gConfig is a simple, light-weight static configuration library that leverages generics and yaml.

[Docs](https://pkg.go.dev/github.com/drshriveer/gtools/gconfig)

### Getting started

```bash
go get -u github.com/drshriveer/gtools/gconfig
```

#### Requirements:

-	golang 1.19+
-	[genum](github.com/drshriveer/gtools/genum) - for dimensions

### Features

-	**Generics:** This library uses generics to fetch configuration values from a yaml file. This works with primitives, slices, maps<sup>†</sup>, and structs (supporting yaml). The same key can be resolved into multiple types. *<sup>†</sup> - note: maps with dimensional keys do not currently work*
-	**Internal Type Caching:** After a setting has been parsed into a type it is cached along with that type information for future resolution.  
-	**Dimensions:** A single configuration file may multiple "dimensions" that are resolved at runtime based on program flags to determine the variation of a setting to vend. Differentiating setting variables by environment/stage (e.g. Development, Beta, Prod) is a great example of how this can be leveraged.
	-	**Auto-flagging:** The configuration library will automatically turn dimensions into flags and parse them! (unless otherwise specified)
	-	**Env Parsing:** The configuration library will automatically parse dimensions environment variables.
	-	**GetDimension:** Extract a Dimension value via `gconfig.GetDimension[my.DimensionType](cfg)`.
-	**Template Environmental Variables:** You can reference environmental variables in a config.yaml as values themselves! e.g. `var: ${{env: MY_ENV_VAR}}` will look up the variable `MY_ENV_VAR`.
	-	Defaults to use when an environmental variable is missing using the syntax `{{ env:MY_ENV_VAR | default_value }}`.
	-	Note: any `"` characters will be trimmed from default values... e.g. `{{ env:MY_ENV_VAR | "" }}` would default to an empty string.
-	**Environmental Overrides:** (TODO) In some cases it is useful to override a single static configuration variable in a specific environment. This can be done though the use of environmental variables.

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

-	don't forget to run the generate script ;).

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
    accessToken:
      Development: hard-coded-fake-token-for-local-development-only
      default: ${{env:SECRET_ACCESS_TOKEN}}
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
	// Extract the current dimension value (parsed from flags or environment variables).
	stage := gconfig.GetDimension[environment.Stage](cfg)

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

-	Environmental overrides.
-	Support for Maps with keys of dimension types.
