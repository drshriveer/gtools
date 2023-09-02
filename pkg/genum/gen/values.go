package gen

// Values implements sort.Interface. Values must be consistently sorted to
// keep diffs to a minimum.
type Values []Value

func (s Values) Len() int {
	return len(s)
}
func (s Values) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Values) Less(i, j int) bool {
	return s[i].Less(s[j])
}

func (s Values) ValueDeduplicatedSet() Values {
	if len(s) < 2 {
		return s
	}
	result := make(Values, 0, len(Values{}))
	result = append(result, s[0])
	lastValue := s[0].Value
	addedDeprecated := s[0].IsDeprecated
	for i := 1; i < len(s); i++ {
		curr := s[i]
		if lastValue != curr.Value {
			result = append(result, curr)
			lastValue = curr.Value
			addedDeprecated = curr.IsDeprecated
		} else if addedDeprecated && !curr.IsDeprecated {
			result[len(result)-1] = curr
		}
	}
	return result
}

type Value struct {
	Name         string
	Value        uint64
	Signed       bool
	IsDeprecated bool
	Line         int
}

func (v Value) Less(vIn Value) bool {
	if v.Signed || vIn.Signed {
		v1, v2 := int64(v.Value), int64(vIn.Value)
		if v1 == v2 {
			return v.Name < vIn.Name
		}
		return v1 < v2
	}
	v1, v2 := v.Value, vIn.Value
	if v1 == v2 {
		return v.Name < vIn.Name
	}
	return v1 < v2
}
