package gerrors

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
	// gerrSkip is the number of lines to skip for a base GError.
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
