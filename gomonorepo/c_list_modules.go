package gomonorepo

import (
	"context"
	"fmt"
)

// ListModulesCommand is the `list-modules` command instance which can be added
// to the main command line parser. See description below for details.
var ListModulesCommand = &listModCommand{
	EmbeddedCommand: EmbeddedCommand{
		CmdName: "list-modules",
		Short:   "List all go modules, and their dependencies also defined in the mono repo.",
	},
}

type listModCommand struct {
	EmbeddedCommand
	Format string `long:"format" choice:"info" choice:"location-only-new-line" default:"info"`
}

func (x *listModCommand) RunCommand(ctx context.Context, opts *AppOptions) error {
	modTree, err := listAllModules(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to list modules: %w", err)
	}

	switch x.Format {
	case "info":
		opts.Infof("Found %d go.mod files in directory %q:", len(modTree.AllModulesMap), opts.GetRoot())
		for _, mod := range modTree.AllModules {
			opts.Printf("\t - %s (%s)\n", mod.ModFile.Module.Mod.Path, mod.ModFilePath)
			for _, dep := range mod.DependsOn {
				opts.Printf("\t\t - %s", dep.ModFile.Module.Mod.Path)
			}
		}
	case "location-only-new-line":
		for _, mod := range modTree.AllModules {
			opts.Printf("%s", mod.ModFilePath)
		}
	}

	return nil
}
