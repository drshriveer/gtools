package internal_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/drshriveer/gtools/gsort/gen"
	"github.com/drshriveer/gtools/gsort/internal"
)

func TestGenerate(t *testing.T) {
	g := gen.Generate{
		InFile:  "./sortable.go",
		OutFile: "sortable.gsort.go",
		Types:   map[string]string{"Sortable": "GSortSortable"},
	}
	require.NoError(t, g.Parse())
	require.NoError(t, g.Write())
}

func TestSortable_Sort(t *testing.T) {
	tests := []struct {
		description    string
		input          []internal.Sortable
		expectedOutput []internal.Sortable
	}{}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {

		})
	}
}
