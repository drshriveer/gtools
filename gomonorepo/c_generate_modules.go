package gomonorepo

import (
	"context"
	"os"
)

const generateModulesDesc = `Invoke go generate in the mono repo.
When a '--parent' argument is provided, the command will run generate against the 
modules that changed since the parent commit... Which is not necessarily what you want 
when following the go:generate directive; consider not supplying this flag.
Note: this command expects go and git are installed.
`

var GenerateModulesCommand = &genModulesCommand{
	EmbeddedCommand: EmbeddedCommand{
		CmdName: "generate",
		Short:   "Invoke go generate in the mono repo.",
		Long:    generateModulesDesc,
	},
}

type genModulesCommand struct {
	EmbeddedCommand
	ParentCommitOpt
}

// TODO: This command isn't smart yet, it would be friggen wonderful if we could search for
// `go:generate` directives and undestand if their templates or the underlying references changed.
func (x *genModulesCommand) RunCommand(ctx context.Context, opts *GlobalOptions) error {
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

func (x *genModulesCommand) testModule(ctx context.Context, m *Module) (commandResult, error) {
	args := []string{
		"go",
		"generate",
		"-C",
		m.ModDirectory,
		"-x",
		"./...",
	}
	return runCommand(ctx, args), nil
}
