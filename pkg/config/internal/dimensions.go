package internal

//go:generate genum -types=DimensionOne,DimensionTwo
type DimensionOne int

const (
	D1a DimensionOne = iota
	D1b
	D1c
	D1d
)

type DimensionTwo int

const (
	D2a DimensionTwo = iota
	D2b
	D2c
	D2d
	D2e
)
