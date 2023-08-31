package internal_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/drshriveer/gcommon/pkg/genum/gen"
)

func TestEnumerableWithTraits(t *testing.T) {
	generator := gen.Generate{
		InFile:        "./enumerable_with_traits.go",
		OutFile:       "./enumerable_with_traits.genum.go",
		EnumTypeNames: []string{"EnumerableWithTraits", "Creatures"},
		GenJSON:       true,
		GenYAML:       true,
		GenText:       true,
	}

	require.NoError(t, generator.Parse())
}
