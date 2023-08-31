# gEnum
### Getting started

```bash
go install github.com/drshriveer/gscommon/pkg/genum/genum
```

```justfile
generate: _tools-generate
    @go generate ./...

_tools-generate:
    @go install github.com/drshriveer/gcommon/pkg/enum/genum
```

#### Docs

// TODO: do we have to install this package ? 
// can it be referenced another way?


### Features
- **Enum Interface**
- [Traits](#traits) - A way
- **Parsing**
- String generators


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

###### Traits

```go
//go:generate genum -types=Creatures
type Creatures int

const (
    NotCreature                             = Creatures(iota)
    Cat, Cat_NumLegs, Cat_IsMammal          = Creatures(iota), 4, true
    Dog, Dog_NumLegs, Dog_IsMammal          = Creatures(iota), 4, true
    Ant, Ant_NumLegs, Ant_IsMammal          = Creatures(iota), 6, false
    Spider, Spider_NumLegs, Spider_IsMammal = Creatures(iota), 8, false
    Human, Human_NumLegs, Human_IsMammal    = Creatures(iota), 2, true
)
```

Will generate with the following functions, in addition to the basic functions.
```go
func (c Creatures) NumLegs() int { ... }
func (c Creatures) IsMammal() bool { ... }
```

###### Options
