package gerrors

import (
	"fmt"
)

type stackType int

const (
	noStack stackType = 0
	// FIXME: sourceOnly could just be 1, but setting it higher means we have a few stack elements
	// to search for a source external to this package.
	// this is perhaps a terrible idea.
	sourceOnly   stackType = 4
	defaultStack stackType = 32
)

const (
	// gerrSkip is the number of lines to skip for a base GError
	gerrSkip = 4
	// factorySkip is the number of stack elements to skip for an err factory.
	factorySkip = 5
)

// Error exposes methods that can be used as an error type.
type Error interface {
	// Error implements the error interface; it returns a formatted message of the style
	// "Name: <name>, DTag: <detailTag>, Src: <source>, Message: <message> \n <stack>
	Error() string

	// Is implements the required errors.Is interface.
	Is(error) bool

	// ExtMsgf returns a copy of the embedded error with diagnostic info and the
	// message extended with additional context.
	ExtMsgf(format string, elems ...any) Error

	// DTagExtMsgf returns a copy of the embedded error with diagnostic info, a detail tag,
	// and the message extended with additional context.
	DTagExtMsgf(detailTag string, format string, elems ...any) Error

	// WithDTag returns a copy of the embedded error with diagnostic info and a detail tag.
	WithDTag(detailTag string) Error

	// ErrMessage returns the error's message.
	ErrMessage() string

	// ErrSource returns the source of the error.
	ErrSource() string

	// ErrName returns the name of the error.
	ErrName() string

	// DTag returns the error detail tag.
	DTag() string

	// FIXME: Maybe this
	Unwrap() error // XXX: what would we unwrap to? a separate unknown source? a factory?

	_embededGError() *GError
}

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

	// srcFactory holds a reference back to the factory error that created this message.
	// This unfortunate wrapping is required for switching.
	srcFactory Error

	// srcError holds a reference back to the original error - this is only populated in
	// case of a Convert() call.
	srcError error

	// extensionString is a Key: Value string that is added to the Error() print out.
	// This string is constructed via the FactoryOf method if there is a struct tag that
	// indicates it should be included.
	extensionString string

	skipLines int
}

func (e *GError) _embededGError() *GError {
	return e
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
		result += "SourceInfo: " + e.Source + separator
	}
	if len(e.extensionString) > 0 {
		result += e.extensionString
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
func (e *GError) Is(err error) bool {
	if e.srcFactory != nil {
		return e.srcFactory.Is(err)
	}
	if e == err ||
		e.srcFactory != nil && e.srcError == err ||
		e.srcError != nil && e.srcError == err {
		return true
	}
	switch v := err.(type) {
	case Error:
		return e.Is(v.Unwrap())
	case Factory:
		return v.Is(e)
	}

	return false
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

func (e *GError) ErrMessage() string {
	return e.Message
}

func (e *GError) ErrSource() string {
	return e.Source
}

func (e *GError) ErrName() string {
	return e.Name
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

	if e.skipLines == 0 {
		clone.stack = makeStack(int(st), gerrSkip)
	} else {
		clone.stack = makeStack(int(st), e.skipLines)
	}

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
	// case GError:
	// 	return &v
	case *GError:
		return v
	}

	clone := e.clone(defaultStack)
	clone.Message += fmt.Sprintf(" originalError: %+v", err)
	e.srcError = e

	return clone
}

func (e *GError) DTag() string {
	return e.detailTag
}

// Unwrap is for unwrapping errors to get to the source.
func (e *GError) Unwrap() error {
	if e.srcFactory != nil {
		return e.srcFactory
	}
	return e
}
