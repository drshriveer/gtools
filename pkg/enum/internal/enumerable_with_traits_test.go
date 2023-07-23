package internal_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/drshriveer/gcommon/pkg/enum/gen"
)

func TestEnumerableWithTraits(t *testing.T) {
	generator := gen.Generate{
		InFile:        "./enumerable_with_traits.go",
		OutFile:       "./enumerable_with_traits.genum.go",
		EnumTypeNames: []string{"EnumerableWithTraits"},
		GenJSON:       true,
		GenYAML:       true,
		GenText:       true,
	}

	require.NoError(t, generator.Parse())
	require.NoError(t, generator.Write())

}
