//nolint:revive // test only
package internal

//go:generate genum -types=EnumerableWithParsableTraits -parsableByTraits=Parsable1,Parsable2,OtherEnum,TypedString

type EnumerableWithParsableTraits int

const (
	P1, _NonParsable, _Parsable1, _Parsable2, _OtherEnum, _TypedString = EnumerableWithParsableTraits(iota), "non-parsable", 1, "1", Enum1Value0, OtherType("typedStr1")
	P2, _, _, _, _, _                                                  = EnumerableWithParsableTraits(iota), "non-parsable", 2, "2", Enum1Value1, OtherType("typedStr2")
	P3, _, _, _, _, _                                                  = EnumerableWithParsableTraits(iota), "non-parsable", 3, "3", Enum1Value2, OtherType("typedStr3")
)
