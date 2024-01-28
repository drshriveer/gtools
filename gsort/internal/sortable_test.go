package internal_test

import (
	"path/filepath"
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

		expectedError error
		expectedFile  bool
	}{
		{
			description:  "sortable success",
			typeName:     "Sortable",
			expectedFile: true,
		},
		{
			description:  "multi sortable success",
			typeName:     "MultiSort",
			expectedFile: true,
		},
		{
			description:  "fails because there are no properties to sort",
			typeName:     "NotSortable",
			expectedFile: false,
		},
		// add more tests some day.
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()
			tempFile := filepath.Join(t.TempDir(), "sortable.gsort.go")
			g := gen.Generate{
				InFile:     "./sortable.go",
				OutFile:    tempFile,
				Types:      []string{test.typeName},
				UsePointer: true,
			}
			err := g.Parse()
			if test.expectedError != nil {
				require.ErrorIs(t, err, test.expectedError)
				return
			}
			require.NoError(t, err)
			require.NoError(t, g.Write())
			if test.expectedFile {
				assert.FileExists(t, tempFile)
			} else {
				assert.NoFileExists(t, tempFile)
			}
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
