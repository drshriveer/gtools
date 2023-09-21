# gError

gError is an opinionated error model for Go that takes a metric-first standpoint on errors.

[Docs](https://pkg.go.dev/github.com/drshriveer/gtool/gerror)

### Getting started

Install with:

```bash
go install github.com/drshriveer/gtool/gerror@latest
```

### Features

- **Named Errors** - Errors have names that can be emitted in metrics. 
- **Source Identity** - Errors can have a static source or can dynamically determine their source when returned.
  - TODO: certain types of binary builds may limit the introspection capabilities; document this here.
- **Stack traces** - gError support stack traces if desired no need to depend on something like [pkg/errors](https://pkg.go.dev/github.com/pkg/errors). Errors ensure they have stacks but do duplicate stacks.
- **Detail Tags** - Errors support metric-safe detail tags  
- **ErrorFactory** - Factories aid in all of the above. 

### Tenants 

Below are the tenants that lead to this library's development.
Keeping them in mind will better aid in understanding how to use this gError effectively.

- **All errors in an application should be consistently gError**
  - This goes along with "Errors should treat metrics as a first-class citizen".
  - _Convert_ errors from external libraries into gError using a factory's `Convert` method.
  - Implement client interceptors to automatically convert errors into the correct type.
    - TODO: provide example.
- **Returned errors should be tested**
  - Raw string matching and error wrapping make testing errors brittle. Verify the error returned is the error you intended.
- **Errors must treat metrics as a first-class citizen**
  - That means errors need to be _Named_ and have a concept of their _Source_.
  - Consider pairing gError with something like [gowrap](https://github.com/hexdigest/gowrap) to generate instrumented interfaces that emit metrics when an error is encountered.
    - TODO: provide example.
- **Errors should be handleable in switch statements**
  - Specific errors may require special handling. Inspecting on individual attributes of an error (status code, error string, error contains string, ec), leads to brittle and even dangerous code, so switching should be made as easy as possible. Thus support switch statements! 
- **Errors should be extensible**
  - Errors sometimes need extra information (e.g. GRPC status codes, HTTP status codes, customer-facing error messages vs internal error messages, etc) that is not included in a base error. For that reason gError are extensible in a case-by-case basis. 
- **Errors should be reusable**
  - It should not be necessary to re-define an error for every use case e.g. ErrInvalidParameter should be valid whether a field is malformed, or a required parameter is missing. However, it should be possible to _distinguish_ between the reason an error was returned from the same path. DetailsTags help with this.
- **Errors should be predefined**
  - Many of the tenants above converge on this point: never return an error created on-the-fly, define them so that they can be tested, reused, and handled as the parent type.
- **Limit error wrapping**
  - Wrapped errors have several drawbacks when it comes to the development experience... Different methods of wrapping (cause vs unwrap), searching the linked list for a specific kind of wrap, etc. More frequently than not, I have seen this lead to brittle code, bugs, and confusion. While unavoidable to a small degree, this library does its best to avoid it.  

### Usage

##### General

**Define an error:**
```go
var InvalidArgument = gerror.FactoryOf(&gerror.GError{
    Name:    "InvalidArgument",
    Message: "this is error 1",
})
```

**Return it with a factory method:**
```go
func FuncName(input InType) error {
	if InType.Field1.Invalid() {
		return InvalidArgument.DTag("Field1")
    }
  return nil
}
```
##### Extend

###### Example: GRPCError

NOTE: example is incomplete ATM. 

```go
// define the type and a generator:
//go:generate gerror --types=GRPCError
type GRPCError struct {
	gerror.GError // embed 
	GRPCStatus      codes.Code    `gerror:"_,print,clone"`
}

func (e *GRPCError) Code() codes.Code {
	return e.codes
} 

func (e *GRPCError) Staus() grpcProtos.StatusError {
	// TODO: write this func correctly
}

```

##### Client Interceptor Example

// TODO

#### GoWrap Example

// TODO

### Limitations:

- Still need `errors.Unwrap(err)` before equality check (without `errors.Is`) or switch statement.  
- Internal code is bonkers.
- ErrSource is derived off random(ish) rules. Need to better understand internals to improve. 

### TODO:
- Consider factory Config:
  - global or otherwise 
  - Stack sampling  (golang.org/x/time/rate::Sometimes)
    - global, factory, or type (via annotations) specific
  - ALARM ON / Severity
- converge on metric-safe "source" string (or a way to configure this)
- possible to split library into specific versions for grpc / http / etc modules?
- linter:
  - for metric-safe detail tags
  - error name must match variable name
- Revisit later:
  - ExtMessage as a first class citizen or not.
  - How an error string is presented
    - Ordering of wrapped details
  - How to combine wrapped message extensions (with a ` `... or?)
  - How to combine DetailTags (with a `-` or?)
  - Consider when and weather to deep clone.
