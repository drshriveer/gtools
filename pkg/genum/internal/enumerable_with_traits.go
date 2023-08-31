package internal

import (
	stupidTime "time"
)

//go:generate genum -types=EnumerableWithTraits,Creatures
type EnumerableWithTraits int

const (
	E1, E1_Trait, E1_Timeout = EnumerableWithTraits(iota), "trait 1", 5 * stupidTime.Minute
	E2, E2_Trait, E2_Timeout = EnumerableWithTraits(iota), "trait 2", 5 * stupidTime.Minute
	E3, E3_Trait, E3_Timeout = EnumerableWithTraits(iota), "trait 3", 5 * stupidTime.Minute
)

type Creatures int

const (
	NotCreature                             = Creatures(iota)
	Cat, Cat_NumLegs, Cat_IsMammal          = Creatures(iota), 4, true
	Dog, Dog_NumLegs, Dog_IsMammal          = Creatures(iota), 4, true
	Ant, Ant_NumLegs, Ant_IsMammal          = Creatures(iota), 6, false
	Spider, Spider_NumLegs, Spider_IsMammal = Creatures(iota), 8, false
	Human, Human_NumLegs, Human_IsMammal    = Creatures(iota), 2, true
)
