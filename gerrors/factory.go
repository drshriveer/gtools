package gerrors

// The Factory interface exposes only methods that can be used for cloning an error.
// But all errors implement this by default.
// This allows for dynamic and mutable errors without modifying the base.
type Factory interface {
	// Base returns a copy of the embedded error without modifications.
	Base() Error

	// Src returns a copy of the embedded error with Source populated if needed.
	// Source is a limited stack.
	Src() Error

	// Stack returns a copy of the embedded error with a Stack trace and diagnostic info.
	Stack() Error

	// ExtMsgf returns a copy of the embedded error with diagnostic info and the
	// message extended with additional context.
	ExtMsgf(format string, elems ...any) Error

	// DExtMsgf returns a copy of the embedded error with diagnostic info, a detail tag,
	// and the message extended with additional context.
	DExtMsgf(detailTag string, format string, elems ...any) Error

	// DTag returns a copy of the embedded error with diagnostic info and a detail tag.
	DTag(detailTag string) Error

	// Convert will attempt to convert the supplied error into a gError.Error of the
	// Factory's type, including the source errors details in the result's error message.
	// The original error can be retrieved via utility methods.
	Convert(err error) Error

	// Error implements the standard Error interface so that a Factory
	// can be passed into errors.Is() as a target.
	Error() string

	// Is implements the interface for error matching in the standard package (errors.IS).
	Is(error) bool
}

type factoryOf interface {
	Error
	Factory
}

// FactoryOf returns a factory instance of this error.
func FactoryOf[T factoryOf](err T) Factory {
	err._embededGError().isFactory = true
	return err
}
