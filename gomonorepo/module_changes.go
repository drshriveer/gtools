package gomonorepo

import (
	"context"

	"github.com/drshriveer/gtools/set"
)

func listAllChangedModules(
	ctx context.Context,
	opts *AppOptions,
	parentCommit string,
) (*ModuleTree, set.Set[*Module], error) {
	tree, err := listAllModules(ctx, opts)
	if err != nil {
		return nil, nil, err
	}
	mods, err := listAllChangedModulesWithTree(ctx, opts, parentCommit, tree)
	opts.Infof("Detected changes in %d/%d modules will run the command.\n",
		len(mods), len(tree.AllModules))

	return tree, mods, err
}

func listAllChangedModulesWithTree(
	ctx context.Context,
	opts *AppOptions,
	parentCommit string,
	tree *ModuleTree,
) (set.Set[*Module], error) {
	if parentCommit == "" {
		opts.Infof("No parent commit indicated, will run command on all %d modules.\n", len(tree.AllModules))
		return set.Make(tree.AllModules...), nil
	}

	changedFiles, err := listChangedFiles(ctx, opts, parentCommit, false)
	if err != nil {
		return nil, err
	}
	if len(changedFiles) == 0 {
		remote, current, err := getCurrentBranch(ctx)
		if err != nil {
			return nil, err
		}
		if current == parentCommit || remote+"/"+current == parentCommit {
			opts.Infof("Currently on %q with no changes; will run command on all %d moduels.\n",
				current,
				len(tree.AllModules))
			return set.Make(tree.AllModules...), nil
		}
		opts.Infof("No changes detected between %q and %q (HEAD). Nothing to do.\n", parentCommit, current)
		return nil, nil
	}

	changedMods := make(set.Set[*Module], len(tree.AllModules))
	for _, f := range changedFiles {
		mod := tree.ModuleContainingFile(f)
		if mod != nil {
			changedMods.Add(mod)
		}
	}

	numChanged := len(changedMods)
	if opts.Verbose {
		opts.Infof("Detected changes in %d modules.\n", numChanged)
		for mod := range changedMods {
			opts.Printf("\t - %s\n", mod.ModFile.Module.Mod.Path)
		}
	}

	return changedMods, nil
}

func listAllChangedAndDependencies(
	ctx context.Context,
	opts *AppOptions,
	parentCommit string,
) (
	*ModuleTree,
	set.Set[*Module],
	error,
) {
	tree, err := listAllModules(ctx, opts)
	if err != nil {
		return nil, nil, err
	}
	modsToTest, err := listAllChangedAndDependenciesWithTree(ctx, opts, parentCommit, tree)

	return tree, modsToTest, err
}

func listAllChangedAndDependenciesWithTree(
	ctx context.Context,
	opts *AppOptions,
	parentCommit string,
	tree *ModuleTree,
) (
	set.Set[*Module],
	error,
) {
	if parentCommit == "" {
		opts.Infof("No parent commit indicated, will run command on all %d modules.\n", len(tree.AllModules))
		return set.Make(tree.AllModules...), nil
	}

	changedMods, err := listAllChangedModulesWithTree(ctx, opts, parentCommit, tree)
	if err != nil {
		return nil, err
	}
	numChanged := len(changedMods)

	// Add their dependencies
	for _, mod := range changedMods.Slice() {
		mod.AddDependants(changedMods)
	}

	opts.Infof("Detected changes in %d/%d modules, after including dependencies %d/%d modules will run the command.\n",
		numChanged, len(tree.AllModules), len(changedMods), len(tree.AllModules))
	if opts.Verbose {
		for mod := range changedMods {
			opts.Printf("\t - %s\n", mod.ModFile.Module.Mod.Path)
		}
	}

	return changedMods, nil
}
