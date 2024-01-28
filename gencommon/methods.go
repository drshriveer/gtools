package gencommon

import (
	"sort"
)

//go:generate gsort --types Method

// Exported returns ony exported methods.
func (ms Methods) Exported() Methods {
	sort.Sort(ms)
	for i, m := range ms {
		if m.IsExported {
			return ms[i:]
		}
	}
	return ms
}

// Private returns ony Private methods.
func (ms Methods) Private() Methods {
	sort.Sort(ms)
	for i, m := range ms {
		if m.IsExported {
			return ms[:i]
		}
	}

	return ms
}
