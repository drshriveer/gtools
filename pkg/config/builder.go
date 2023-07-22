package config

import "github.com/drshriveer/gcommon/pkg/enum"

type Dimension[T enum.EnumLike] struct {
	Default  enum.Enum[T]
	FlagName string
	ParseEnv bool
}

type Options struct {
	Dimensions []Dimension
}
