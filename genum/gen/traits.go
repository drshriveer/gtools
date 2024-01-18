package gen

import (
	"go/types"

	"github.com/drshriveer/gtools/gencommon"
)

type underlying int

const (
	Unknown underlying = iota
	String
	Uint64
	Int64
	Float64
	Float32
)

// TraitDescs is a sortable slice of TraitDesc.
type TraitDescs []TraitDesc

func (s TraitDescs) Len() int {
	return len(s)
}

func (s TraitDescs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s TraitDescs) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

func (s TraitDescs) getParsableUnderlying(u underlying, excluding func(*TraitDesc) bool) TraitDescs {
	out := make([]TraitDesc, 0, len(s))
	for _, t := range s {
		if t.Parsable && t.hasUnderlying(u) && !excluding(&t) {
			out = append(out, t)
		}
	}
	return out
}

func (s TraitDescs) GetParsableUnderlyingStringForJSON() TraitDescs {
	return s.getParsableUnderlying(String, implementsJSONUnmarshaler)
}

func (s TraitDescs) GetParsableUnderlyingFloat64ForJSON() TraitDescs {
	return s.getParsableUnderlying(Float64, implementsJSONUnmarshaler)
}

func (s TraitDescs) GetParsableUnderlyingFloat32ForJSON() TraitDescs {
	return s.getParsableUnderlying(Float32, implementsJSONUnmarshaler)
}

func (s TraitDescs) GetParsableUnderlyingInt64ForJSON() TraitDescs {
	return s.getParsableUnderlying(Int64, implementsJSONUnmarshaler)
}

func (s TraitDescs) GetParsableUnderlyingUint64ForJSON() TraitDescs {
	return s.getParsableUnderlying(Uint64, implementsJSONUnmarshaler)
}

func (s TraitDescs) GetParsableUnderlyingStringForYAML() TraitDescs {
	return s.getParsableUnderlying(String, implementsYAMLUnmarshaler)
}

func (s TraitDescs) GetParsableUnderlyingFloat64ForYAML() TraitDescs {
	return s.getParsableUnderlying(Float64, implementsYAMLUnmarshaler)
}

func (s TraitDescs) GetParsableUnderlyingFloat32ForYAML() TraitDescs {
	return s.getParsableUnderlying(Float32, implementsYAMLUnmarshaler)
}

func (s TraitDescs) GetParsableUnderlyingInt64ForYAML() TraitDescs {
	return s.getParsableUnderlying(Int64, implementsYAMLUnmarshaler)
}

func (s TraitDescs) GetParsableUnderlyingUint64ForYAML() TraitDescs {
	return s.getParsableUnderlying(Uint64, implementsYAMLUnmarshaler)
}

func (s TraitDescs) GetParsableUnderlyingStringForText() TraitDescs {
	return s.getParsableUnderlying(String, implementsTextUnmarshaler)
}

func (s TraitDescs) GetParsableJSONUnmarshalable() TraitDescs {
	out := make([]TraitDesc, 0, len(s))
	for _, t := range s {
		if t.Parsable && implementsJSONUnmarshaler(&t) {
			out = append(out, t)
		}
	}
	return out
}

func (s TraitDescs) GetParsableYAMLUnmarshalable() TraitDescs {
	out := make([]TraitDesc, 0, len(s))
	for _, t := range s {
		if t.Parsable && implementsYAMLUnmarshaler(&t) {
			out = append(out, t)
		}
	}
	return out
}

func (s TraitDescs) GetParsableTextUnmarshalable() TraitDescs {
	out := make([]TraitDesc, 0, len(s))
	for _, t := range s {
		if t.Parsable && implementsTextUnmarshaler(&t) {
			out = append(out, t)
		}
	}
	return out
}

// TraitDesc define a trait-- this is exposed for template use.
type TraitDesc struct {
	Name     string
	Type     types.Type
	TypeRef  string
	Parsable bool
	Traits   TraitInstances
}

func (td *TraitDesc) extractUnderlying() (underlying, bool) {
	v, ok := td.Type.Underlying().(*types.Basic)
	if !ok {
		return Unknown, false
	}
	switch v.Kind() {
	case
		types.UntypedInt,
		types.Int,
		types.Int8,
		types.Int16,
		types.Int32,
		types.Int64:
		return Int64, true
	case
		types.Uint,
		types.Uint8,
		types.Uint16,
		types.Uint32,
		types.Uint64:
		return Uint64, true
	case
		types.Float32:
		return Float32, true
	case
		types.UntypedFloat,
		types.Float64:
		return Float64, true
	case
		// untyped strings dont need casting
		// so we can skip them here
		// types.UntypedString,
		types.String:
		return String, true
	}
	return Unknown, true
}

func (td *TraitDesc) hasUnderlying(u underlying) bool {
	underlying, ok := td.extractUnderlying()
	if !ok {
		return false
	}
	return underlying == u
}

func implementsJSONUnmarshaler(td *TraitDesc) bool {
	iFace, err := gencommon.FindIFaceDef("encoding/json", "Unmarshaler")
	if err != nil || iFace == nil {
		panic("Failed to find encoding/json.Unmarshaler")
	}
	return gencommon.TypeImplements(td.Type, iFace)
}

func implementsYAMLUnmarshaler(td *TraitDesc) bool {
	iFace, err := gencommon.FindIFaceDef("gopkg.in/yaml.v3", "Unmarshaler")
	if err != nil || iFace == nil {
		panic("Failed to find gopkg.in/yaml.v3.Unmarshaler")
	}
	return gencommon.TypeImplements(td.Type, iFace)
}

func implementsTextUnmarshaler(td *TraitDesc) bool {
	iFace, err := gencommon.FindIFaceDef("encoding", "TextUnmarshaler")
	if err != nil || iFace == nil {
		panic("Failed to find encoding.TextUnmarshaler")
	}
	return types.Implements(td.Type, iFace)
}

// TraitInstances are a sortable slice of `TraitInstance`s.
type TraitInstances []TraitInstance

func (s TraitInstances) Len() int {
	return len(s)
}

func (s TraitInstances) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s TraitInstances) Less(i, j int) bool {
	// sort instances in the same order as values
	// so that they're entered into switch / map statements identically.
	return s[i].OwningValue.Less(s[j].OwningValue)
}

// TraitInstance is an instance of a trait.
type TraitInstance struct {
	OwningValue  Value
	value        string
	variableName string // optional; will be used if exists.
}

// Value safely returns a reference to a constant OR an absolute value.
func (t TraitInstance) Value() string {
	if len(t.variableName) > 0 && t.variableName != "_" {
		return t.variableName
	}
	return t.value
}
