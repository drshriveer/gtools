package internal_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/drshriveer/gcommon/pkg/genum/gen"
	"github.com/drshriveer/gcommon/pkg/genum/internal"
)

func TestEnumerableWithTraits(t *testing.T) {
	generator := gen.Generate{
		InFile:        "./enumerable_with_traits.go",
		OutFile:       "./enumerable_with_traits.genum.go",
		EnumTypeNames: []string{"CreaturesAlt"},
		GenJSON:       true,
		GenYAML:       true,
		GenText:       true,
	}

	require.NoError(t, generator.Parse())
}

func TestCreatures_Traits(t *testing.T) {
	tests := []struct {
		enum     internal.Creatures
		numLegs  int
		isMammal bool
	}{
		{enum: internal.NotCreature, numLegs: 0, isMammal: false},
		{enum: internal.Cat, numLegs: internal.Cat_NumLegs, isMammal: internal.Cat_IsMammal},
		{enum: internal.Feline, numLegs: internal.Cat_NumLegs, isMammal: internal.Cat_IsMammal},
		{enum: internal.Dog, numLegs: internal.Dog_NumLegs, isMammal: internal.Dog_IsMammal},
		{enum: internal.Ant, numLegs: internal.Ant_NumLegs, isMammal: internal.Ant_IsMammal},
		{enum: internal.Spider, numLegs: internal.Spider_NumLegs, isMammal: internal.Spider_IsMammal},
		{enum: internal.Human, numLegs: internal.Human_NumLegs, isMammal: internal.Human_IsMammal},
	}

	for _, test := range tests {
		t.Run(test.enum.String(), func(t *testing.T) {
			assert.Equal(t, test.numLegs, test.enum.NumLegs())
			assert.Equal(t, test.isMammal, test.enum.IsMammal())
		})
	}
}
