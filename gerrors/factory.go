package gerrors

import (
	"fmt"
	"reflect"
	"strings"
)

var gErrT = reflect.TypeOf(GError{})

// The Factory interface exposes only methods that can be used for cloning an error.
// But all errors implement this by default.
// This allows for dynamic and mutable errors without modifying the base.
type Factory interface {
	// Base returns a copy of the embedded error without modifications.
	Base() Error

	// WithStack returns a copy of the embedded error with a Stack trace and diagnostic info.
	WithStack() Error

	// WithSource returns a copy of the embedded error with SourceInfo populated if needed.
	WithSource() Error

	// ExtMsgf returns a copy of the embedded error with diagnostic info and the
	// message extended with additional context.
	ExtMsgf(format string, elems ...any) Error

	// DTagExtMsgf returns a copy of the embedded error with diagnostic info, a detail tag,
	// and the message extended with additional context.
	DTagExtMsgf(detailTag string, format string, elems ...any) Error

	// WithDTag returns a copy of the embedded error with diagnostic info and a detail tag.
	WithDTag(detailTag string) Error

	// Convert will attempt to convert the supplied error into a gError.Error of the
	// Factory's type, including the source errors details in the result's error message.
	// The original error can be retrieved via utility methods.
	Convert(err error) Error

	// Error implements the standard Error interface so that a Factory
	// can be passed into errors.Is() as a target.
	Error() string

	// Is implements the interface for error matching in the standard package (errors.IS).
	Is(error) bool
}

//nolint:errname
type factoryImpl struct {
	ref         Error // for equality
	embedded    *GError
	underlyingT reflect.Type

	embeddedIndex int
	fieldsToClone map[int]reflect.Value
	fieldsToPrint string
}

// FactoryOf takes any kind of error that extends a GError.
func FactoryOf[T Error](err T) Factory {
	if gErr, ok := any(err).(*GError); ok {
		// GErrors _are_ factories.
		return gErr
	}
	underlyingV := reflect.ValueOf(err)
	if underlyingV.Kind() == reflect.Pointer {
		underlyingV = underlyingV.Elem()
	}
	underlyingT := underlyingV.Type()
	embeddedIndex := 0
	fieldsToClone := make(map[int]reflect.Value, 0)
	fieldsToPrint := ""
	for i := 0; i < underlyingT.NumField(); i++ {
		field := underlyingT.Field(i)
		if field.Type.AssignableTo(gErrT) {
			embeddedIndex = i
		}
		tagInfo, ok := field.Tag.Lookup("gerror")
		if !ok {
			continue
		}
		if strings.Contains(tagInfo, "clone") {
			mustBePrimitive(field.Name, field.Type) // or panic.
			fieldsToClone[i] = underlyingV.Field(i)
		}
		if strings.Contains(tagInfo, "print") {
			fieldsToPrint += fmt.Sprintf("%s: %v, ", field.Name, underlyingV.Field(i).Interface())
		}
	}
	embedded := err._embededGError()
	embedded.skipLines = factorySkip
	return &factoryImpl{
		ref:           err,
		embedded:      embedded,
		underlyingT:   underlyingT,
		embeddedIndex: embeddedIndex,
		fieldsToClone: fieldsToClone,
		fieldsToPrint: fieldsToPrint,
	}
}

func mustBePrimitive(fName string, t reflect.Type) {
	// FIXME: need to test with complex types.
	switch t.Kind() {
	case reflect.Invalid,
		reflect.Array,
		reflect.Chan,
		reflect.Func,
		reflect.Map,
		reflect.Slice,
		reflect.Struct:
		panic(fmt.Sprintf("gerror field tags with the 'clone' directive MUST be primitve but field %s is %s", fName, t))
	case reflect.Pointer,
		reflect.Interface,
		reflect.UnsafePointer:
		mustBePrimitive(fName, t.Elem())
	default:
		// good!
		return
	}
}

func (f *factoryImpl) Base() Error {
	gerr := f.embedded.Base()
	return f.cloneUnderlyingWith(gerr)
}

func (f *factoryImpl) WithStack() Error {
	gerr := f.embedded.WithStack()
	return f.cloneUnderlyingWith(gerr)
}

func (f *factoryImpl) WithSource() Error {
	gerr := f.embedded.WithSource()
	return f.cloneUnderlyingWith(gerr)
}

func (f *factoryImpl) ExtMsgf(format string, elems ...any) Error {
	gerr := f.embedded.ExtMsgf(format, elems...)
	return f.cloneUnderlyingWith(gerr)
}

func (f *factoryImpl) DTagExtMsgf(detailTag string, format string, elems ...any) Error {
	gerr := f.embedded.DTagExtMsgf(detailTag, format, elems...)
	return f.cloneUnderlyingWith(gerr)
}

func (f *factoryImpl) WithDTag(detailTag string) Error {
	gerr := f.embedded.WithDTag(detailTag)
	return f.cloneUnderlyingWith(gerr)
}

func (f *factoryImpl) Convert(err error) Error {
	gerr := f.embedded.Convert(err)
	return f.cloneUnderlyingWith(gerr)
}

func (f *factoryImpl) Error() string {
	return "Factory of " + f.embedded.Name
}

func (f *factoryImpl) Is(err error) bool {
	if err == f.ref {
		return true
	}
	return f.ref.Is(err)
}

func (f *factoryImpl) cloneUnderlyingWith(gerr Error) Error {
	gr := gerr.(*GError)
	gr.extensionString = f.fieldsToPrint
	gr.srcFactory = f.embedded

	result := reflect.New(f.underlyingT).Elem()
	result.Field(f.embeddedIndex).Set(reflect.ValueOf(gr).Elem())

	for i, v := range f.fieldsToClone {
		result.Field(i).Set(v)
	}

	return result.Addr().Interface().(Error)
}
