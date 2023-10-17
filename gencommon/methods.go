package gencommon

import (
	"sort"
)

//go:generate gsort --types Method=Methods

// Exported returns ony exported methods.
func (m Methods) Exported() Methods {
	sort.Sort(m)
	for i, m_ := range m {
		if m_.IsExported {
			return m[i:]
		}
	}
	return m
}

// Private returns ony Private methods.
func (m Methods) Private() Methods {
	sort.Sort(m)
	for i, m_ := range m {
		if m_.IsExported {
			return m[:i]
		}
	}

	return m
}
