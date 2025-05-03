package gomonorepo

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/drshriveer/gtools/gsync"
)

type Module struct {
	Mod          *modfile.File
	ModFilePath  string
	ModDirectory string
	DependsOn    []*Module
	DependencyOf []*Module

	// NestedModules is a list of modules that are nested within this module,
	// this does not imply a dependency relationship, but is important for detecting
	// if a file is contained within a module.
	NestedModules []*Module
}

// ContainsFile returns true if a file is contained within this module.
func (x *Module) ContainsFile(f string) bool {
	if !strings.HasPrefix(f, x.ModDirectory) {
		return false
	}
	// Ensure the file is not part of a nested module.
	for _, nested := range x.NestedModules {
		if strings.HasPrefix(f, nested.ModDirectory) {
			return false
		}
	}
	return true
}

func listAllModules(ctx context.Context, rootDir string) (*ModuleTree, error) {
	executor, done := gsync.NewSliceExecutor[*Module](ctx)
	defer done()

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && info.Name() == "go.mod" {
			data, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("failed to read go.mod file at %s: %w", path, err)
			}
			return executor.AddTask(func(context.Context) (*Module, error) {
				f, err := modfile.Parse(path, data, nil)
				if err != nil {
					return nil, fmt.Errorf("failed to parse go.mod file at %s: %w", path, err)
				}
				return &Module{
					Mod:          f,
					ModFilePath:  path,
					ModDirectory: filepath.Dir(path),
				}, nil
			})
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking the directory tree: %w", err)
	}
	mods, err := executor.WaitAndResult()
	if err != nil {
		return nil, fmt.Errorf("error executing tasks: %w", err)
	}

	return postProcess(mods, rootDir)
}

func postProcess(mods []*Module, rootDir string) (*ModuleTree, error) {
	root := &ModuleTree{
		rootDir:       rootDir,
		AllModules:    mods,
		AllModulesMap: make(map[string]*Module, len(mods)),
	}
	var err error
	for _, mod := range mods {
		err = root.AddModule(mod)
		if err != nil {
			return nil, err
		}
	}

	for _, mod := range mods {
		for _, dep := range mod.Mod.Require {
			if depMod, exists := root.AllModulesMap[dep.Mod.Path]; exists {
				mod.DependsOn = append(mod.DependsOn, depMod)
				depMod.DependencyOf = append(depMod.DependencyOf, mod)
			}
		}
	}

	return root, nil
}

type ModuleTree struct {
	AllModules    []*Module
	AllModulesMap map[string]*Module
	rootDir       string
	root          treeNode
}

func (r *ModuleTree) ModuleContainingFile(f string) *Module {
	// Possibly danger: this assumes the file path is relative to the rootDir.
	path := strings.Split(filepath.Dir(f), string(filepath.Separator))
	if len(path) > 0 && path[0] == "." {
		path = path[1:]
	}
	return r.root.findModuleContainingFile(path)
}

func (r *ModuleTree) AddModule(mod *Module) error {
	if _, exists := r.AllModulesMap[mod.Mod.Module.Mod.Path]; exists {
		return fmt.Errorf("duplicate module path found: %q", mod.Mod.Module.Mod.Path)
	}
	r.AllModulesMap[mod.Mod.Module.Mod.Path] = mod

	temp := strings.TrimPrefix(mod.ModDirectory, r.rootDir+"/")
	path := strings.Split(temp, string(filepath.Separator))
	r.root.addModuleWithPath(mod, path)
	return nil
}

type treeNode struct {
	directoryName string
	moduleAtPath  *Module
	children      []*treeNode
}

func (r *treeNode) addModuleWithPath(mod *Module, path []string) {
	if len(path) == 0 {
		r.moduleAtPath = mod
		return
	}

	if r.moduleAtPath != nil {
		r.moduleAtPath.NestedModules = append(r.moduleAtPath.NestedModules, mod)
	}

	node := r.nodeWithDirName(path[0])
	if node == nil {
		node = &treeNode{directoryName: path[0]}
		r.children = append(r.children, node)
	}
	node.addModuleWithPath(mod, path[1:])
}

func (r *treeNode) findModuleContainingFile(path []string) *Module {
	if len(path) == 0 {
		return r.moduleAtPath
	}
	node := r.nodeWithDirName(path[0])
	if node == nil {
		// either we walked the tree and didn't find a match,
		// or we found a match but it has no children,
		// which means we do not need to search any further
		// and can find the nearest parent module.
		return r.moduleAtPath
	} else if len(node.children) == 0 {
		return node.moduleAtPath
	}

	// Or keep searching down the tree
	return node.findModuleContainingFile(path[1:])
}

func (r *treeNode) nodeWithDirName(dirName string) *treeNode {
	for _, node := range r.children {
		if node.directoryName == dirName {
			return node
		}
	}
	return nil
}
