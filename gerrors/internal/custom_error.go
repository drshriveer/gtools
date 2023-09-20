package internal

import (
	"time"

	"github.com/drshriveer/gtools/gerrors"
)

//go:generate gerror --types=GRPCError

// GRPCError is just a test.
type GRPCError struct {
	gerrors.GError
	GRPCStatus      Status        `gerror:"_,print,clone"`
	CustomerMessage string        `gerror:"_,print"`
	Timeout         time.Duration `gerror:"_,clone"`

	// Do not print, or create a factory for
	DoNotPrint string
}

// ErrExtendedExample is an example error.
var ErrExtendedExample = gerrors.FactoryOf(&GRPCError{
	GError: gerrors.GError{
		Name:    "ErrExtendedExample",
		Message: "extended error example",
	},
	GRPCStatus:      InvalidArgument,
	CustomerMessage: "Print this message",
	DoNotPrint:      "this is for internal issue only",
})

// L1 is layer one for testing stack traces.
func L1() error {
	return L2()
}

// L2 is layer one for testing stack traces.
func L2() error {
	return L3()
}

// L3 is layer one for testing stack traces.
func L3() error {
	return ErrExtendedExample.Stack()
}

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
