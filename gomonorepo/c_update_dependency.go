package gomonorepo

import (
	"context"
	"os"
)

// UpdateDependencyCommand is the `lint` command instance which can be added
// to the main command line parser. See description above for details.
var UpdateDependencyCommand = &updateDependencyCommand{
	EmbeddedCommand: EmbeddedCommand{
		CmdName: "update-pkgs",
		Short:   "Update all modules containing the packages to the version specified.",
	},
}

type updateDependencyCommand struct {
	EmbeddedCommand
	Packages []string `long:"pkg" description:"Packages to upgrade across all modules that depend on it."`
}

type updateDep struct {
	mod          *Module
	pkgsToUpdate []string
}

func (x *updateDependencyCommand) RunCommand(ctx context.Context, opts *AppOptions) error {
	tree, err := listAllModules(ctx, opts)
	if err != nil {
		return err
	}
	updateDeps := make([]updateDep, 0)
	for _, mod := range tree.AllModules {
		var update updateDep
		for _, pkg := range x.Packages {
			if mod.Requires(pkg) {
				update.mod = mod
				update.pkgsToUpdate = append(update.pkgsToUpdate, pkg)
			}
		}
		if update.mod != nil {
			updateDeps = append(updateDeps, update)
		}
	}

	opts.Infof("Found %d modules to update %d packages.", len(updateDeps), len(x.Packages))

	success, err := invokeOnElement(ctx, opts, updateDeps, x.runPerModule)
	if err != nil {
		return err
	}
	if !success {
		os.Exit(1)
	}
	return nil
}

func (x *updateDependencyCommand) runPerModule(ctx context.Context, u updateDep) (commandResult, error) {
	var finalResult commandResult
	for i, pkg := range u.pkgsToUpdate {
		args := []string{
			"go",
			"get",
			"-C",
			u.mod.ModRoot,
			"-u",
			pkg,
		}
		if i == 0 {
			finalResult = runCommand(ctx, args)
		} else {
			cr := runCommand(ctx, args)
			finalResult.join(&cr)
		}
	}
	return finalResult, nil
}
