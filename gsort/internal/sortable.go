package internal

// Sortable is for testing.
//
//go:generate gsort -types Sortable,MultiSort
type Sortable struct {
	Category     Category `gsort:"Sortables,1,String()"`
	Property1    string   `gsort:"Sortables,2"`
	Property2    int      `gsort:"Sortables,3" gsort:"*SortOnPriority2,1"`
	property3    int      `gsort:"Sortables,4"`
	UnsortedProp string
}

// MultiSort is for testing.
type MultiSort struct {
	Property1 string `gsort:"SortByProp1,1" gsort:"SortByProp2,2"`
	Property2 string `gsort:"SortByProp2,1"`
}

// NotSortable is for testing.
type NotSortable struct {
	Prop1 string
}

// SortBool is for testing.
type SortBool struct {
	Category  bool   `gsort:"1"`
	Property1 string `gsort:"2"`
}

// Category exists to test the stringify tag.
type Category int

// These are just test enums.
const (
	Unset Category = iota
	ACategory
	BCategory
	CCategory
)

func (c Category) String() string {
	switch c {
	case Unset:
		return "Unset"
	case ACategory:
		return "ACategory"
	case BCategory:
		return "BCategory"
	case CCategory:
		return "CCategory"
	}
	return "failed"
}
