package gerrors

import (
	"fmt"
)

type stackType int

const (
	noStack stackType = 0
	// FIXME: sourceOnly could just be 1, but setting it higher means we have a few stack elements
	// to search for a source external to this package.
	// this is perhaps a terrible idea... .especially with generated code.
	sourceOnly   stackType = 4
	defaultStack stackType = 32
)

// Error exposes methods that can be used as an error type.
type Error interface {
	// Error implements the error interface; it returns a formatted message of the style
	// "Name: <name>, DTag: <detailTag>, Src: <source>, Message: <message> \n <stack>
	Error() string

	// Is implements the required errors.Is interface.
	Is(error) bool

	// ExtMsgf returns a copy of the underlying error with diagnostic info and the
	// message extended with additional context.
	ExtMsgf(format string, elems ...any) Error

	// DTagExtMsgf returns a copy of the underlying error with diagnostic info, a detail tag,
	// and the message extended with additional context.
	DTagExtMsgf(detailTag string, format string, elems ...any) Error

	// WithDTag returns a copy of the underlying error with diagnostic info and a detail tag.
	WithDTag(detailTag string) Error

	// // Message is the unmodified message string of the error.
	// Message() string
	//
	// // SourceInfo is the unmodified source string of the error.
	// SourceInfo() string
	//
	// // SourceInfo is the unmodified name string of the error.
	// Name() string

	// SourceInfo is the unmodified DetailTag string of the error.
	DTag() string

	// FIXME: Maybe this
	Unwrap() error // XXX: what would we unwrap to? a separate unknown source? a factory?
}

// The Factory interface exposes only methods that can be used for cloning an error.
// But all errors implement this by default.
// This allows for dynamic and mutable errors without modifying the base.
type Factory interface {
	// Factory implements the error interface to permit switching.
	Error() string

	// Base returns a copy of the underlying error without modifications.
	Base() Error

	// WithStack returns a copy of the underlying error with a Stack trace and diagnostic info.
	WithStack() Error

	// WithSource returns a copy of the underlying error with SourceInfo populated if needed.
	WithSource() Error

	// ExtMsgf returns a copy of the underlying error with diagnostic info and the
	// message extended with additional context.
	ExtMsgf(format string, elems ...any) Error

	// DTagExtMsgf returns a copy of the underlying error with diagnostic info, a detail tag,
	// and the message extended with additional context.
	DTagExtMsgf(detailTag string, format string, elems ...any) Error

	// WithDTag returns a copy of the underlying error with diagnostic info and a detail tag.
	WithDTag(detailTag string) Error

	// Convert will attempt to convert the supplied error into a gError.Error of the
	// Factory's type, including the source errors details in the result's error message.
	// The original error can be retrieved via utility methods.
	Convert(err error) Error
}

// GError is a base error type which may represent an actual error or a factory.
// GError itself combines both GError and Factory implementations for two reasons
// - to aid in extending errors
// -
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

	// FIXME: VERY tempting. || RAW
	UseFullSack bool

	// detailTag is a metric-safe 'tag' that can distinguish between different uses of the same error.
	detailTag string

	// stack is the stack trace info.
	stack Stack

	// srcFactory holds a reference back to the factory error that created this message.
	// This unfortunate wrapping is required for switching.
	srcFactory *GError

	// srcError holds a reference back to the original error - this is only populated in
	// case of a Convert() call.
	srcError error
}

// Error implements the "error" interface.
func (e GError) Error() string {
	const separator = ", "
	result := ""
	if len(e.Name) > 0 {
		result += "Name: " + e.Name + separator
	}
	if len(e.detailTag) > 0 {
		result += "DTag: " + e.detailTag + separator
	}
	if len(e.Source) > 0 {
		result += "SourceInfo: " + e.Source + separator
	}
	result += "Message: " + e.Message

	// FIXME: I think I need a check in here to know the difference between stack types.
	// e.g. if there is only one element, I don't think i care about the rest.
	if len(e.stack) > 0 {
		result += "\n" + e.stack.String()
	}

	return result
}

// Is implements the required errors.Is interface.
// FIXME: this is definitely broken.
func (e GError) Is(err error) bool {
	gerr, ok := err.(GError)
	if !ok {
		return false
	}

	// this is a possiblity
	// return e.srcFactory == err.srcFactory
	// or whatever criteria
	return e.Message == gerr.Message && e.Name == gerr.Name
}

// WithStack is a factory method for cloning the base error with a full sack trace.
func (e *GError) WithStack() Error {
	return e.clone(defaultStack)
}

// WithSource is a factory method for cloning the base error and adding source info only (a limited stack trace).
func (e *GError) WithSource() Error {
	return e.clone(sourceOnly)
}

// Base clones the base error but does not add any tracing info.
func (e *GError) Base() Error {
	return e.clone(noStack)
}

// ExtMsgf clones the base error and adds an extended message.
func (e *GError) ExtMsgf(format string, elems ...any) Error {
	clone := e.clone(defaultStack)
	clone.Message = fmt.Sprintf(clone.Message+" "+format, elems...)
	return clone
}

// DTagExtMsgf clones the base error and adds an extended message and metric tag.
func (e *GError) DTagExtMsgf(dTag string, format string, elems ...any) Error {
	clone := e.clone(defaultStack)
	clone.detailTag += "-" + dTag
	clone.Message = fmt.Sprintf(clone.Message+" "+format, elems...)
	return clone
}

// DTag clones the base error and adds a metric tag.
func (e *GError) WithDTag(mTag string) Error {
	clone := e.clone(defaultStack)
	clone.detailTag = mTag
	return clone
}

func (e *GError) clone(st stackType) *GError {
	clone := &GError{
		Name:       e.Name,
		Message:    e.Message,
		Source:     e.Source,
		detailTag:  e.detailTag,
		srcFactory: e,
		stack:      nil,
		srcError:   nil,
	}
	if e.srcFactory != nil {
		clone.srcFactory = e.srcFactory
	}

	if st == noStack {
		return clone
	} else if st == sourceOnly && len(clone.Source) > 0 {
		return clone
	}

	clone.stack = makeStack(int(st), 4)
	// XXX: Fix this: try to generate first stack outside of package.
	clone.Source = (clone.stack)[0].Metric()
	if st == sourceOnly {
		clone.stack = nil
	}
	return clone
}

// Convert attempts translates a non-gerror of an unknown kind into this base error.
func (e *GError) Convert(err error) Error {
	switch v := err.(type) {
	case GError:
		return &v
	case *GError:
		return v
	}

	clone := e.clone(defaultStack)
	clone.Message += fmt.Sprintf(" originalError: %+v", err)
	e.srcError = e

	return clone
}

// Unwrap is for unwrapping errors to get to the source.
func (e *GError) Unwrap() error {
	if e.srcFactory != nil {
		return e.srcFactory
	}
	return e
}

func (e *GError) DTag() string {
	return e.detailTag
}
