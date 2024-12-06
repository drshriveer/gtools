package gerror

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
	factoryRef factoryOf

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
	return CloneBase(e, NoStack, "", "", "", nil)
}

// SourceOnly clones the error and ensures Source is populated.
func (e *GError) SourceOnly() Error {
	return CloneBase(e, SourceStack, "", "", "", nil)
}

// Stack clones the error and ensures there is a Stack. Source will also be populated
// if not already set.
func (e *GError) Stack() Error {
	return CloneBase(e, DefaultStack, "", "", "", nil)
}

// Src clones the error with a custom source.
func (e *GError) Src(src string) Error {
	return CloneBase(e, SourceStack, "", src, "", nil)
}

// DTag clones the error with a detailTag, and will populate Source if needed.
func (e *GError) DTag(dTag string) Error {
	return CloneBase(e, SourceStack, dTag, "", "", nil)
}

// Msg clones the error, extends its message, and will populate a Source if needed.
func (e *GError) Msg(format string, elems ...any) Error {
	return CloneBase(e, SourceStack, "", "", fmt.Sprintf(format, elems...), nil)
}

// SrcDTagMsg clones the error, adds a Detail tag, custom source, and extends its message.
func (e *GError) SrcDTagMsg(src, dTag, format string, elems ...any) Error {
	return CloneBase(e, SourceStack, dTag, src, fmt.Sprintf(format, elems...), nil)
}

// SrcDTag clones the error, adds a detail tag and source.
func (e *GError) SrcDTag(src, dTag string) Error {
	return CloneBase(e, SourceStack, dTag, src, "", nil)
}

// SrcMsg clones the error, adds a source, and extends its message.
func (e *GError) SrcMsg(src, format string, elems ...any) Error {
	return CloneBase(e, SourceStack, "", src, fmt.Sprintf(format, elems...), nil)
}

// DTagMsg clones the error, adds a detail tag, and extends its message.
func (e *GError) DTagMsg(dTag, format string, elems ...any) Error {
	return CloneBase(e, SourceStack, dTag, "", fmt.Sprintf(format, elems...), nil)
}

// SrcS is the same as Src but also includes a full StackTrace.
func (e *GError) SrcS(src string) Error {
	return CloneBase(e, DefaultStack, "", src, "", nil)
}

// DTagS is the same as DTag but also includes a full StackTrace.
func (e *GError) DTagS(dTag string) Error {
	return CloneBase(e, DefaultStack, dTag, "", "", nil)
}

// MsgS is the same as Msg but also includes a full StackTrace.
func (e *GError) MsgS(format string, elems ...any) Error {
	return CloneBase(e, DefaultStack, "", "", fmt.Sprintf(format, elems...), nil)
}

// SrcDTagMsgS is the same as DTagSrcMsg but also includes a full StackTrace.
func (e *GError) SrcDTagMsgS(src, dTag, format string, elems ...any) Error {
	return CloneBase(e, DefaultStack, dTag, src, fmt.Sprintf(format, elems...), nil)
}

// SrcDTagS is the same as DTagSrc but also includes a full StackTrace.
func (e *GError) SrcDTagS(src, dTag string) Error {
	return CloneBase(e, DefaultStack, dTag, src, "", nil)
}

// SrcMsgS is the same as SrcMsg but also includes a full StackTrace.
func (e *GError) SrcMsgS(src, format string, elems ...any) Error {
	return CloneBase(e, DefaultStack, "", src, fmt.Sprintf(format, elems...), nil)
}

// DTagMsgS is the same as DTagMsg but also includes a full StackTrace.
func (e *GError) DTagMsgS(dTag, format string, elems ...any) Error {
	return CloneBase(e, DefaultStack, dTag, "", fmt.Sprintf(format, elems...), nil)
}

// Convert will attempt to convert the supplied error into a gError.Error of the
// Factory's type, including the source errors details in the result's error message.
// The original error's equality can be checked with errors.Is().
func (e *GError) Convert(err error) Error {
	if gerr, ok := err.(Error); ok {
		return gerr
	}
	return CloneBase(e, SourceStack, "", "", fmt.Sprintf("originalError: %+v", err), err)
}

// ConvertS is the same as Convert but includes a full StackTrace.
func (e *GError) ConvertS(err error) Error {
	if gerr, ok := err.(Error); ok {
		return gerr
	}
	return CloneBase(e, DefaultStack, "", "", fmt.Sprintf("originalError: %+v", err), err)
}

// Error implements the "error" interface.
func (e *GError) Error() string {
	const separator = ", "
	result := ""
	if e.Name != "" {
		result += "Name: " + e.Name + separator
	}
	if e.detailTag != "" {
		result += "DTag: " + e.detailTag + separator
	}
	if e.Source != "" {
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
func ExtractFactoryReference(err error) Factory {
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
