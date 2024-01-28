GSort
=====

Is a generator that generates sort interfaces with ordered fields.

### Features

-	auto generate sortable variations of structs via struct tags
-	support multiple sorters per struct. **Note:** to use this feature you may need to disable static check `SA5008`. i.e. in `.golangci.yaml`:

```yaml
issues:
  exclude-rules:
    # gsort handles duplicate tags
    - linters:
        - staticcheck
      text: "SA5008: duplicate struct tag \"gsort\""
```

### Getting started

Install with:

```bash
go install github.com/drshriveer/gtool/gsort/cmd/gsort@latest
```

##### Usage:

Add `gsort` struct tag(s) in the following format:

```go
type ... struct {
    FiledToSortOn string `gsort:"<required:TypeNameToGenerate>,<required:Priority>,<optional:Accessor>"`
}
```

-	`TypeNameToGenerate` is the type name to use for the generated sortable structure.
	-	Prefix this with an optional `*` to indicate that a pointer to the struct should be generated.
-	`Priority` must be specified as an integer; this indicates the relative sort priority of the field in cases where there are multiple fields to sort on.
-	`Accessor` is an optional attribute that indicates a method of the type to use in the sortable computation.

**Example:**

```go
//go:generate gsort -types=Sortable

type Category int

type Sortable struct {
	Category     Category `gsort:"Sortables,1,String()"`
	Property1    string   `gsort:"Sortables,2"`
	Property2    int      `gsort:"Sortables,3" gsort:"SortOnPriority2UsingPointers,1"`
	UnsortedProp string
}

```

**Generates:**

```go
// SortOnPriority2UsingPointers implements a sort.Sort interface for Sortable.
type SortOnPriority2UsingPointers []*Sortable

func (s SortOnPriority2UsingPointers) Len() int {
	return len(s)
}
func (s SortOnPriority2UsingPointers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortOnPriority2UsingPointers) Less(i, j int) bool {
	return s[i].Property2 < s[j].Property2
}

// Sortables implements a sort.Sort interface for Sortable.
type Sortables []Sortable

func (s Sortables) Len() int {
	return len(s)
}
func (s Sortables) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Sortables) Less(i, j int) bool {
	if s[i].Category.String() == s[j].Category.String() {
		if s[i].Property1 == s[j].Property1 {
			if s[i].Property2 == s[j].Property2 {
				return s[i].property3 < s[j].property3
			}
			return s[i].Property2 < s[j].Property2
		}
		return s[i].Property1 < s[j].Property1
	}
	return s[i].Category.String() < s[j].Category.String()
}

```

### TODO:

-	improve documentation. -
