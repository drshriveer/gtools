package internal

import (
	"github.com/drshriveer/gtools/gerrors"
)

// GRPCError is just a test.
//
//go:generate gerror types=GRPCError
type GRPCError struct {
	gerrors.GError
	GRPCStatus      Status `gerror:"print,factory"`
	CustomerMessage string `gerror:"print"`

	// Do not print, or create a factory for
	DoNotPrint string
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
