package enum

// Enum is the base interface that all generated enums will implement.
type Enum[T ~int | ~uint] interface {

	// IsValid returns true if the num is valid.
	IsValid(T) bool

	// Values returns all valid values of an enum.
	Values() []T

	// returns the string value of an enum.
	String() string

	// ParseString converts text into a type if valid.
	// returns true if the enum is valid, and false otherwise.
	ParseString(text string) (T, bool)
}
