gEnum
=====

gEnum is an enum code generator for golang inspired by projects like [enumer](https://github.com/dmarkham/enumer).

[Docs](https://pkg.go.dev/github.com/drshriveer/gtools/genum)

### Getting started

Install with:

```bash
go install github.com/drshriveer/gtools/genum/cmd/genum@latest
```

Add `go:generate` command directive to file with enum definition:

```go
//go:generate genum -types=EnumType1,EnumType2
```

### Features

-	**Enum Interface** - All enums implement a common interface that can be referenced directly; useful when an enum type is required.
-	**Marshalers** - enums are generated with yaml/v3, json, and text unmarshalers.
-	[Traits](#traits) - Tie constant values to enums as first-class citizens!

##### Generated Methods

Below are all the methods that are generated by default:

```go
func (e MyEnum) IsValid() bool {...}
func (MyEnum) Values() []MyEnum {...}
func (MyEnum) StringValues() []string {...}
func (e MyEnum) String() string {...}
func (e MyEnum) ParseString(text string) (MyEnum, error) {...}
func (e MyEnum) MarshalJSON() ([]byte, error) {...}
func (e *MyEnum) UnmarshalJSON(data []byte) error {...}
func (e MyEnum) MarshalText() ([]byte, error) {...}
func (e *MyEnum) UnmarshalText(text []byte) error {...}
func (e MyEnum) MarshalYAML() (any, error) {...}
func (e *MyEnum) UnmarshalYAML(value *yaml.Node) error {...}
func (MyEnum) IsEnum() {}
```

### Usage

###### Basic

```go
//go:generate genum -types=Creatures
type Creatures int

const (
	NotCreature Creatures = iota
	Cat
	Dog
	Ant
	Spider
	Human
)
```

-	Define the enum (`Creatures` above) in a file (e.g. `filename.go`\)
-	Add the generate directive: `//go:generate genum -types=Creatures`.
-	Run the `go generate` command.
-	Code will be generated and written to file `<filanme>.genum.go` in the same package.

###### Traits

Traits tie other constant values to an enum values. Traits must be defined on the same line as an enum value. The generation code will generate a method for the trait of the TraitName.  
TraitNames are derived from the trait constants of the lowest valued enum. They may be prefixed with `_` (e.g. `_NumLegs`) so that they are not exposed out of the package.

```go
//go:generate genum -types=Creatures
type Creatures int

const (
	NotCreature, _NumLegs, _IsMammal = Creatures(iota), 0, false
	Cat, _, _                        = Creatures(iota), 4, true
	Dog, _, _                        = Creatures(iota), 4, true
	Ant, _, _                        = Creatures(iota), 6, false
	Spider, _, _                     = Creatures(iota), 8, false
	Human, _, _                      = Creatures(iota), 2, true
)
```

Will generate with the following trait methods, in addition to the basic functions.

```go
func (c Creatures) NumLegs() int { ... }
func (c Creatures) IsMammal() bool { ... }
```

Genums can also be parsed by their traits by using the `--parsableByTraits=TraitName1,TraitName2` flag. When using this flag code generation will fail if trait values can be parsed into multiple enums; uniqueness is required. Furthermore, there may be edge cases where traits do not parse consistently between various parsers... Durations for example. We will try to fix these in subsequent updates; if you discover any, please file an issue asap.

###### Duplicate Values

Duplicated enum values present a small challenge to code; it is not always possible to distinguish between identical values. For example, when turning an enum into string form. In such cases the generator will consistently choose one value as the "primary" value. To force a primary value, mark all others as `Deprecated:`.

###### Options

```bash
→ genum --help
Usage of ./bin/genum:
  -disableTraits
        disable trait syntax inspection (default false)
  -in string
        path to input file (defaults to go:generate context)
  -json
        generate json marshal methods (default true) (default true)
  -out string
        name of output file (defaults to go:generate context filename.enum.go)
  -text
        generate text marshal methods (default true) (default true)
  -types string
        [Required] comma-separated names of types to generate enum code for
  -yaml
        generate yaml marshal methods (default true) (default true)
  -caseInsensitive
        string parsing of enum names will be case insensitive (default false)
  -parsableByTraits string
        comma-separated list of trait names which will generate their own parser
```

###### Limitations

1.	Enum definitions must be in a single file.
2.	Currently no string transformation support.
3.	[Duplicate Values](#duplicate-values) can cause some issues; prefer not to use them.
4.	In some cases parsing by traits may have unexpected behaviors. Durations, for example.  

### TODO:

-	generate tests for parsing traits to prevent unintentional bugs and to flag issues that may arise
