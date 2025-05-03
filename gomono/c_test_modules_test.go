package gomono

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTestModulesCommand_RunCommand(t *testing.T) {
	t.Skip("this tests exists for debugging")
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	opts := &GlobalOptions{
		Root: "/Users/gs/Dev/gtools",
	}
	cmd := &lintModulesCommand{
		ParentCommitOpt: ParentCommitOpt{
			ParentCommit: "HEAD",
		},
	}
	require.NoError(t, cmd.RunCommand(ctx, opts))
}
