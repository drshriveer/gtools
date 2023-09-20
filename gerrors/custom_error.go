package gerrors

//go:generate gerror --types=GRPCError

// GRPCError is just a test.
type GRPCError struct {
	GError
	GRPCStatus      Status `gerror:"_,print,factory"`
	CustomerMessage string `gerror:"_,print"`

	// Do not print, or create a factory for
	DoNotPrint string
}

//
// var ErrExtendedExample gerrors.Factory = &GRPCError{
// 	GError: gerrors.GError{
// 		Name:    "ErrExtendedExample",
// 		Message: "extended error example",
// 	},
// 	GRPCStatus:      InvalidArgument,
// 	CustomerMessage: "Print this message",
// 	DoNotPrint:      "this is for internal issue only",
// }
//
// func L1() error {
// 	return L2()
// }
// func L2() error {
// 	return L3()
// }
// func L3() error {
// 	return ErrExtendedExample.Stack()
// }

// Status is just a sample for testing.
type Status int

// This block exports for testing only.
const (
	OK Status = iota
	Canceled
	Unknown
	InvalidArgument
	DeadlineExceeded
)

func (s Status) String() string {
	switch s {
	case OK:
		return "OK"
	case Canceled:
		return "Canceled"
	case Unknown:
		return "Unknown"
	case InvalidArgument:
		return "InvalidArgument"
	case DeadlineExceeded:
		return "DeadlineExceeded"
	}
	return "UNKNOWN"
}
