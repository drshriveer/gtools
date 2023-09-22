package gerror

// ErrUnknown converts any non-gError into a gerror.
var ErrUnknown = FactoryOf(&GError{
	Name:    "ErrUnknown",
	Message: "tried to operate on non gerror.Error",
})

// ExtMsgf attempts to extend an error's message.
func ExtMsgf(err error, format string, args ...any) error {
	gerr, ok := err.(Factory)
	if !ok {
		return ErrUnknown.Convert(err)
	}

	return gerr.Msg(format, args...)
}
