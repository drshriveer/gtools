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
		EnumTypeNames: []string{"Creatures"},
		GenJSON:       true,
		GenYAML:       true,
		GenText:       true,
	}

	require.NoError(t, generator.Parse())
	require.NoError(t, generator.Write())
}

// FIXME!! gavin !! re-enable tests with traits.

//
// func TestCreatures_Traits(t *testing.T) {
// 	tests := []struct {
// 		enum     internal.Creatures
// 		numLegs  int
// 		isMammal bool
// 	}{
// 		{enum: internal.NotCreature, numLegs: 0, isMammal: false},
// 		{enum: internal.Cat, numLegs: internal.CatLegs, isMammal: true},
// 		{enum: internal.Feline, numLegs: internal.CatLegs, isMammal: true},
// 		{enum: internal.Dog, numLegs: internal.DogLegs, isMammal: true},
// 		{enum: internal.Ant, numLegs: internal.AntLegs, isMammal: false},
// 		{enum: internal.Spider, numLegs: internal.SpiderLegs, isMammal: false},
// 		{enum: internal.Human, numLegs: internal.HumanLegs, isMammal: true},
// 	}
//
// 	for _, test := range tests {
// 		t.Run(test.enum.String(), func(t *testing.T) {
// 			assert.Equal(t, test.numLegs, test.enum.NumCreatureLegs())
// 			assert.Equal(t, test.isMammal, test.enum.IsCreatureMammal())
// 		})
// 	}
// }
