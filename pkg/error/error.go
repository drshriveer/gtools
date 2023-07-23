package error

// G세늠 G세넘

import (
	"errors"
)

// FIXME: GAVIN!!
//  - what other kinds of diagnostic info can I get?
//        could I return the fileName? MethodName? PackageName?
//  - what if a factory added metrics for all the times it was created?
//  - what if a factory optimized itself to cache frequent stack traces?
//  - Consider: for the use of switch statements, maybe there does need to be
//    a concept of wrapper-- but that really can be the "error" where the rest
//    is the factory and we are essentially only comparing the factories
//  - what if there was a _service_ stack (appended to stack upon passing instrumented layers)
//    - cool? not really sure at this point what could be done with that data
//    - maybe something clever with error graphs?
//    - can grafana make a waterfall from
//  - what if errors could indicate whether or not they're alarm worthy?
//    - could be done by setting up a special namespace?
//    - alarms.UnacceptableErrors.gshriver.<pkg_name|src_source>.method.name-detail
//    - alarms.UnacceptableErrors.team.<pkg_name|src_source>.method.name-detail
//    - even could include tags of people to identify... team, individual,
//
//  - what could we change if context was an input to factory methods)

// The Factory interface exposes only methods that can be used for cloning an error.
// But all errors implement this by default.
// This allows for dynamic and mutable errors without modifying the base.
type Factory interface {

	// E returns a copy of the underlying error with a stack trace and diagnostic info.
	E() GError

	// Raw returns a copy of the underlying error without a stack trace or diagnostic info.
	Raw() GError

	// Merge returns a copy of the underlying error with a stack trace and diagnostic info.
	// merges set properties into the result.z
	Merge(gError GError) GError

	// FIXME: lint rules for formats! (govet)
	// Include adds an additional error message.
	Include(format string, elems any) GError

	// DInclude (DetailInclude) returns a copy of the underlying error with a stack trace,
	// diagnostic info, an updated metric tag and a message
	// FIXME: lint rules for formats! (govet)
	DInclude(metricTag string, format string, elems any) GError

	// MetricTag returns a copy of the underlying error with a stack trace, diagnostic info,
	// and a metric tag.
	MetricTag(tag string) GError

	// Will try to convert an error into
	Convert(err error) GError
}

type MonitoredError struct {
	Count       int
	AlertSpaces []string
}

type ErrorInterface interface {
	// humm... how do I want to re-implement this to make it better
	// than before?
	// 1. equality checks are important
	// 2. stack traces and wrapping / dependency on pkgErrors?
	// 		- ENSURE Stack (but don't add more!)
	// 3. factory pattern?
	// 4. global registration? (for permanent ids ?)
	// 5. EXTENSABILITY!

	// Maybe the Error interface is
	error
	// grpc.StatusError
	ExtendTag()
	Code()
	MetricString() string
	Message() string

	// Standard error package interfaces:
	Error() string
	Is(error) bool
	Unwrap() error
}

type GError struct {
	// The things that have to be equal:
	Source  string
	Name    string
	Message string
	// Code  codes.Code

	// The things that don't have to be equal:
	ExtMessage string
	MetricTag  string

	stack []string
	// comi
	srcFactory *GError
}

func (e GError) Equal(err GError) bool {
	// this is a possiblity
	// return e.srcFactory == err.srcFactory
	// or whatever critera
	return e.Source == err.Source
}

func T() {
	errors.Join()
	// pkgErrors.Wrap()
	// v := Wrapped{}
}

func Equal(err1, err2 error) {

}
