package internal

import (
	stupidTime "time"
)

//go:generate genum -types=EnumerableWithTraits
type EnumerableWithTraits int

const (
	E1, E1_Trait, E1_Timeout = EnumerableWithTraits(iota), "trait 1", 5 * stupidTime.Minute
	E2, E2_Trait, E2_Timeout = EnumerableWithTraits(iota), "trait 2", 5 * stupidTime.Minute
	E3, E3_Trait, E3_Timeout = EnumerableWithTraits(iota), "trait 3", 5 * stupidTime.Minute
)
