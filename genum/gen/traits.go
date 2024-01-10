package gen

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

// TraitDesc define a trait-- this is exposed for template use.
type TraitDesc struct {
	Name     string
	TypeRef  string
	Parsable bool
	Traits   TraitInstances
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
