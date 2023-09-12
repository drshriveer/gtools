package internal_test

import (
	"encoding"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/drshriveer/gtools/rutils"

	"github.com/drshriveer/gtools/genum"
	"github.com/drshriveer/gtools/genum/gen"
	"github.com/drshriveer/gtools/genum/internal"
)

func TestSimpleEnumGeneration(t *testing.T) {
	generator := gen.Generate{
		InFile:        "./simple_enum.go",
		OutFile:       "./simple_enum.genum.go",
		EnumTypeNames: []string{"MyEnum", "MyEnum2", "MyEnum3"},
		GenJSON:       true,
		GenYAML:       true,
		GenText:       true,
	}

	require.NoError(t, generator.Parse())
}

func TestImplementsEnumInterface(t *testing.T) {
	assert.Implements(t, (*genum.Enum)(nil), internal.Enum1Value0)
	assert.Implements(t, (*genum.Enum)(nil), internal.Enum2Value0)
	assert.Implements(t, (*genum.Enum)(nil), internal.Enum3Value0)
}

type enumTest[T genum.EnumLike] struct {
	enum                genum.TypedEnum[T]
	sName               string
	invalid             bool
	duplicateDefinition bool
}

func (e *enumTest[T]) roundTripJson(t *testing.T) {
	require.Implements(t, (*json.Unmarshaler)(nil), reflect.New(reflect.TypeOf(e.enum)).Interface())
	require.Implements(t, (*json.Marshaler)(nil), e.enum)
	bytes, err := json.Marshal(e.enum)
	require.NoError(t, err)
	assert.Equal(t, `"`+e.sName+`"`, string(bytes))
	ptrVal := reflect.New(reflect.TypeOf(e.enum)).Interface()
	require.NoError(t, json.Unmarshal(bytes, ptrVal))
	assert.Equal(t, e.enum, rutils.Unptr(ptrVal))
}

func (e *enumTest[T]) roundTripYaml(t *testing.T) {
	require.Implements(t, (*yaml.Unmarshaler)(nil), reflect.New(reflect.TypeOf(e.enum)).Interface())
	require.Implements(t, (*yaml.Marshaler)(nil), e.enum)
	bytes, err := yaml.Marshal(e.enum)
	require.NoError(t, err)
	assert.Equal(t, e.sName+"\n", string(bytes))
	ptrVal := reflect.New(reflect.TypeOf(e.enum)).Interface()
	require.NoError(t, yaml.Unmarshal(bytes, ptrVal))
	assert.Equal(t, e.enum, rutils.Unptr(ptrVal))
}

func (e *enumTest[T]) roundTripText(t *testing.T) {
	require.Implements(t, (*encoding.TextUnmarshaler)(nil), reflect.New(reflect.TypeOf(e.enum)).Interface())
	require.Implements(t, (*encoding.TextMarshaler)(nil), e.enum)
	bytes, err := e.enum.(encoding.TextMarshaler).MarshalText()
	require.NoError(t, err)
	assert.Equal(t, e.sName, string(bytes))
	ptrVal := reflect.New(reflect.TypeOf(e.enum)).Interface().(encoding.TextUnmarshaler)
	require.NoError(t, ptrVal.UnmarshalText(bytes))
	assert.Equal(t, e.enum, rutils.Unptr(ptrVal))
}

func testRunner[T genum.EnumLike](t *testing.T, tests []enumTest[T]) {
	t.Run("general type test", func(t *testing.T) {
		validValues := make([]genum.Enum, 0, len(tests))
		validValueStrings := make([]string, 0, len(tests))
		for _, test := range tests {
			if !test.invalid {
				if !test.duplicateDefinition {
					validValues = append(validValues, test.enum)
					validValueStrings = append(validValueStrings, test.sName)
				}
			}
		}

		v := rutils.Unptr(reflect.New(reflect.TypeOf(*new(T))).Interface()).(genum.TypedEnum[T])
		actualStringValues := v.StringValues()
		assert.Len(t, actualStringValues, len(validValueStrings))
		assert.ElementsMatch(t, actualStringValues, validValueStrings)
		actualValues := v.Values()
		assert.Len(t, actualValues, len(validValues))
		assert.ElementsMatch(t, actualValues, validValues)

	})

	for _, test := range tests {
		t.Run(test.enum.String(), func(t *testing.T) {
			if test.invalid {
				assert.False(t, test.enum.IsValid())
				return
			}

			assert.True(t, test.enum.IsValid())
			assert.Equal(t, test.sName, test.enum.String())
			v := rutils.Unptr(reflect.New(reflect.TypeOf(*new(T))).Interface()).(genum.TypedEnum[T])
			r, err := v.ParseString(test.sName)
			require.NoError(t, err)
			assert.Equal(t, test.enum, r)

			// marshals:
			test.roundTripJson(t)
			test.roundTripYaml(t)
			test.roundTripText(t)
		})
	}
}

func TestMyEnum3(t *testing.T) {
	tests := []enumTest[internal.MyEnum3]{
		{enum: internal.Enum3Value0, sName: "Enum3Value0"},
		{enum: internal.Enum3Value1, sName: "Enum3Value1"},
		{enum: internal.Enum3Value2, sName: "Enum3Value2"},
		{enum: internal.Enum3Value3, sName: "Enum3Value3"},
		{enum: internal.Enum3Value4, sName: "Enum3Value4"},
		{enum: internal.Enum3Value5, sName: "Enum3Value5"},
		{enum: internal.Enum3Value6, sName: "Enum3Value6"},
		{enum: internal.Enum3Value7, sName: "Enum3Value7"},
		{enum: internal.Enum3Value8, sName: "Enum3Value8"},
		{enum: internal.Enum3Value9, sName: "Enum3Value9"},
		{enum: internal.Enum3Value10, sName: "Enum3Value10"},
		{enum: internal.Enum3Value11, sName: "Enum3Value11"},
		{enum: internal.Enum3Value12, sName: "Enum3Value12"},
		{enum: internal.Enum3Value13, sName: "Enum3Value13"},
		{enum: internal.MyEnum3(14), invalid: true},
		{enum: internal.Enum3Value15, sName: "Enum3Value15"},
		{enum: internal.Enum3Value16, sName: "Enum3Value16"},
		{enum: internal.MyEnum3(99), invalid: true},
	}

	testRunner[internal.MyEnum3](t, tests)
}

func TestMyEnum2(t *testing.T) {
	tests := []enumTest[internal.MyEnum2]{
		{
			enum:  internal.Enum2Value0,
			sName: "Enum2Value0",
		},
		{
			enum:  internal.Enum2Value1,
			sName: "Enum2Value1",
		},
		{
			enum:    internal.MyEnum2(2),
			invalid: true,
		},
	}

	testRunner[internal.MyEnum2](t, tests)
}

func TestMyEnum(t *testing.T) {
	tests := []enumTest[internal.MyEnum]{
		{
			enum:  internal.Enum1Value0,
			sName: "Enum1Value0",
		},
		{
			enum:                internal.Enum1Value0Complication1,
			sName:               "Enum1Value0",
			duplicateDefinition: true,
		},
		{
			enum:  internal.Enum1Value1,
			sName: "Enum1Value1",
		},
		{
			enum:                internal.Enum1Value1Complication1,
			sName:               "Enum1Value1",
			duplicateDefinition: true,
		},
		{
			enum:  internal.Enum1Value2,
			sName: "Enum1Value2",
		},
		{
			enum:                internal.Enum1Value2Complication1,
			sName:               "Enum1Value2",
			duplicateDefinition: true,
		},
		{
			enum:  internal.Enum1Value7,
			sName: "Enum1Value7",
		},
		{
			enum:  internal.Enum1IntentionallyNegative,
			sName: "Enum1IntentionallyNegative",
		},
		{
			enum:    internal.MyEnum(99),
			invalid: true,
		},
		{
			enum:    internal.MyEnum(6),
			invalid: true,
		},
	}

	testRunner[internal.MyEnum](t, tests)
}
