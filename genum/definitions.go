package genum

// EnumLike is a generic type for something that looks like an enum.
type EnumLike interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Enum is the base interface all generated enums implement.
type Enum interface {
	// IsValid returns true if the enum is valid.
	IsValid() bool

	// StringValues returns a list of values as strings.
	StringValues() []string

	// returns the string value of an enum.
	String() string

	// IsEnum does nothing but help define the interface.
	IsEnum()

	// ParseGeneric calls TypedEnum.Parse but returns the result
	// in the generic Enum interface. Which is useful when you are only able to work with
	// the un-typed interface.
	ParseGeneric(input any) (Enum, error)
}

// TypedEnum is extended, generic interface that enums extend.
// I'm not actually sure where this will be useful...
type TypedEnum[T EnumLike] interface {
	Enum

	// Values returns all valid values of an enum.
	Values() []T
}
