package gencommon_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/drshriveer/gtools/gencommon"
	"github.com/drshriveer/gtools/set"
)

func TestFindInterface(t *testing.T) {
	t.Parallel()
	tests := []struct {
		description     string
		expectedMethods []string
		options         gencommon.ParseIFaceOption
	}{
		{
			description:     "Default options",
			expectedMethods: []string{"ParentMethod", "BazMethod"},
		},
		{
			description:     "IncludePrivate options",
			options:         gencommon.IncludePrivate,
			expectedMethods: []string{"ParentMethod", "BazMethod", "pooMethod"},
		},
		{
			description: "IncludeEmbedded options",
			options:     gencommon.IncludeEmbedded,
			expectedMethods: []string{
				"ParentMethod",
				"BazMethod",
				"BarMethod",
				"AMethod",
				"BMethod",
				"CMethod",
				"DMethod",
			},
		},
		{
			description: "IncludeEmbedded & IncludePrivate options",
			options:     gencommon.IncludeEmbedded | gencommon.IncludePrivate,
			expectedMethods: []string{
				"ParentMethod",
				"BazMethod",
				"pooMethod",
				"BarMethod",
				"AMethod",
				"BMethod",
				"CMethod",
				"DMethod",
				"bPrivate",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			t.Parallel()
			pwd := filepath.Join(os.Getenv("PWD"), "internal", "test_types.go")
			pkgs, pkg, _, imports, err := gencommon.LoadPackages(pwd)
			require.NoError(t, err)
			iface, err := gencommon.FindInterface(imports, pkgs, pkg.PkgPath, "TypeToGenerate", test.options)
			require.NoError(t, err)

			assert.Equal(t,
				"// TypeToGenerate has a comment.\n// SecondLine of expected comment.",
				iface.Comments.String())

			expected := set.Make(test.expectedMethods...)
			for _, m := range iface.Methods {
				assert.Truef(t, expected.Remove(m.Name), "found duplicate OR unexpected method %q", m.Name)
				// check that CMethod has a comment:
				if m.Name == "CMethod" {
					assert.NotEmpty(t, m.Comments)
				}
			}
			assert.Emptyf(t, expected.Slice(), "methods that were not found!")

		})
	}
}
