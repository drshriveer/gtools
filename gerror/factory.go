package gerror

// The Factory interface exposes only methods that can be used for cloning an error.
// But all errors implement this by default.
// This allows for dynamic and mutable errors without modifying the base.
type Factory interface {
	// Base clones the error without modifications.
	Base() Error

	// SourceOnly clones the error and ensures Source is populated.
	SourceOnly() Error

	// Stack clones the error and ensures there is a Stack. Source will also be populated
	// if not already set.
	Stack() Error

	// Src clones the error with a custom source.
	Src(src string) Error

	// DTag clones the error with a detailTag, and will populate Source if needed.
	DTag(dTag string) Error

	// Msg clones the error, extends its message, and will populate a Source if needed.
	Msg(fmt string, elems ...any) Error

	// DTagSrcMsg clones the error, adds a Detail tag, custom source, and extends its message.
	SrcDTagMsg(src, dTag, fmt string, elems ...any) Error

	// SrcDTag clones the error, adds a detail tag and source.
	SrcDTag(src, dTag string) Error

	// SrcMsg clones the error, adds a source, and extends its message.
	SrcMsg(src, fmt string, elems ...any) Error

	// DTagSrc clones the error, adds a detail tag, and extends its message.
	DTagMsg(dTag, fmt string, elems ...any) Error

	// SrcS is the same as Src but also includes a full StackTrace.
	SrcS(src string) Error

	// DTagS is the same as DTag but also includes a full StackTrace.
	DTagS(dTag string) Error

	// MsgS is the same as Msg but also includes a full StackTrace.
	MsgS(fmt string, elems ...any) Error

	// SrcDTagMsgS is the same as DTagSrcMsg but also includes a full StackTrace.
	SrcDTagMsgS(src, dTag, fmt string, elems ...any) Error

	// SrcDTagS is the same as DTagSrc but also includes a full StackTrace.
	SrcDTagS(src, dTag string) Error

	// SrcMsgS is the same as SrcMsg but also includes a full StackTrace.
	SrcMsgS(src, fmt string, elems ...any) Error

	// DTagSrcS is the same as DTagMsg but also includes a full StackTrace.
	DTagMsgS(dTag, fmt string, elems ...any) Error

	// Convert will attempt to convert the supplied error into a gError.Error of the
	// Factory's type, including the source errors details in the result's error message.
	// The original error's equality can be checked with errors.Is().
	Convert(err error) Error

	// ConvertS is the same as Convert but includes a full StackTrace.
	ConvertS(err error) Error

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
	dTag string,
	source string,
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

	// handle source:
	if len(source) > 0 && len(clone.Source) == 0 {
		clone.Source = source
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
		clone.factoryRef = err
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

	clone.stack = makeStack(stackType, defaultSkip)
	if len(clone.Source) == 0 {
		clone.Source = clone.stack.NearestExternal().Metric()
		if stackType == SourceStack {
			clone.stack = nil
		}
	}

	return clone
}
