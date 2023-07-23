package gen

import (
	"go/types"
	"sort"
	"strings"
)

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

type TraitDesc struct {
	Name          string
	Type          string
	PackageImport string
	Traits        TraitInstances
}

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

type TraitInstance struct {
	OwningValue  Value
	VariableName string

	// Includes all the type info.
	Type types.Type
}

// PackageImportPath returns the package import path of the underlying type;
// TODO: at some point we need to be able to resolve conflicts from inconsistent naming.
func (t TraitInstance) PackageImportPath() string {
	named, ok := t.Type.(*types.Named)
	if !ok {
		return ""
	}
	return named.Obj().Pkg().Path()
}

func (t TraitInstance) TypeRef() string {
	// I don't feel great about this, but some *types.Basic (string) come out as
	// "untyped something" so I guess we'll just trim that part off...
	// not sure how else to get the correct type.
	return strings.TrimPrefix(t.Type.String(), "untyped ")
}

func traitsFromMap(in map[string]TraitInstances) TraitDescs {
	result := make(TraitDescs, 0, len(in))
	for name, traits := range in {
		sort.Sort(traits)
		result = append(result, TraitDesc{
			Name:          name,
			Type:          traits[0].TypeRef(),
			PackageImport: traits[0].PackageImportPath(),
			Traits:        traits,
		},
		)
	}
	sort.Sort(result)
	return result
}
