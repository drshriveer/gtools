package internal

//go:generate genum -types=CreaturesAlt
type CreaturesAlt int

const (
	NotCreaturesAlt, _NumLegs, _IsMammal, _Uint64Thing = CreaturesAlt(iota), 0, false, uint64(0)
	CatAlt, _, _, _                                    = CreaturesAlt(iota), 4, true, uint64(65320)
	AntAlt, _, _, _                                    = CreaturesAlt(iota), 6, false, uint64(320)
)
