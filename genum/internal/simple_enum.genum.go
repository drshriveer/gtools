// Code generated by genum DO NOT EDIT.
package internal

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/drshriveer/gtools/genum"
	"gopkg.in/yaml.v3"
)

var _MyEnumValues = []MyEnum{
	Enum1IntentionallyNegative,
	Enum1Value0,
	Enum1Value1,
	Enum1Value2,
	Enum1Value7,
}

// IsValid returns true if the enum value is, in fact, valid.
func (e MyEnum) IsValid() bool {
	for _, v := range _MyEnumValues {
		if v == e {
			return true
		}
	}
	return false
}

// Values returns a list of all potential values of this enum.
func (MyEnum) Values() []MyEnum {
	return slices.Clone(_MyEnumValues)
}

// StringValues returns a list of all potential values of this enum as strings.
// Note: This does not return duplicates.
func (MyEnum) StringValues() []string {
	return []string{
		"Enum1IntentionallyNegative",
		"Enum1Value0",
		"Enum1Value1",
		"Enum1Value2",
		"Enum1Value7",
	}
}

// String returns a string representation of this enum.
// Note: in the case of duplicate values only the first alphabetical definition will be choosen.
func (e MyEnum) String() string {
	switch e {
	case Enum1IntentionallyNegative:
		return "Enum1IntentionallyNegative"
	case Enum1Value0:
		return "Enum1Value0"
	case Enum1Value1:
		return "Enum1Value1"
	case Enum1Value2:
		return "Enum1Value2"
	case Enum1Value7:
		return "Enum1Value7"
	default:
		return fmt.Sprintf("UndefinedMyEnum:%d", e)
	}
}

// ParseString will return a value as defined in string form.
func (e MyEnum) ParseString(text string) (MyEnum, error) {
	return ParseMyEnum(text)
}

// ParseMyEnum will attempt to parse the value of a MyEnum from either its string form
// or any value of a trait flagged with the --parsableByTrait flag
func ParseMyEnum(input any) (MyEnum, error) {
	switch input {
	case "Enum1IntentionallyNegative":
		return Enum1IntentionallyNegative, nil
	case "Enum1Value0":
		return Enum1Value0, nil
	case "Enum1Value0Complication1":
		return Enum1Value0Complication1, nil
	case "Enum1Value1":
		return Enum1Value1, nil
	case "Enum1Value1Complication1":
		return Enum1Value1Complication1, nil
	case "Enum1Value2":
		return Enum1Value2, nil
	case "Enum1Value2Complication1":
		return Enum1Value2Complication1, nil
	case "Enum1Value7":
		return Enum1Value7, nil
	default:
		return 0, fmt.Errorf("`%+v` could not be parsed to enum of type MyEnum", input)
	}
}

// ParseGeneric calls TypedEnum.Parse but returns the result
// in the generic genum.Enum interface. Which is useful when you are only able to work with
// the un-typed interface.
func (e MyEnum) ParseGeneric(input any) (genum.Enum, error) {
	return ParseMyEnum(input)
}

// MarshalJSON implements the json.Marshaler interface for MyEnum.
func (e MyEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for MyEnum.
func (e *MyEnum) UnmarshalJSON(data []byte) error {
	// We always support strings.
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		var err error
		*e, err = ParseMyEnum(s)
		if err == nil {
			return nil
		}
	}

	// native parsing

	return fmt.Errorf("unable to unmarshal MyEnum from `%v`", data)
}

// MarshalText implements the encoding.TextMarshaler interface for MyEnum.
func (e MyEnum) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for MyEnum.
func (e *MyEnum) UnmarshalText(text []byte) error {
	s := string(text)
	var err error
	*e, err = ParseMyEnum(s)
	if err == nil {
		return nil
	}

	return fmt.Errorf("unable to unmarshal MyEnum from `%s`", s)
}

// MarshalYAML implements a YAML Marshaler for MyEnum.
func (e MyEnum) MarshalYAML() (any, error) {
	return e.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for MyEnum.
func (e *MyEnum) UnmarshalYAML(value *yaml.Node) error {
	var err error

	// first try and parse as a string
	*e, err = ParseMyEnum(value.Value)
	if err == nil {
		return nil
	}

	// then try and parse for any string-like traits

	// native parsing

	return fmt.Errorf("unable to unmarshal MyEnum from yaml `%s`", value.Value)
}

// IsEnum implements an empty function required to implement Enum.
func (MyEnum) IsEnum() {}

var _MyEnum2Values = []MyEnum2{
	Enum2Value0,
	Enum2Value1,
}

// IsValid returns true if the enum value is, in fact, valid.
func (e MyEnum2) IsValid() bool {
	for _, v := range _MyEnum2Values {
		if v == e {
			return true
		}
	}
	return false
}

// Values returns a list of all potential values of this enum.
func (MyEnum2) Values() []MyEnum2 {
	return slices.Clone(_MyEnum2Values)
}

// StringValues returns a list of all potential values of this enum as strings.
// Note: This does not return duplicates.
func (MyEnum2) StringValues() []string {
	return []string{
		"Enum2Value0",
		"Enum2Value1",
	}
}

// String returns a string representation of this enum.
// Note: in the case of duplicate values only the first alphabetical definition will be choosen.
func (e MyEnum2) String() string {
	switch e {
	case Enum2Value0:
		return "Enum2Value0"
	case Enum2Value1:
		return "Enum2Value1"
	default:
		return fmt.Sprintf("UndefinedMyEnum2:%d", e)
	}
}

// ParseString will return a value as defined in string form.
func (e MyEnum2) ParseString(text string) (MyEnum2, error) {
	return ParseMyEnum2(text)
}

// ParseMyEnum2 will attempt to parse the value of a MyEnum2 from either its string form
// or any value of a trait flagged with the --parsableByTrait flag
func ParseMyEnum2(input any) (MyEnum2, error) {
	switch input {
	case "Enum2Value0":
		return Enum2Value0, nil
	case "Enum2Value1":
		return Enum2Value1, nil
	default:
		return 0, fmt.Errorf("`%+v` could not be parsed to enum of type MyEnum2", input)
	}
}

// ParseGeneric calls TypedEnum.Parse but returns the result
// in the generic genum.Enum interface. Which is useful when you are only able to work with
// the un-typed interface.
func (e MyEnum2) ParseGeneric(input any) (genum.Enum, error) {
	return ParseMyEnum2(input)
}

// MarshalJSON implements the json.Marshaler interface for MyEnum2.
func (e MyEnum2) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for MyEnum2.
func (e *MyEnum2) UnmarshalJSON(data []byte) error {
	// We always support strings.
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		var err error
		*e, err = ParseMyEnum2(s)
		if err == nil {
			return nil
		}
	}

	// native parsing

	return fmt.Errorf("unable to unmarshal MyEnum2 from `%v`", data)
}

// MarshalText implements the encoding.TextMarshaler interface for MyEnum2.
func (e MyEnum2) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for MyEnum2.
func (e *MyEnum2) UnmarshalText(text []byte) error {
	s := string(text)
	var err error
	*e, err = ParseMyEnum2(s)
	if err == nil {
		return nil
	}

	return fmt.Errorf("unable to unmarshal MyEnum2 from `%s`", s)
}

// MarshalYAML implements a YAML Marshaler for MyEnum2.
func (e MyEnum2) MarshalYAML() (any, error) {
	return e.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for MyEnum2.
func (e *MyEnum2) UnmarshalYAML(value *yaml.Node) error {
	var err error

	// first try and parse as a string
	*e, err = ParseMyEnum2(value.Value)
	if err == nil {
		return nil
	}

	// then try and parse for any string-like traits

	// native parsing

	return fmt.Errorf("unable to unmarshal MyEnum2 from yaml `%s`", value.Value)
}

// IsEnum implements an empty function required to implement Enum.
func (MyEnum2) IsEnum() {}

var _MyEnum3Values = []MyEnum3{
	Enum3Value0,
	Enum3Value1,
	Enum3Value2,
	Enum3Value3,
	Enum3Value4,
	Enum3Value5,
	Enum3Value6,
	Enum3Value7,
	Enum3Value8,
	Enum3Value9,
	Enum3Value10,
	Enum3Value11,
	Enum3Value12,
	Enum3Value13,
	Enum3Value15,
	Enum3Value16,
}

// IsValid returns true if the enum value is, in fact, valid.
func (e MyEnum3) IsValid() bool {
	_, ok := slices.BinarySearch(_MyEnum3Values, e)
	return ok
}

// Values returns a list of all potential values of this enum.
func (MyEnum3) Values() []MyEnum3 {
	return slices.Clone(_MyEnum3Values)
}

// StringValues returns a list of all potential values of this enum as strings.
// Note: This does not return duplicates.
func (MyEnum3) StringValues() []string {
	return []string{
		"Enum3Value0",
		"Enum3Value1",
		"Enum3Value2",
		"Enum3Value3",
		"Enum3Value4",
		"Enum3Value5",
		"Enum3Value6",
		"Enum3Value7",
		"Enum3Value8",
		"Enum3Value9",
		"Enum3Value10",
		"Enum3Value11",
		"Enum3Value12",
		"Enum3Value13",
		"Enum3Value15",
		"Enum3Value16",
	}
}

// String returns a string representation of this enum.
// Note: in the case of duplicate values only the first alphabetical definition will be choosen.
func (e MyEnum3) String() string {
	switch e {
	case Enum3Value0:
		return "Enum3Value0"
	case Enum3Value1:
		return "Enum3Value1"
	case Enum3Value2:
		return "Enum3Value2"
	case Enum3Value3:
		return "Enum3Value3"
	case Enum3Value4:
		return "Enum3Value4"
	case Enum3Value5:
		return "Enum3Value5"
	case Enum3Value6:
		return "Enum3Value6"
	case Enum3Value7:
		return "Enum3Value7"
	case Enum3Value8:
		return "Enum3Value8"
	case Enum3Value9:
		return "Enum3Value9"
	case Enum3Value10:
		return "Enum3Value10"
	case Enum3Value11:
		return "Enum3Value11"
	case Enum3Value12:
		return "Enum3Value12"
	case Enum3Value13:
		return "Enum3Value13"
	case Enum3Value15:
		return "Enum3Value15"
	case Enum3Value16:
		return "Enum3Value16"
	default:
		return fmt.Sprintf("UndefinedMyEnum3:%d", e)
	}
}

// ParseString will return a value as defined in string form.
func (e MyEnum3) ParseString(text string) (MyEnum3, error) {
	return ParseMyEnum3(text)
}

// ParseMyEnum3 will attempt to parse the value of a MyEnum3 from either its string form
// or any value of a trait flagged with the --parsableByTrait flag
func ParseMyEnum3(input any) (MyEnum3, error) {
	switch input {
	case "Enum3Value0":
		return Enum3Value0, nil
	case "Enum3Value1":
		return Enum3Value1, nil
	case "Enum3Value2":
		return Enum3Value2, nil
	case "Enum3Value3":
		return Enum3Value3, nil
	case "Enum3Value4":
		return Enum3Value4, nil
	case "Enum3Value5":
		return Enum3Value5, nil
	case "Enum3Value6":
		return Enum3Value6, nil
	case "Enum3Value7":
		return Enum3Value7, nil
	case "Enum3Value8":
		return Enum3Value8, nil
	case "Enum3Value9":
		return Enum3Value9, nil
	case "Enum3Value10":
		return Enum3Value10, nil
	case "Enum3Value11":
		return Enum3Value11, nil
	case "Enum3Value12":
		return Enum3Value12, nil
	case "Enum3Value13":
		return Enum3Value13, nil
	case "Enum3Value15":
		return Enum3Value15, nil
	case "Enum3Value16":
		return Enum3Value16, nil
	default:
		return 0, fmt.Errorf("`%+v` could not be parsed to enum of type MyEnum3", input)
	}
}

// ParseGeneric calls TypedEnum.Parse but returns the result
// in the generic genum.Enum interface. Which is useful when you are only able to work with
// the un-typed interface.
func (e MyEnum3) ParseGeneric(input any) (genum.Enum, error) {
	return ParseMyEnum3(input)
}

// MarshalJSON implements the json.Marshaler interface for MyEnum3.
func (e MyEnum3) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for MyEnum3.
func (e *MyEnum3) UnmarshalJSON(data []byte) error {
	// We always support strings.
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		var err error
		*e, err = ParseMyEnum3(s)
		if err == nil {
			return nil
		}
	}

	// native parsing

	return fmt.Errorf("unable to unmarshal MyEnum3 from `%v`", data)
}

// MarshalText implements the encoding.TextMarshaler interface for MyEnum3.
func (e MyEnum3) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for MyEnum3.
func (e *MyEnum3) UnmarshalText(text []byte) error {
	s := string(text)
	var err error
	*e, err = ParseMyEnum3(s)
	if err == nil {
		return nil
	}

	return fmt.Errorf("unable to unmarshal MyEnum3 from `%s`", s)
}

// MarshalYAML implements a YAML Marshaler for MyEnum3.
func (e MyEnum3) MarshalYAML() (any, error) {
	return e.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for MyEnum3.
func (e *MyEnum3) UnmarshalYAML(value *yaml.Node) error {
	var err error

	// first try and parse as a string
	*e, err = ParseMyEnum3(value.Value)
	if err == nil {
		return nil
	}

	// then try and parse for any string-like traits

	// native parsing

	return fmt.Errorf("unable to unmarshal MyEnum3 from yaml `%s`", value.Value)
}

// IsEnum implements an empty function required to implement Enum.
func (MyEnum3) IsEnum() {}
