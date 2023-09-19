package internal_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/drshriveer/gtools/gerrors/gen"
)

func TestSimpleEnumGeneration(t *testing.T) {
	generator := gen.Generate{
		InFile:  "./custom_error.go",
		OutFile: "./custom_error.gerror.go",
		Types:   []string{"GRPCError"},
	}

	require.NoError(t, generator.Parse())
}
