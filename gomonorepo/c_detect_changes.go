package gomonorepo

import (
	"context"
	"errors"
)

// DetectedChangesCommand is the `detect-changes` command instance which can be added
// to the main command line parser. See description below for details.
var DetectedChangesCommand = &detectChanges{
	EmbeddedCommand: EmbeddedCommand{
		CmdName: "detect-changes",
		Short:   "detect changed modules and their dependencies.",
	},
}

type detectChanges struct {
	EmbeddedCommand
	ParentCommitOpt
}

func (x *detectChanges) RunCommand(ctx context.Context, opts *AppOptions) error {
	if x.ParentCommit == "" {
		return errors.New("no parent commit specified; use the --parent flag")
	}
	opts.Verbose = true
	_, _, err := listAllChangedAndDependencies(ctx, opts, x.ParentCommit)
	return err
}
