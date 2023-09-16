package gerrors

// Unwrap finds the base error.
func Unwrap(err error) error {
	switch v := err.(type) {
	case GError:
		if v.srcFactory != nil {
			return v.srcFactory
		}
	case *GError:
		if v.srcFactory != nil {
			return v.srcFactory
		}
	}
	// dunno if any further unwrapping is required...
	return err
}
