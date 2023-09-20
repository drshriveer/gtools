// Code generated by gerrors DO NOT EDIT.
package internal

// Sortables implements the sort.Sort interface for Sortable.
type Sortables []Sortable

func (s Sortables) Len() int {
	return len(s)
}
func (s Sortables) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s Sortables) Less(i, j int) bool {
	if s[i].Category.String() == s[j].Category.String() {
		if s[i].Property1 == s[j].Property1 {
			return s[i].Property2 < s[j].Property2
		}
		return s[i].Property1 < s[j].Property1
	}
	return s[i].Category.String() < s[j].Category.String()
}