package gen

import (
	"errors"
	"github.com/fatih/structtag"
	"go/types"
	"sort"
	"strconv"

	"github.com/drshriveer/gtools/set"
)

type SorterDesc struct {
	TypeName string
	SortType string
	Fields   SortFieldDescs
}

// PriorityTree produces a lopsided tree that expresses how to compare values.
// Exposed for use in templates.
func (s SorterDesc) PriorityTree() *CompareLine {
	sort.Sort(s.Fields)
	result := &CompareLine{}
	current := result
	for i, v := range s.Fields {
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

func createSorterDesc(obj types.Object, typeName, sortableTypeName string) (*SorterDesc, error) {
	if obj == nil {
		return nil, errors.New(typeName + " was not found in AST")
	}

	strukt, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return nil, errors.New(typeName + " is not a struct")
	}

	// pull out tags and ordering info.
	sortFields := make(SortFieldDescs, 0)
	for i := 0; i < strukt.NumFields(); i++ {
		sfd, err := sortFieldDescFromTag(strukt.Field(i).Name(), strukt.Tag(i))
		if err != nil {
			return nil, err
		} else if sfd != nil {
			sortFields = append(sortFields, sfd)
		}
	}

	// validate.
	if err := sortFields.Validate(); err != nil {
		return nil, err
	}

	return &SorterDesc{
		TypeName: typeName,
		SortType: sortableTypeName,
		Fields:   sortFields,
	}, nil
}

type SortFieldDescs []*SortFieldDesc

func (s SortFieldDescs) Len() int {
	return len(s)
}
func (s SortFieldDescs) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortFieldDescs) Less(i, j int) bool {
	return s[i].Priority < s[j].Priority
}

func (s SortFieldDescs) Validate() error {
	// TODO: could add more diagnostic info here, but too lazy for now.
	known := set.Make[int]()
	for _, fd := range s {
		if !known.Add(fd.Priority) {
			return errors.New("multiple fields have same sort priority")
		}
	}

	return nil
}

type SortFieldDesc struct {
	FieldName      string
	CustomAccessor string
	Priority       int
}

func sortFieldDescFromTag(fName, tagLine string) (*SortFieldDesc, error) {
	tags, err := structtag.Parse(tagLine)
	if err != nil { // error returned when not found
		return nil, nil
	}

	sortTags, err := tags.Get("gsort")
	if err != nil { // error returned when not found
		return nil, nil
	}

	result := &SortFieldDesc{
		FieldName: fName,
	}
	result.Priority, err = strconv.Atoi(sortTags.Name)
	if err != nil {
		return nil, errors.New("first element must be an int indicating sort priority! fieldName: " + fName)
	}

	if len(sortTags.Options) == 1 {
		result.CustomAccessor = sortTags.Options[0]
	} else if len(sortTags.Options) > 1 {
		return nil, errors.New("too many gsort options found! fieldName: " + fName)
	}

	return result, nil
}

type CompareLine struct {
	Accessor string
	Nest     *CompareLine
}

func (c CompareLine) HasNest() bool {
	return c.Nest != nil
}
