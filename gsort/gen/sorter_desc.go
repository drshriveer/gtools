package gen

//go:generate gsort --types=SorterDesc,SortFieldDesc

import (
	"errors"
	"go/types"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/drshriveer/gtools/set"
)

// SorterDesc is a description of a sorter.
type SorterDesc struct {
	// The underlying type name (that we're making sortable).
	TypeName string `gsort:"SorterDescs,1"`

	// The name of the sortable type (sorted in case we support generating different sorters per type)
	SortTypeName string `gsort:"SorterDescs,2"`
	Fields       SortFieldDescs
}

// PriorityTree produces a lopsided tree that expresses how to compare values.
// Exposed for use in templates.
func (s SorterDesc) PriorityTree() *CompareLine {
	sort.Sort(s.Fields)
	result := &CompareLine{}
	current := result
	for i, v := range s.Fields {
		current.IsBool = v.FieldType.String() == "bool"
		current.Accessor = v.FieldName
		if len(v.CustomAccessor) > 0 {
			current.Accessor += "." + v.CustomAccessor
		}
		if len(s.Fields)-1 > i {
			current.Nest = &CompareLine{}
			current = current.Nest
		}
	}

	return result
}

func createSorterDesc(obj types.Object, typeName string) (SorterDescs, error) {
	if obj == nil {
		return nil, errors.New(typeName + " was not found in AST")
	}

	strukt, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return nil, errors.New(typeName + " is not a struct")
	}

	// pull out tags and ordering info.
	descs := make(map[string]*SorterDesc)
	sortFields := make(SortFieldDescs, 0)
	for i := 0; i < strukt.NumFields(); i++ {
		sField := strukt.Field(i)
		sfds, err := sortFieldDescFromTag(sField, strukt.Tag(i))
		if err != nil {
			return nil, err
		}
		for _, fd := range sfds {
			desc, ok := descs[fd.SortTypeName]
			if ok {
				desc.Fields = append(desc.Fields, fd)
			} else {
				desc = &SorterDesc{
					TypeName:     typeName,
					SortTypeName: fd.SortTypeName,
					Fields:       SortFieldDescs{fd},
				}
			}
			sort.Sort(desc.Fields)
			descs[fd.SortTypeName] = desc
		}
	}

	for _, desc := range descs {
		if err := desc.Fields.Validate(); err != nil {
			return nil, err
		}
	}

	result := make(SorterDescs, 0, len(sortFields))
	for _, desc := range descs {
		result = append(result, desc)
	}

	return result, nil
}

// Validate returns an error if anything is broken.
func (s SortFieldDescs) Validate() error {
	if len(s) == 0 {
		return errors.New("no sort attributes defined")
	}
	// TODO: could add more diagnostic info here, but too lazy for now.
	known := set.Make[int]()
	for _, fd := range s {
		if !known.Add(fd.Priority) {
			return errors.New("multiple fields have same sort priority")
		}
	}

	return nil
}

// SortFieldDesc describes a single field used for sorting.
type SortFieldDesc struct {
	FieldName      string
	FieldType      types.Type
	CustomAccessor string
	SortTypeName   string
	Priority       int `gsort:"SortFieldDescs,1"`
}

func sfdFromLine(options string) (*SortFieldDesc, error) {
	sfd := &SortFieldDesc{}
	tuple := strings.Split(options, ",")
	if len(tuple) < 2 {
		return nil, errors.New("minimum two tag options required; name of type to generate, and field priority")
	} else if len(tuple) > 3 {
		return nil, errors.New("maximum three tag options allowed; name of type to generate, field priority, optional accessor")
	}
	sfd.SortTypeName = tuple[0]

	var err error
	sfd.Priority, err = strconv.Atoi(tuple[1])
	if err != nil {
		return nil, errors.New("second option must be an int indicating sort priority! found: " + tuple[1])
	}
	if len(tuple) == 3 {
		sfd.CustomAccessor = tuple[2]
	}
	return sfd, nil
}

func sortFieldDescFromTag(sFiled *types.Var, tagLine string) ([]*SortFieldDesc, error) {
	remaining := reflect.StructTag(tagLine)
	result := make([]*SortFieldDesc, 0)
	for options, ok := remaining.Lookup("gsort"); ok; options, ok = remaining.Lookup("gsort") {
		sfd, err := sfdFromLine(options)
		if err != nil {
			return nil, err
		}
		sfd.FieldName = sFiled.Name()
		sfd.FieldType = sFiled.Type()
		result = append(result, sfd)
		remaining = reflect.StructTag(strings.Replace(string(remaining), `gsort:"`+options+`"`, "", 1))
	}
	return result, nil
}

// CompareLine is what's actually used by the template to generate if/else statements.
type CompareLine struct {
	// IsBool indicates this is a bool for making a different kind of comparison.
	IsBool bool
	// Accessor is how to access the field that sorts things.
	Accessor string
	// Nest is another if/template call to be nested in an if statement.
	Nest *CompareLine
}

// HasNest is a helper for templates.
func (c CompareLine) HasNest() bool {
	return c.Nest != nil
}

// String returns the Comparison like as a string.
func (c CompareLine) String() string {
	if c.IsBool {
		return "s[j]." + c.Accessor
	}
	return "s[i]." + c.Accessor + " < " + "s[j]." + c.Accessor
}
