package internal

//go:generate gensort -types=Sortable
type Sortable struct {
	Category     Category `gsort:"1,String()"`
	Property1    string   `gsort:"2"`
	Property2    int      `gsort:"3"`
	UnsortedProp string
}

// Category exists to test the stringify tag
type Category int

const (
	unset Category = iota
	aCategory
	bCategory
)

func (c Category) String() string {
	switch c {
	case unset:
		return "unset"
	case aCategory:
		return "aCategory"
	case bCategory:
		return "bCategory"
	}
	return "failed"
}
