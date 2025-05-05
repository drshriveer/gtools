package gomonorepo

import (
	"context"
	"os"
	"strings"
)

const testModulesDesc = `Invoke go tests in the mono repo.
When a '--parent' argument is provided, the command will run tests against the 
modules that changed since the parent commit, and their dependencies recursively.
If any of the tests fail, the command will return a non-zero exit code.
Note: this command expects go and git are installed.
`

// TestModulesCommand is the `test` command instance which can be added
// to the main command line parser. See description above for details.
var TestModulesCommand = &testModulesCommand{
	EmbeddedCommand: EmbeddedCommand{
		CmdName: "test",
		Short:   "Invoke go tests in the mono repo.",
		Long:    testModulesDesc,
	},
}

type testModulesCommand struct {
	EmbeddedCommand
	ParentCommitOpt
	Fags []string `long:"flags" short:"f" description:"Flags to pass to through to the test command, these are passed in BEFORE the package list." default:"-race" default:"-count=1" default:"-cover"`
	Args string   `long:"args" description:"Pass the remainder of the command line (everything after -args) to the test binary, uninterpreted and unchanged. This is passed AFTER the package list."`
}

func (x *testModulesCommand) RunCommand(ctx context.Context, opts *AppOptions) error {
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

	_, mods, err := listAllChangedAndDependencies(ctx, opts, x.ParentCommit)
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

func (x *testModulesCommand) runPerModule(ctx context.Context, m *Module) (commandResult, error) {
	return x.runPerTarget(ctx, m.ModRoot)
}

func (x *testModulesCommand) runPerTarget(ctx context.Context, target string) (commandResult, error) {
	args := make([]string, 2, 5)
	args[0] = "go"
	args[1] = "test"
	args = append(args, x.Fags...)
	args = append(args, ensureRecursivePath(target))
	if x.Args != "" {
		args = append(args, "-args")
		args = append(args, strings.Fields(x.Args)...)
	}
	return runCommand(ctx, args), nil
}
