package gen

type Values []Value

func (s Values) Len() int {
	return len(s)
}

func (s Values) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Values) Less(i, j int) bool {
	if s[i].Signed || s[j].Signed {
		v1, v2 := int64(s[i].Value), int64(s[j].Value)
		if v1 == v2 {
			return s[i].Name < s[j].Name
		}
		return v1 < v2
	}
	v1, v2 := s[i].Value, s[j].Value
	if v1 == v2 {
		return s[i].Name < s[j].Name
	}
	return v1 < v2
}

func (s Values) ValueDeduplicatedSet() Values {
	if len(s) < 2 {
		return s
	}
	result := make(Values, 0, len(Values{}))
	result = append(result, s[0])
	lastValue := s[0].Value
	for i := 1; i < len(s); i++ {
		curr := s[i]
		if lastValue != curr.Value {
			result = append(result, curr)
			lastValue = curr.Value
		}
	}
	return result
}

type Value struct {
	Name   string
	Value  uint64
	Signed bool
}
