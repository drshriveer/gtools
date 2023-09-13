package set

var setVal = struct{}{}

// FIXME LATER:
// 1. add tests
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
func (s *Set[T]) Slice() []T {
	// TODO: Benchmark if this is really more efficient.
	result := make([]T, len(*s))
	i := 0
	for v := range *s {
		result[i] = v
		i++
	}

	return result
}

// Add returns true if any were added.
func (s *Set[T]) Add(items ...T) bool {
	if s == nil {
		*s = make(Set[T])
	}
	// TODO: benchmark variations of this.
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
func (s *Set[T]) AddSet(items Set[T]) bool {
	if s == nil {
		*s = make(Set[T])
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
func (s *Set[T]) Remove(items ...T) bool {
	// TODO: benchmark variations of this.
	removed := false
	for _, item := range items {
		if !removed {
			if _, ok := (*s)[item]; ok {
				removed = true
			}
		}
		delete(*s, item)
	}

	return removed
}

// RemoveSet removes the input set from the current set and returns true if any items were removed.
func (s *Set[T]) RemoveSet(items Set[T]) bool {
	removed := false
	for item := range items {
		if !removed {
			if _, ok := (*s)[item]; ok {
				removed = true
			}
		}
		delete(*s, item)
	}

	return removed
}

// Has returns true if ALL items are contained.
func (s *Set[T]) Has(items ...T) bool {
	// TODO: use more efficient variation of this!
	if s == nil || len(*s) < len(items) {
		return false
	}
	for _, item := range items {
		if _, ok := (*s)[item]; !ok {
			return false
		}
	}
	return true
}
