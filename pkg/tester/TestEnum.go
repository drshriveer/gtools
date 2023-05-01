package tester

// MyEnum is pretty awesome
type MyEnum int

// MyEnum2 is dopy
// Deprecated:
type MyEnum2 uint64

const (
	// UNSET is the default value and is completely unset.
	UNSET MyEnum = iota

	// ValueOne does a thing
	ValueOne

	// ValueTwo does a thing
	ValueTwo

	// ValueSeven is a special thing
	ValueSeven MyEnum = 7

	UnrelatedValue           = "my string!"
	EnumTwoValueZero MyEnum2 = iota
	EnumTwoValueOne
	IntentionallyNegative MyEnum = -1
)

// These should just be treated as alternative definitions.
const (
	// Deprecated: old value, don't use
	EnumOneComplicationZero MyEnum = iota
	EnumOneComplicationOne  MyEnum = iota
	EnumOneComplicationTwo  MyEnum = iota
)

// These should just be treated as alternative definitions.
const (
	// Deprecated: old names for the same thing.
	EnumTwoComplicationZero, EnumTwoComplicationOne MyEnum = iota, iota
	EnumTwoComplicationTwo, EnumTwoComplicationThree
)

// These should just be treated as alternative definitions.
const (
	// Deprecated: old names for the same thing.
	EnumThreeComplicationZero, EnumThreeComplicationOne MyEnum = 7, iota
	EnumThreeComplicationTwo, EnumThreeComplicationThree
)
