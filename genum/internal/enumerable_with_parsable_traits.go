//nolint:revive // test only
package internal

//go:generate genum -types=EnumerableWithParsableTraits -parsableByTraits=Parsable,AlsoSometimesNonParsable,OtherEnum

type EnumerableWithParsableTraits int

const (
	P1, _NonParsable, _SometimesNonParsable, _AlsoSometimesNonParsable, _Parsable, _OtherEnum = EnumerableWithParsableTraits(iota), "non-parsable", 3, 1, "1", Enum1Value0
	P2, _, _, _, _, _                                                                         = EnumerableWithParsableTraits(iota), "non-parsable", 2, 2, "2", Enum1Value1
	P3, _, _, _, _, _                                                                         = EnumerableWithParsableTraits(iota), "non-parsable", 1, 3, "3", Enum1Value2
)
