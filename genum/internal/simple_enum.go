package internal

//go:generate genum -types=MyEnum,MyEnum2,MyEnum3

// MyEnum is a mess of a definition;
// - multiple constants resolve to the same value
// - definitions are spread across multiple blocks.
type MyEnum int

// MyEnum2 is simple, but still a little messy as it is defined in the middle
// of MyEnum.
type MyEnum2 uint64

const (
	// Enum1Value0 is the default value and is completely unset.
	Enum1Value0 MyEnum = iota
	Enum1Value1
	Enum1Value2

	// Enum1Value7 is a special thing
	Enum1Value7 MyEnum = 7

	UnrelatedValue         = "my string!"
	Enum2Value0    MyEnum2 = iota
	Enum2Value1
	Enum1IntentionallyNegative MyEnum = -1
)

// These should just be treated as alternative definitions.
const (
	// Deprecated: old value, don't use
	Enum1Value0Complication1 MyEnum = iota
	Enum1Value1Complication1 MyEnum = iota
	Enum1Value2Complication1 MyEnum = iota
)

// MyEnum3 is a simple, well-formed enum with nothing special.
type MyEnum3 int

const (
	Enum3Value0 MyEnum3 = iota
	Enum3Value1
	Enum3Value2
	Enum3Value3
	Enum3Value4
	Enum3Value5
	Enum3Value6
	Enum3Value7
	Enum3Value8
	Enum3Value9
	Enum3Value10
	Enum3Value11
	Enum3Value12
	Enum3Value13
	_ // Enum3Value14 // intentionally missing!
	Enum3Value15
	Enum3Value16
)
