package internal_test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/drshriveer/gtools/gsort/gen"
	"github.com/drshriveer/gtools/gsort/internal"
)

func TestGenerate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		description string
		typeName    string

		// TODO: use proper errors for this... I didn't have them when I wrote it,
		//  so not doing that now.
		expectedError bool
	}{
		{
			description: "sortable success",
			typeName:    "Sortable",
		},
		{
			description:   "fails because there are no properties to sort",
			typeName:      "NotSortable",
			expectedError: true,
		},
		// add more tests some day.
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()
			g := gen.Generate{
				InFile:     "./sortable.go",
				OutFile:    "sortable.gsort.go",
				Types:      map[string]string{test.typeName: test.typeName + "s"},
				UsePointer: true,
			}
			err := g.Parse()
			if test.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			// require.NoError(t, g.Write())
		})
	}
}

func TestSortable_Sort(t *testing.T) {
	t.Parallel()
	tests := []struct {
		description string
		input       internal.Sortables
		expected    internal.Sortables
	}{
		{
			description: "sort on category only",
			input: internal.Sortables{
				{Category: internal.BCategory},
				{Category: internal.CCategory},
				{Category: internal.ACategory},
			},
			expected: internal.Sortables{
				{Category: internal.ACategory},
				{Category: internal.BCategory},
				{Category: internal.CCategory},
			},
		},
		{
			description: "tie break with prop 1",
			input: internal.Sortables{
				{Category: internal.CCategory, Property1: "b"},
				{Category: internal.CCategory, Property1: "d"},
				{Category: internal.CCategory, Property1: "c"},
				{Category: internal.ACategory},
				{Category: internal.CCategory, Property1: "a"},
			},
			expected: internal.Sortables{
				{Category: internal.ACategory},
				{Category: internal.CCategory, Property1: "a"},
				{Category: internal.CCategory, Property1: "b"},
				{Category: internal.CCategory, Property1: "c"},
				{Category: internal.CCategory, Property1: "d"},
			},
		},
		{
			description: "tie break with prop 2",
			input: internal.Sortables{
				{Category: internal.CCategory, Property2: 2},
				{Category: internal.CCategory, Property2: 4},
				{Category: internal.CCategory, Property2: 3},
				{Category: internal.ACategory},
				{Category: internal.CCategory, Property2: 1},
			},
			expected: internal.Sortables{
				{Category: internal.ACategory},
				{Category: internal.CCategory, Property2: 1},
				{Category: internal.CCategory, Property2: 2},
				{Category: internal.CCategory, Property2: 3},
				{Category: internal.CCategory, Property2: 4},
			},
		},
		{
			description: "tie break layers with prop 2",
			input: internal.Sortables{
				{Category: internal.CCategory, Property1: "a", Property2: 2},
				{Category: internal.CCategory, Property1: "b", Property2: 3},
				{Category: internal.CCategory, Property1: "a", Property2: 4},
				{Category: internal.CCategory, Property1: "z", Property2: 3},
				{Category: internal.ACategory},
				{Category: internal.CCategory, Property2: 4},
				{Category: internal.CCategory, Property2: 1},
			},
			expected: internal.Sortables{
				{Category: internal.ACategory},
				{Category: internal.CCategory, Property2: 1},
				{Category: internal.CCategory, Property2: 4},
				{Category: internal.CCategory, Property1: "a", Property2: 2},
				{Category: internal.CCategory, Property1: "a", Property2: 4},
				{Category: internal.CCategory, Property1: "b", Property2: 3},
				{Category: internal.CCategory, Property1: "z", Property2: 3},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()
			sort.Sort(test.input)
			for i, expected := range test.expected {
				assert.Equal(t, expected, test.input[i])
			}
		})
	}
}
