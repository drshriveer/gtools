package set

import (
	"encoding/json"

	"gopkg.in/yaml.v3"
)

var setVal = struct{}{}

// FIXME LATER:
// 2. benchmark especially for Add(singleItem T, moreItems T) optimization.
// 3. consider separating them all... Add(single) AddAll(multi), Has, HasAll, HasAny

// Set is a generic se type.
type Set[T comparable] map[T]struct{}

// Make creates a set.
func Make[T comparable](items ...T) Set[T] {
	s := make(Set[T], len(items))
	for _, item := range items {
		s[item] = setVal
	}
	return s
}

// Slice returns a slice copy of this struct.
func (s Set[T]) Slice() []T {
	if len(s) == 0 {
		return nil
	}
	result := make([]T, len(s))
	i := 0
	for v := range s {
		result[i] = v
		i++
	}

	return result
}

// Add returns true if any were added.
func (s *Set[T]) Add(items ...T) bool {
	if *s == nil {
		*s = make(Set[T], len(items))
	}
	added := false
	for _, item := range items {
		if !added {
			if _, ok := (*s)[item]; !ok {
				added = true
			}
		}
		(*s)[item] = setVal
	}

	return added
}

// AddSet adds an input set to the current set and returns true if any were added.
// This is nil-safe.
func (s *Set[T]) AddSet(items Set[T]) bool {
	if *s == nil {
		*s = make(Set[T], len(items))
	}
	added := false
	for item := range items {
		if !added {
			if _, ok := (*s)[item]; !ok {
				added = true
			}
		}
		(*s)[item] = setVal
	}

	return added
}

// Remove returns true if any were removed.
func (s Set[T]) Remove(items ...T) bool {
	if len(s) == 0 {
		return false
	}
	removed := false
	for _, item := range items {
		if !removed {
			if _, ok := s[item]; ok {
				removed = true
			}
		}
		delete(s, item)
	}

	return removed
}

// RemoveSet removes the input set from the current set and returns true if any items were removed.
func (s Set[T]) RemoveSet(items Set[T]) bool {
	if len(s) == 0 {
		return false
	}
	removed := false
	for item := range items {
		if !removed {
			if _, ok := s[item]; ok {
				removed = true
			}
		}
		delete(s, item)
	}

	return removed
}

// Has returns true if ALL items are contained.
func (s Set[T]) Has(items ...T) bool {
	if len(s) == 0 || len(s) < len(items) {
		return false
	}
	for _, item := range items {
		if _, ok := s[item]; !ok {
			return false
		}
	}
	return true
}

// Has returns true if ANY items are contained.
func (s Set[T]) HasAny(items ...T) bool {
	if len(s) == 0 {
		return false
	}
	for _, item := range items {
		if _, ok := s[item]; ok {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json marshal interface.
func (s Set[T]) MarshalJSON() ([]byte, error) {
	v := s.Slice()
	return json.Marshal(v)
}

// UnmarshalJSON implements the json marshal interface.
func (s *Set[T]) UnmarshalJSON(data []byte) error {
	var v []T

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	s.Add(v...)
	return nil
}

// MarshalYAML implements the yaml marshal interface.
func (s Set[T]) MarshalYAML() (any, error) {
	return s.Slice(), nil
}

// UnmarshalYAML implements the yaml unmarshal interface.
func (s *Set[T]) UnmarshalYAML(value *yaml.Node) error {
	temp := []T{}
	err := value.Decode(&temp)
	if err != nil {
		return err
	}
	s.Add(temp...)
	return nil
}
