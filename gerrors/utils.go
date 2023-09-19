package gerrors

// ExtMsgf attempts to extend an error's message.
func ExtMsgf(err error, format string, args ...any) error {
	if gerr, ok := err.(Error); ok {
		return gerr.ExtMsgf(format, args...)
	}
	// TODO: consider trying to call a Convert to a base "internal" error
	return err // failed
}
