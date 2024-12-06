package set

type anyUint interface {
	~uint64 | ~uint32 | ~uint16 | ~uint8 | ~uint
}

// BitSet is a "Set" that operates on a bit-flag type.
type BitSet[T anyUint] uint64

// MakeBitSet creates bit set from the provided items.
func MakeBitSet[T anyUint](items ...T) BitSet[T] {
	result := BitSet[T](0)
	for _, item := range items {
		result |= BitSet[T](item)
	}
	return result
}

// Add returns true if any were added.
func (s *BitSet[T]) Add(items ...T) bool {
	added := false
	resultS := *s
	for _, item := range items {
		asFlag := BitSet[T](item)
		added = added || resultS&asFlag != asFlag
		resultS |= asFlag
	}
	*s = resultS
	return added
}

// Remove returns true if any were removed.
func (s *BitSet[T]) Remove(items ...T) bool {
	removed := false
	resultS := *s
	for _, item := range items {
		asFlag := BitSet[T](item)
		removed = removed || resultS&asFlag == asFlag
		resultS &= ^asFlag
	}

	*s = resultS
	return removed
}

// MaskOf returns a new BitSet containing only the flags that are set in the provided value.
func (s BitSet[T]) MaskOf(in T) BitSet[T] {
	return s & BitSet[T](in)
}

// Has returns true if the provided flag is set.
func (s BitSet[T]) Has(flag T) bool {
	asFlag := BitSet[T](flag)
	return s&asFlag == asFlag
}

// HasAny returns true if any of the flags are present.
func (s BitSet[T]) HasAny(flags ...T) bool {
	for _, flag := range flags {
		if s.Has(flag) {
			return true
		}
	}
	return false
}
