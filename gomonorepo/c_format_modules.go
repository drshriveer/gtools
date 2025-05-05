package gomonorepo

import (
	"context"
	"os"
	"strings"
)

const formatModulesDesc = `Invoke golangci-lint fmt in the mono repo.
When a '--parent' argument is provided, the command will run lint against the 
modules that changed since the parent commit. 
Note: This command expects golangci-lint is installed at a version >= v2.0.0.
`

// FormatModulesCommand is the `fmt` command instance which can be added
// to the main command line parser. See description above for details.
var FormatModulesCommand = &formatModulesCommand{
	EmbeddedCommand: EmbeddedCommand{
		CmdName: "fmt",
		Short:   "Invoke format command in the mono repo.",
		Long:    formatModulesDesc,
	},
}

type formatModulesCommand struct {
	EmbeddedCommand
	ParentCommitOpt

	Fags string `long:"flags" description:"Flags to pass to through to the format command."`
}

func (x *formatModulesCommand) RunCommand(ctx context.Context, opts *AppOptions) error {
	_, mods, err := listAllChangedModules(ctx, opts, x.ParentCommit)
	if err != nil {
		return err
	}

	success, err := invokeOnElement(ctx, opts, mods.Slice(), x.fmtModule)
	if err != nil {
		return err
	}
	if !success {
		os.Exit(1)
	}
	return nil
}

func (x *formatModulesCommand) fmtModule(ctx context.Context, m *Module) (commandResult, error) {
	args := make([]string, 2, 5)
	args[0] = "golangci-lint"
	args[1] = "fmt"
	if x.Fags != "" {
		args = append(args, strings.Fields(x.Fags)...)
	}
	args = append(args, m.ModRoot)
	return runCommand(ctx, args), nil
}
