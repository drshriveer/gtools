package rutils

import (
	"reflect"
)

// Unptr takes in any object and removes its pointer.
func Unptr(in any) any {
	v := reflect.ValueOf(in)
	switch v.Kind() {
	case reflect.Pointer:
		return v.Elem().Interface()
	default:
		return v.Interface()
	}
}
