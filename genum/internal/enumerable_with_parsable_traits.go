//nolint:revive // test only
package internal

//go:generate genum -types=EnumerableWithParsableTraits -parsableByTraits=Parsable,AlsoSometimesNonParsable

type EnumerableWithParsableTraits int

const (
	P1, _NonParsable, _SometimesNonParsable, _AlsoSometimesNonParsable, _Parsable = EnumerableWithParsableTraits(iota), "non-parsable", 3, 1, "1"
	P2, _, _, _, _                                                                = EnumerableWithParsableTraits(iota), "non-parsable", 2, 2, "2"
	P3, _, _, _, _                                                                = EnumerableWithParsableTraits(iota), "non-parsable", 1, 3, "3"
)
