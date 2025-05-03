package gomonorepo

import (
	"context"
	"fmt"
)

var ListModulesCommand = &listModCommand{
	EmbeddedCommand: EmbeddedCommand{
		CmdName: "list-modules",
		Short:   "Recursively list all go modules, and their dependencies also defined in the mono repo.",
	},
}

type listModCommand struct {
	EmbeddedCommand
	Format string `long:"format" choice:"info" choice:"location-only-new-line" default:"info"`
}

func (x *listModCommand) RunCommand(ctx context.Context, opts *GlobalOptions) error {
	modTree, err := listAllModules(ctx, opts.GetRoot())
	if err != nil {
		return fmt.Errorf("failed to list modules: %w", err)
	}

	switch x.Format {
	case "info":
		opts.Infof("Found %d go.mod files in directory %q:", len(modTree.AllModulesMap), opts.GetRoot())
		for _, mod := range modTree.AllModules {
			opts.Printf("\t - " + mod.Mod.Module.Mod.Path + "\n")
			opts.Printf("\t   " + mod.ModFilePath + "\n")
		}
	case "location-only-new-line":
		for _, mod := range modTree.AllModules {
			opts.Printf(mod.ModFilePath)
		}
	}

	return nil
}
