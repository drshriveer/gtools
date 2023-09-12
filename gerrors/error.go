package gerrors

import (
	"fmt"
)

type stackType int

const (
	noStack      stackType = 0
	sourceOnly   stackType = 1
	defaultStack stackType = 32
)

// G세늠 G세넘

// FIXME: GAVIN!!
//  - what other kinds of diagnostic info can I get?
//        could I return the fileName? MethodName? PackageName?
//  - what if a factory added metrics for all the times it was created?
//  - what if a factory optimized itself to cache frequent Stack traces?
//  - Consider: for the use of switch statements, maybe there does need to be
//    a concept of wrapper-- but that really can be the "error" where the rest
//    is the factory and we are essentially only comparing the factories
//  - what if there was a _service_ Stack (appended to Stack upon passing instrumented layers)
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
	error
	// WithStack returns a copy of the underlying error with a Stack trace and diagnostic info.
	WithStack() GError

	// WithSource returns a copy of the underlying error with diagnostic info.
	WithSource() GError

	// Raw returns a copy of the underlying error without a Stack trace or diagnostic info.
	Raw() GError

	// Merge returns a copy of the underlying error with a Stack trace and diagnostic info.
	// merges set properties into the result.
	Merge(gError GError) GError

	// FIXME: lint rules for formats! (govet)
	// Include adds an additional error message.
	Include(format string, elems ...any) GError

	// DInclude (DetailInclude) returns a copy of the underlying error with a Stack trace,
	// diagnostic info, an updated metric tag and a message
	// FIXME: lint rules for formats! (govet)
	DInclude(metricTag string, format string, elems ...any) GError

	// MetricTag returns a copy of the underlying error with a Stack trace, diagnostic info,
	// and a metric tag.
	MetricTag(mTag string) GError

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
	// 2. Stack traces and wrapping / dependency on pkgErrors?
	// 		- ENSURE Stack (but don't add more!)
	// 3. factory pattern?
	// 4. global registration? (for permanent ids ?)
	// 5. EXTENSABILITY!

	// Maybe the Error interface is
	error
	// grpc.StatusError
	ExtendTag(string)
	Code()
	MetricString() string
	Message() string

	// return the factory for comparison
	SrcFactory() Factory

	// Standard error package interfaces:
	Error() string
	Is(error) bool
	Unwrap() error // XXX: what would we unwrap to? a separate unknown source? a factory?
}

type GError struct {
	// The things that have to be equal:
	Name    string
	Message string

	// The things that don't have to be equal:
	Source     string
	ExtMessage string
	MTag       string

	// trace info:
	Stack *Stack

	// consider using his as the true equality check for slices.
	// e.g. switch gerrors.Unwrap(err) ...
	//
	SrcFactory *GError
}

// implements the "error" interface.
func (e GError) Error() string {
	// FIXME:
	//   - only include elements if they're present...
	//   - source should not be so dumb?
	return fmt.Sprintf("Name: %s, MTag: %s, Src: %s, Message: %s %s \n%s",
		e.Name, e.MTag, e.Source, e.Message, e.ExtMessage, e.Stack)
}

func (e GError) Is(err error) bool {
	gerr, ok := err.(GError)
	if !ok {
		return false
	}
	// this is a possiblity
	// return e.SrcFactory == err.SrcFactory
	// or whatever critera
	return e.Message == gerr.Message && e.Name == gerr.Name
}

func (e *GError) WithStack() GError {
	return e.clone(defaultStack)
}

func (e *GError) WithSource() GError {
	return e.clone(sourceOnly)
}

func (e *GError) Raw() GError {
	return e.clone(noStack)
}

func (e *GError) Merge(gError GError) GError {
	clone := e.clone(noStack)
	if len(gError.Name) > 0 {
		clone.Name = gError.Name
	}
	if len(gError.Message) > 0 {
		clone.Message = gError.Message
	}
	if len(gError.Source) > 0 {
		clone.Source = gError.Source
	}
	if len(gError.ExtMessage) > 0 {
		if e.ExtMessage != "" {
			clone.ExtMessage = fmt.Sprintf("%s %s", e.ExtMessage, gError.ExtMessage)
		} else {
			clone.ExtMessage = gError.ExtMessage
		}
	}
	if len(gError.MTag) > 0 {
		e.MTag = gError.MTag
	}
	return clone
}

func (e *GError) Include(format string, elems ...any) GError {
	clone := e.clone(defaultStack)
	clone.ExtMessage = fmt.Sprintf(format, elems...)
	return clone
}

func (e *GError) DInclude(mTag string, format string, elems ...any) GError {
	clone := e.clone(defaultStack)
	clone.ExtMessage = fmt.Sprintf(format, elems...)
	clone.MTag = mTag
	return clone
}

func (e *GError) MetricTag(mTag string) GError {
	clone := e.clone(defaultStack)
	clone.MTag = mTag
	return clone
}

func (e *GError) clone(st stackType) GError {
	clone := *e
	clone.SrcFactory = e
	clone.ExtMessage = ""
	clone.MTag = ""

	if st == noStack || (st == sourceOnly && clone.Source != "") {
		clone.Stack = nil
	} else {
		clone.Stack = makeStack(int(st), 4)
		clone.Source = (*clone.Stack)[0].Metric()
		if st == sourceOnly {
			clone.Stack = nil
		}
	}
	return clone
}

func (e *GError) Convert(err error) GError {
	switch v := err.(type) {
	case GError:
		return v
	case *GError:
		return *v
	}

	clone := e.clone(defaultStack)
	clone.ExtMessage = fmt.Sprintf("originalError: %+v", err)

	return clone
}

func (e *GError) Unwrap() error {
	if e.SrcFactory != nil {
		return e.SrcFactory
	}
	return e
}

func Unwrap(err error) error {
	switch v := err.(type) {
	case GError:
		if v.SrcFactory != nil {
			return v.SrcFactory
		}
	case *GError:
		if v.SrcFactory != nil {
			return v.SrcFactory
		}
	}
	// dunno if any further unwrapping is required...
	return err
}
