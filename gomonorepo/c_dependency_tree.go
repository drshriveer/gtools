package gomonorepo

import (
	"context"
	"fmt"
)

var ListDependencyTree = &listDepCommand{
	EmbeddedCommand: EmbeddedCommand{
		CmdName: "list-dependency-tree",
		Short:   "List the dependency structure of modules in the monorepo.",
	},
}

type listDepCommand struct {
	EmbeddedCommand
}

func (x *listDepCommand) RunCommand(ctx context.Context, opts *GlobalOptions) error {
	modTree, err := listAllModules(ctx, opts.GetRoot())
	if err != nil {
		return fmt.Errorf("failed to list modules: %w", err)
	}

	opts.Infof("Found %d go.mod files in directory %q:", len(modTree.AllModules), opts.GetRoot())
	for _, mod := range modTree.AllModules {
		opts.Printf("\t - %q has %d dependencies\n", mod.Mod.Module.Mod.Path, len(mod.DependsOn))
		for _, dep := range mod.DependsOn {
			opts.Printf("\t\t - " + dep.Mod.Module.Mod.Path + "\n")
		}
	}
	return nil
}
