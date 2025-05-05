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
	Fags string `long:"flags" description:"Flags to pass to through to the test command." default:"-race -count=1 -cover"`
}

func (x *testModulesCommand) RunCommand(ctx context.Context, opts *AppOptions) error {
	_, mods, err := listAllChangedAndDependencies(ctx, opts, x.ParentCommit)
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

func (x *testModulesCommand) testModule(ctx context.Context, m *Module) (commandResult, error) {
	args := make([]string, 2, 5)
	args[0] = "go"
	args[1] = "test"
	if x.Fags != "" {
		args = append(args, strings.Fields(x.Fags)...)
	}
	args = append(args, m.ModRoot)
	return runCommand(ctx, args), nil
}
