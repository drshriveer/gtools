package internal

import (
	stupidTime "time"
)

type OtherType string

//go:generate genum -types=EnumerableWithTraits,Creatures
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
	Feline, _, _                                     = Cat, 5, false
	Feline2                                          = Cat
	SeaAnemone                                       = Creatures(iota)
)
