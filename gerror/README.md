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

###### General

```go
// Define an error: 
var ErrInvalidArgument = GErrorFactory{
    Name:    "ErrMyError1",
    Message: "this is error 1",
}

```
###### Extend

**GRPC**

**ClientError**

##### Client Interceptor Example

// TODO

### GoWrap Example

// TODO


- Define the enum (`Creatures` above) in a file (e.g. `filename.go`)
- Add the generate directive: `//go:generate genum -types=Creatures`.
- Run the `go generate` command.
- Code will be generated and written to file `<filanme>.genum.go` in the same package.

### Limitations: 
- Does not print extended gerror fields that have been mutated. i.e. mutations after an error is created via factory are limited to message and detail tag.
- cloned fields _must be immutable_.

### TODO:
- Consider factory Config:
  - global or otherwise 
  - Stack sampling  (golang.org/x/time/rate::Sometimes)
  - ALARM ON / Severity
  - stack configurations: always, never, source only
  - 
- converge on metric-safe "source" string (or a way to configure this)
- consider metric-aware error factories (return count, ec)
- possible to split library into specific versions for grpc / http / etc modules?
- linter:
  - for metric-safe detail tags
  - error name must match variable name
- TEST EXTENSIONS
- Later revisit / converge on:
  - ExtMessage as a first class citizen or not.
  - How an error string is presented
    - Ordering of wrapped details
    - 
  - How to combine wrapped message extensions
  - How to combine DetailTags
- Consider when and weather to deep clone.
- Lint use of err1 == err2 vs errors.Is