// Code generated by genum DO NOT EDIT.
package internal

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	stupidTime "time"

	"gopkg.in/yaml.v3"
)

var _EnumerableWithTraitsValues = []EnumerableWithTraits{
	E1,
	E2,
	E3,
}

// Timeout returns the enum's associated trait of the same name.
// If no trait exists for the enumeration a default value will be returned.
func (e EnumerableWithTraits) Timeout() stupidTime.Duration {
	switch e {
	case E1:
		return _Timeout
	case E2:
		return 1 * stupidTime.Minute
	case E3:
		return 2 * stupidTime.Minute
	}

	return *new(stupidTime.Duration)
}

// Trait returns the enum's associated trait of the same name.
// If no trait exists for the enumeration a default value will be returned.
func (e EnumerableWithTraits) Trait() string {
	switch e {
	case E1:
		return _Trait
	case E2:
		return "trait 2"
	case E3:
		return "trait 3"
	}

	return *new(string)
}

// TypedStringTrait returns the enum's associated trait of the same name.
// If no trait exists for the enumeration a default value will be returned.
func (e EnumerableWithTraits) TypedStringTrait() OtherType {
	switch e {
	case E1:
		return _TypedStringTrait
	case E2:
		return OtherType("OtherType2")
	case E3:
		return OtherType("OtherType3")
	}

	return *new(OtherType)
}

// IsValid returns true if the enum value is, in fact, valid.
func (e EnumerableWithTraits) IsValid() bool {
	for _, v := range _EnumerableWithTraitsValues {
		if v == e {
			return true
		}
	}
	return false
}

// Values returns a list of all potential values of this enum.
func (EnumerableWithTraits) Values() []EnumerableWithTraits {
	return slices.Clone(_EnumerableWithTraitsValues)
}

// StringValues returns a list of all potential values of this enum as strings.
// Note: This does not return duplicates.
func (EnumerableWithTraits) StringValues() []string {
	return []string{
		"E1",
		"E2",
		"E3",
	}
}

// String returns a string representation of this enum.
// Note: in the case of duplicate values only the first alphabetical definition will be choosen.
func (e EnumerableWithTraits) String() string {
	switch e {
	case E1:
		return "E1"
	case E2:
		return "E2"
	case E3:
		return "E3"
	default:
		return fmt.Sprintf("UndefinedEnumerableWithTraits:%d", e)
	}
}

// ParseString will return a value as defined in string form.
func (e EnumerableWithTraits) ParseString(text string) (EnumerableWithTraits, error) {
	switch text {
	case "E1":
		return E1, nil
	case "E2":
		return E2, nil
	case "E3":
		return E3, nil
	default:
		return 0, fmt.Errorf("`%s` is not a valid enum of type EnumerableWithTraits", text)
	}
}

// MarshalJSON implements the json.Marshaler interface for EnumerableWithTraits.
func (e EnumerableWithTraits) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for EnumerableWithTraits.
func (e *EnumerableWithTraits) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		var err error
		*e, err = EnumerableWithTraits(0).ParseString(s)
		return err
	}
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*e = EnumerableWithTraits(i)
		if e.IsValid() {
			return nil
		}
	}

	return fmt.Errorf("unable to unmarshal EnumerableWithTraits from `%v`", data)
}

// MarshalText implements the encoding.TextMarshaler interface for EnumerableWithTraits.
func (e EnumerableWithTraits) MarshalText() ([]byte, error) {
	return []byte(e.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for EnumerableWithTraits.
func (e *EnumerableWithTraits) UnmarshalText(text []byte) error {
	var err error
	*e, err = EnumerableWithTraits(0).ParseString(string(text))
	return err
}

// MarshalYAML implements a YAML Marshaler for EnumerableWithTraits.
func (e EnumerableWithTraits) MarshalYAML() (any, error) {
	return e.String(), nil
}

// UnmarshalYAML implements a YAML Unmarshaler for EnumerableWithTraits.
func (e *EnumerableWithTraits) UnmarshalYAML(value *yaml.Node) error {
	i, err := strconv.ParseInt(value.Value, 10, 64)
	if err == nil {
		*e = EnumerableWithTraits(i)
	} else {
		*e, err = EnumerableWithTraits(0).ParseString(value.Value)
	}
	if err != nil {
		return err
	} else if e.IsValid() {
		return nil
	}
	return fmt.Errorf("unable to unmarshal EnumerableWithTraits from yaml `%s`", value.Value)
}

// IsEnum implements an empty function required to implement Enum.
func (EnumerableWithTraits) IsEnum() {}