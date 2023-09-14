//nolint:revive // test only
package internal

//go:generate genum -types=EnumerableWithTraits,Creatures

import (
	stupidTime "time"
)

type OtherType string

type EnumerableWithTraits int

const (
	E1, _Trait, _Timeout, _TypedStringTrait = EnumerableWithTraits(iota), "trait 1", 5 * stupidTime.Minute, OtherType("OtherType0")
	E2, _, _, _                             = EnumerableWithTraits(iota), "trait 2", 1 * stupidTime.Minute, OtherType("OtherType2")
	E3, _, _, _                             = EnumerableWithTraits(iota), "trait 3", 2 * stupidTime.Minute, OtherType("OtherType3")
)

type Creatures int

const (
	NotCreature, _NumCreatureLegs, _IsCreatureMammal = Creatures(iota), 0, false
	Cat, CatLegs, _                                  = Creatures(iota), 4, true
	Dog, DogLegs, _                                  = Creatures(iota), 4, true
	Ant, AntLegs, _                                  = Creatures(iota), 6, false
	Spider, SpiderLegs, _                            = Creatures(iota), 8, false
	Human, HumanLegs, _                              = Creatures(iota), 2, true
	// Feline traits will be ignored in favor of cat traits.
	Feline, _, _ = Cat, 5, false
	Feline2      = Cat
	SeaAnemone   = Creatures(iota)
)

// This enum fails to generate because it has an inconsisitnet number of traits.
type ErrEnum1 int

const (
	ErrEnum1V1, _Trait1, _Trait2 = ErrEnum1(iota), 1, "hi"
	ErrEnum1V2, _                = ErrEnum1(iota), 1
)

// This enum fails to generate because it has no trait names.
type ErrEnum2 int

const (
	ErrEnum2V1, _, _ = ErrEnum2(iota), 1, "hi"
	ErrEnum2V2, _, _ = ErrEnum2(iota), 1, "bye"
)
