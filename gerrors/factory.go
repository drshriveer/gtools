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

// CloneBase is used by factory methods.
func CloneBase[T factoryOf](
	err T,
	stackType StackType,
	stackSkip StackSkip,
	dTag string,
	extMsg string, // PRE FORMATED!
	srcError error, // Be careful with this...
) *GError {
	base := err._embededGError()
	clone := &GError{
		Name:       base.Name,
		Message:    base.Message,
		Source:     base.Source,
		detailTag:  base.detailTag,
		factoryRef: base.factoryRef,
		stack:      base.stack,
		srcError:   base.srcError,
	}

	// handle detail tags:
	if len(dTag) > 0 {
		if len(clone.detailTag) == 0 {
			clone.detailTag = dTag
		} else {
			clone.detailTag += "-" + dTag
		}
	}

	// handle message extension:
	if len(extMsg) > 0 {
		clone.Message += " " + extMsg
	}

	// handle error inheritance:
	if clone.factoryRef == nil && base.isFactory {
		clone.factoryRef = base
	}

	if clone.srcError == nil && srcError != nil {
		clone.srcError = srcError
	}

	// If we already have a stack, don't want one, or want a source and already have it
	// skip stacks.
	if len(clone.stack) > 0 ||
		stackType == NoStack ||
		stackType == SourceStack && len(clone.Source) > 0 {
		return clone
	}

	clone.stack = makeStack(stackType, stackSkip)
	if len(clone.Source) == 0 {
		clone.Source = clone.stack.NearestExternal().Metric()
		if stackType == SourceStack {
			clone.stack = nil
		}
	}

	return clone
}
