package gomonorepo

import (
	"context"
	"os"
)

const lintModulesDesc = `Invoke lint on modules in the mono repo.
When a '--parent' argument is provided, the command will run lint against the 
modules that changed since the parent commit. 
Note: this command expects golangci-lint is installed.
`

// LintModulesCommand is the `lint` command instance which can be added
// to the main command line parser. See description above for details.
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

	Fags []string `long:"flags" short:"f" description:"Flags to pass to through to the lint command."`
}

func (x *lintModulesCommand) RunCommand(ctx context.Context, opts *AppOptions) error {
	focus, ok := opts.GetFocusDir()
	if ok {
		cr, err := x.runPerTarget(ctx, focus)
		if err != nil {
			return err
		}
		cr.Print()
		if !cr.succeeded {
			os.Exit(1)
		}
		return nil
	}

	_, mods, err := listAllChangedModules(ctx, opts, x.ParentCommit)
	if err != nil {
		return err
	}

	success, err := invokeOnElement(ctx, opts, mods.Slice(), x.runPerModule)
	if err != nil {
		return err
	}
	if !success {
		os.Exit(1)
	}
	return nil
}

func (x *lintModulesCommand) runPerModule(ctx context.Context, m *Module) (commandResult, error) {
	return x.runPerTarget(ctx, m.ModRoot)
}

func (x *lintModulesCommand) runPerTarget(ctx context.Context, dir string) (commandResult, error) {
	args := make([]string, 2, 5)
	args[0] = "golangci-lint"
	args[1] = "run"
	args = append(args, x.Fags...)
	args = append(args, dir)
	return runCommand(ctx, args), nil
}
