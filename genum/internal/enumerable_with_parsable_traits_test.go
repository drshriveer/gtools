package internal_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drshriveer/gtools/genum/gen"
	"github.com/drshriveer/gtools/genum/internal"
)

func TestGenerate_EnumerableWithParsableTraits(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc           string
		parsableTraits []string
		expectErr      bool
	}{
		{
			desc: "No parsable traits should function as usual",
		},
		{
			desc:           "One parsable trait",
			parsableTraits: []string{"Parsable"},
		},
		{
			desc:           "Two parsable traits",
			parsableTraits: []string{"Parsable", "AlsoParsable"},
		},
		{
			desc:           "Non-existent trait is ignored",
			parsableTraits: []string{"Parsable", "AlsoParsable", "NotARealTrait"},
		},
		{
			desc:           "Trait with non-unique values throws an error",
			parsableTraits: []string{"NonParsable"},
			expectErr:      true,
		},
		{
			desc:           "Sometimes non parsable trait is fine on its own",
			parsableTraits: []string{"SometimesNonParsable"},
		},
		{
			desc:           "Sometimes non parsable trait throws an error when it overlaps with another trait",
			parsableTraits: []string{"SometimesNonParsable", "AlsoSometimesNonParsable"},
			expectErr:      true,
		},
	}
	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {

			generator := gen.Generate{
				InFile:           "./enumerable_with_parsable_traits.go",
				OutFile:          "./enumerable_with_parsable_traits.genum.go",
				Types:            []string{"EnumerableWithParsableTraits"},
				GenJSON:          true,
				GenYAML:          true,
				GenText:          true,
				ParsableByTraits: test.parsableTraits,
			}

			err := generator.Parse()
			if test.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
			}

		})

	}

}

func TestEnumerableWithParsableTraits_Parser(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc     string
		val      any
		expected internal.EnumerableWithParsableTraits
	}{
		{desc: "parse P1 by string trait", val: "1", expected: internal.P1},
		{desc: "parse P1 by int trait", val: 1, expected: internal.P1},
		{desc: "parse P1 by enum string value", val: "P1", expected: internal.P1},

		{desc: "parse P2 by string trait", val: "2", expected: internal.P2},
		{desc: "parse P2 by int trait", val: 2, expected: internal.P2},
		{desc: "parse P2 by enum string value", val: "P2", expected: internal.P2},

		{desc: "parse P3 by string trait", val: "3", expected: internal.P3},
		{desc: "parse P3 by int trait", val: 3, expected: internal.P3},
		{desc: "parse P3 by enum string value", val: "P3", expected: internal.P3},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			res, err := internal.ParseEnumerableWithParsableTraits(test.val)
			require.NoError(t, err)
			assert.Equal(t, res, test.expected)
		})
	}
}
