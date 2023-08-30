package internal_test

import (
	"encoding"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/drshriveer/gcommon/pkg/enum"
	"github.com/drshriveer/gcommon/pkg/enum/gen"
	"github.com/drshriveer/gcommon/pkg/enum/internal"
)

func TestEnumerableWithTraits(t *testing.T) {
	generator := gen.Generate{
		InFile:        "./enumerable_with_traits.go",
		OutFile:       "./enumerable_with_traits.genum.go",
		EnumTypeNames: []string{"EnumerableWithTraits"},
		GenJSON:       true,
		GenYAML:       true,
		GenText:       true,
	}

	require.NoError(t, generator.Parse())
	require.NoError(t, generator.Write())

	t.Run("types", func(t *testing.T) {
		var value enum.Enum = internal.E1

		eType := reflect.TypeOf(value)
		eKind := eType.Kind()

		value2 := reflect.New(eType).Interface()

		println(fmt.Sprintf("etype: %s, eKind: %s, value2: %T:%s ", eType, eKind, value2, value2))

		umashaller, ok := value2.(encoding.TextUnmarshaler)
		require.True(t, ok)
		require.NoError(t, umashaller.UnmarshalText([]byte("E2")))
		assert.Equal(t, internal.E2, unptr(value2))

	})
}

func genericTest[T any]() T {
	return *new(T)
}

func unptr(in any) any {
	v := reflect.ValueOf(in)
	switch v.Kind() {
	case reflect.Pointer:
		return v.Elem().Interface()
	default:
		return v.Interface()
	}
}
