# GSort

Is a generator that generates sort interfaces with ordered fields.

##### Usage:

```go
type Category int

//go:generate gsort -types=Sortable
type Sortable struct {
    Category     Category `gsort:"1,String()"`
    Property1    string   `gsort:"2"`
    Property2    int      `gsort:"3"`
    UnsortedProp string
}

```

Generates: 

```go
// GSortSortable implements the sort.Sort interface for Sortable.
type GSortSortable []Sortable

func (s GSortSortable) Len() int {
	return len(s)
}
func (s GSortSortable) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s GSortSortable) Less(i, j int) bool {
	if s[i].Category.String() == s[j].Category.String() {
		if s[i].Property1 == s[j].Property1 {
			return s[i].Property2 < s[j].Property2
		}
		return s[i].Property1 < s[j].Property1
	}
	return s[i].Category.String() < s[j].Category.String()
}

```

### TODO: 
- arrive at a syntax that works to generate varients of sortable data.
- improve documentation.
- 