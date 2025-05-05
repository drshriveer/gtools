package gomonorepo

import (
	"context"
	"fmt"
	"os"
)

// TidyModulesCommand is the `list-modules` command instance which can be added
// to the main command line parser. See description below for details.
var TidyModulesCommand = &tidyModCommand{
	EmbeddedCommand: EmbeddedCommand{
		CmdName: "tidy",
		Short:   "Run 'go mod tidy' on all go modules in the monorepo. If a go.work file is found, this will also be tidied.",
	},
}

type tidyModCommand struct {
	EmbeddedCommand
}

func (x *tidyModCommand) RunCommand(ctx context.Context, opts *AppOptions) error {
	focus, ok := opts.GetFocusDir()
	if ok {
		cr, err := x.runPerModTarget(ctx, focus)
		if err != nil {
			return err
		}
		cr.Print()
		if !cr.succeeded {
			os.Exit(1)
		}
		return nil
	}

	modTree, err := listAllModules(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to list modules: %w", err)
	}

	success, err := invokeOnElement(ctx, opts, modTree.AllModules, x.runPerMod)
	if err != nil {
		return err
	} else if !success {
		os.Exit(1)
	}

	success, err = invokeOnElement(ctx, opts, modTree.AllWorkFiles, x.runPerWork)
	if err != nil {
		return err
	} else if !success {
		os.Exit(1)
	}

	return nil
}

func (x *tidyModCommand) runPerModTarget(ctx context.Context, target string) (commandResult, error) {
	args := []string{
		"go",
		"mod",
		"tidy",
		"-C",
		target,
	}
	return runCommand(ctx, args), nil
}

func (x *tidyModCommand) runPerMod(ctx context.Context, m *Module) (commandResult, error) {
	return x.runPerModTarget(ctx, m.ModRoot)
}

func (x *tidyModCommand) runPerWork(ctx context.Context, m *WorkFile) (commandResult, error) {
	args := []string{
		"go",
		"work",
		"sync",
		"-C",
		m.WorkRoot,
	}
	return runCommand(ctx, args), nil
}
