package gomonorepo

import (
	"context"
	"os"
	"strings"
)

const lintModulesDesc = `Invoke lint on modules in the mono repo.
When a '--parent' argument is provided, the command will run lint against the 
modules that changed since the parent commit. 
Note: this command expects golangci-lint is installed.
`

var LintModulesCommand = &lintModulesCommand{
	EmbeddedCommand: EmbeddedCommand{
		CmdName: "lint",
		Short:   "Invoke lint command in the mono repo.",
		Long:    lintModulesDesc,
	},
}

type lintModulesCommand struct {
	EmbeddedCommand
	ParentCommitOpt

	Fags string `long:"flags" description:"Flags to pass to through to the lint command."`
}

func (x *lintModulesCommand) RunCommand(ctx context.Context, opts *GlobalOptions) error {
	_, mods, err := listAllChangedModules(ctx, opts, x.ParentCommit)
	if err != nil {
		return err
	}

	success, err := invokeOnModules(ctx, opts, mods.Slice(), x.testModule)
	if err != nil {
		return err
	}
	if !success {
		os.Exit(1)
	}
	return nil
}

func (x *lintModulesCommand) testModule(ctx context.Context, m *Module) (commandResult, error) {
	args := make([]string, 2, 5)
	args[0] = "golangci-lint"
	args[1] = "run"
	if x.Fags != "" {
		args = append(args, strings.Fields(x.Fags)...)
	}
	args = append(args, m.ModDirectory)
	return runCommand(ctx, args), nil
}
