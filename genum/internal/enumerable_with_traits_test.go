package internal_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/drshriveer/gtools/genum"
	"github.com/drshriveer/gtools/genum/gen"
	"github.com/drshriveer/gtools/genum/internal"
)

func TestGenerate_EnumerableWithTraits(t *testing.T) {
	generator := gen.Generate{
		InFile:          "./enumerable_with_traits.go",
		OutFile:         "./enumerable_with_traits.genum.go",
		Types:           []string{"EnumWithPackageImports"},
		GenJSON:         true,
		GenYAML:         true,
		GenText:         true,
		CaseInsensitive: true,
	}

	require.NoError(t, generator.Parse())
}

func TestEnumerableWithTraitsGeneration(t *testing.T) {
	tests := []struct {
		description   string
		enumName      string
		disableTraits bool
		expectError   bool
	}{
		{
			description: "EnumerableWithTraits parses successfully",
			enumName:    "EnumerableWithTraits",
			expectError: false,
		},
		{
			description: "Creatures parses successfully",
			enumName:    "Creatures",
			expectError: false,
		},
		{
			description: "ErrEnum1 fails due to inconsistent traits",
			enumName:    "ErrEnum1",
			expectError: true,
		},
		{
			description: "ErrEnum2 fails due to no trait names",
			enumName:    "ErrEnum2",
			expectError: true,
		},
		{
			description:   "ErrEnum1 succeeds with traits disabled",
			enumName:      "ErrEnum1",
			disableTraits: true,
			expectError:   false,
		},
		{
			description:   "ErrEnum2 succeeds with traits disabled",
			enumName:      "ErrEnum2",
			disableTraits: true,
			expectError:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			generator := gen.Generate{
				InFile:        "./enumerable_with_traits.go",
				OutFile:       "./enumerable_with_traits.genum.go",
				Types:         []string{test.enumName},
				GenJSON:       true,
				GenYAML:       true,
				GenText:       true,
				DisableTraits: test.disableTraits,
			}

			err := generator.Parse()
			if test.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestCreatures_Traits(t *testing.T) {
	tests := []struct {
		enum     internal.Creatures
		numLegs  int
		isMammal bool
	}{
		{enum: internal.NotCreature, numLegs: 0, isMammal: false},
		{enum: internal.Cat, numLegs: internal.CatLegs, isMammal: true},
		{enum: internal.Feline, numLegs: internal.CatLegs, isMammal: true},
		{enum: internal.Feline2, numLegs: internal.CatLegs, isMammal: true},
		{enum: internal.Dog, numLegs: internal.DogLegs, isMammal: true},
		{enum: internal.Ant, numLegs: internal.AntLegs, isMammal: false},
		{enum: internal.Spider, numLegs: internal.SpiderLegs, isMammal: false},
		{enum: internal.Human, numLegs: internal.HumanLegs, isMammal: true},
		{enum: internal.SeaAnemone, numLegs: 0, isMammal: false},
	}

	for _, test := range tests {
		t.Run(test.enum.String(), func(t *testing.T) {
			assert.Implements(t, (*genum.Enum)(nil), test.enum)
			assert.True(t, test.enum.IsValid())
			assert.Equal(t, test.numLegs, test.enum.NumCreatureLegs())
			assert.Equal(t, test.isMammal, test.enum.IsCreatureMammal())
		})
	}
}

func TestEnumerableWithTraits_Traits(t *testing.T) {
	tests := []struct {
		enum       internal.EnumerableWithTraits
		trait      string
		timeout    time.Duration
		typedTrait internal.OtherType
	}{
		{
			enum:       internal.E1,
			trait:      "trait 1",
			timeout:    5 * time.Minute,
			typedTrait: "OtherType0",
		},
		{
			enum:       internal.E2,
			trait:      "trait 2",
			timeout:    1 * time.Minute,
			typedTrait: "OtherType2",
		},
		{
			enum:       internal.E3,
			trait:      "trait 3",
			timeout:    2 * time.Minute,
			typedTrait: "OtherType3",
		},
	}

	for _, test := range tests {
		t.Run(test.enum.String(), func(t *testing.T) {
			assert.Implements(t, (*genum.Enum)(nil), test.enum)
			assert.True(t, test.enum.IsValid())
			assert.Equal(t, test.trait, test.enum.Trait())
			assert.Equal(t, test.timeout, test.enum.Timeout())
			assert.Equal(t, test.typedTrait, test.enum.TypedStringTrait())
		})
	}
}
