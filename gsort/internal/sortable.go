package internal

// Sortable is for testing.
//
//go:generate gsort -types Sortable=Sortables SortBool=SortedWithBools
type Sortable struct {
	Category     Category `gsort:"1,String()"`
	Property1    string   `gsort:"2"`
	Property2    int      `gsort:"3"`
	property3    int      `gsort:"4"`
	UnsortedProp string
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
