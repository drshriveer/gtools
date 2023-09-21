package gerrors

import (
	"fmt"
)

// GError is a base error type that can be extended and turned into a factory.
type GError struct {
	// The Name property is the literal name of the error as it will be represented in metrics.
	// Generally, this should match the name of the error variable.
	Name string

	// Message is the base message of an error all errors will have.
	// This message may be extended or modified through various functions.
	// _Never_ programmatically modify this message.
	Message string

	// Source is an optional attribute that identifies where an error originated.
	// If defined in a factory source is static.
	// If not supplied source will be derived unless using a raw error.
	// A derived source includes packageName, typeName (if applicable), and methodName.
	Source string

	// detailTag is a metric-safe 'tag' that can distinguish between different uses of the same error.
	detailTag string

	// stack is the stack trace info.
	stack Stack

	// factoryRef holds a reference back to the factory error that created this message.
	// This unfortunate wrapping is required for switching.
	factoryRef *GError

	// srcError holds a reference back to the original error - this is only populated in
	// case of a Convert() call.
	srcError error

	isFactory bool
}

// ErrMessage returns the error's message.
func (e *GError) ErrMessage() string {
	return e.Message
}

// ErrSource returns the source of the error.
func (e *GError) ErrSource() string {
	return e.Source
}

// ErrName returns the name of the error.
func (e *GError) ErrName() string {
	return e.Name
}

// ErrDetailTag returns the metric-safe detail-tag of the error.
func (e *GError) ErrDetailTag() string {
	return e.detailTag
}

// ErrStack returns an error stack (if available).
func (e *GError) ErrStack() Stack {
	return e.stack
}

// Base clones the base error but does not add any tracing info.
func (e *GError) Base() Error {
	return CloneBase(e, NoStack, "", "", nil)
}

// Convert attempts translates a non-gerror of an unknown kind into this base error.
func (e *GError) Convert(err error) Error {
	if gerr, ok := err.(Error); ok {
		return gerr
	}
	clone := CloneBase(e, DefaultStack, "", fmt.Sprintf("originalError: %+v", err), err)
	return clone
}

// DTag clones the base error and adds or extends a metric tag.
func (e *GError) DTag(dTag string) Error {
	return CloneBase(e, DefaultStack, dTag, "", nil)
}

// ExtMsgf clones the base error and adds an extended message.
func (e *GError) ExtMsgf(format string, elems ...any) Error {
	return CloneBase(e, DefaultStack, "", fmt.Sprintf(format, elems...), nil)
}

// DExtMsgf clones the base error and adds an extended message and metric tag.
func (e *GError) DExtMsgf(dTag string, format string, elems ...any) Error {
	return CloneBase(e, DefaultStack, dTag, fmt.Sprintf(format, elems...), nil)
}

// Src returns a copy of the embedded error with Source populated if needed.
func (e *GError) Src() Error {
	return CloneBase(e, SourceStack, "", "", nil)
}

// Stack is a factory method for cloning the base error with a full sack trace.
func (e *GError) Stack() Error {
	return CloneBase(e, DefaultStack, "", "", nil)
}

// Error implements the "error" interface.
func (e *GError) Error() string {
	const separator = ", "
	result := ""
	if len(e.Name) > 0 {
		result += "Name: " + e.Name + separator
	}
	if len(e.detailTag) > 0 {
		result += "DTag: " + e.detailTag + separator
	}
	if len(e.Source) > 0 {
		result += "Source: " + e.Source + separator
	}
	result += "Message: " + e.Message

	// Note: right now if we have a "source stack", we actually remove the stack after calculations.
	if len(e.stack) > 0 {
		result += "\n" + e.stack.String()
	}

	return result
}

// Unwrap is for unwrapping errors to get to the source.
func (e *GError) Unwrap() error {
	if e.factoryRef != nil {
		return e.factoryRef
	}
	return nil
}

// Is implements the required errors.Is interface.
func (e *GError) Is(err error) bool {
	if e.isFactory && e == ExtractFactoryReference(err) {
		return true
	}
	if e == err ||
		e.factoryRef != nil && e.factoryRef == err ||
		e.srcError != nil && e.srcError == err {
		return true
	}
	gerr, ok := err.(Error)
	if !ok {
		return false
	}
	if unwrapped := gerr.Unwrap(); unwrapped != nil {
		return e.Is(unwrapped)
	}

	return false
}

func (e *GError) _embededGError() *GError {
	return e
}

// ExtractFactoryReference pulls out a factory reference if one exists or returns nil.
func ExtractFactoryReference(err error) *GError {
	gerr, ok := err.(Error)
	if !ok {
		return nil
	}
	embedded := gerr._embededGError()
	if embedded.isFactory {
		return embedded
	}
	return embedded.factoryRef
}
