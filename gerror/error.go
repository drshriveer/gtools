package gerror

// Error exposes methods that can be used as an error type.
type Error interface {
	// Error implements the error interface; it returns a formatted message of the style
	// "Name: <name>, DTag: <detailTag>, Src: <source>, Message: <message> \n <stack>
	Error() string

	// Is implements the errors.Is interface.
	// This works on converted errors to compare against an external source,
	// as well as error factories.
	Is(error) bool

	// Unwrap implements errors.Unwrap interface and works along side errors.Is.
	// It will unwrap to the underlying Factory, NOT to a converted error.
	// Error.Is can still validate converted errors.
	Unwrap() error

	// ErrMessage returns the error's message.
	ErrMessage() string

	// ErrSource returns the source of the error.
	ErrSource() string

	// ErrName returns the name of the error.
	ErrName() string

	// ErrDetailTag returns the detailTag (if set).
	ErrDetailTag() string

	// ErrStack returns an error stack (if available).
	ErrStack() Stack

	_embededGError() *GError
}
