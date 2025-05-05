package gomonorepo

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"

	"github.com/drshriveer/gtools/gsync"
	"github.com/drshriveer/gtools/set"
)

// Module represents a Go module, its dependencies,
// and other parsed information from the module.
// This is likely to grow considerably in the future.
type Module struct {
	// ModFile is the parsed go.mod file.
	ModFile *modfile.File

	// ModFilePath is the path to the go.mod file, including the go.mod file suffix.
	ModFilePath string

	// ModRoot is the directory containing the go.mod file; or the root of the module.
	ModRoot string

	// DependsOn is a list of modules that this module depends on.
	DependsOn []*Module

	// DependencyOf is a list of modules that depend on this module.
	DependencyOf []*Module

	// NestedModules is a list of modules that are nested within this module,
	// this does not imply a dependency relationship, but is important for detecting
	// if a file is contained within a module.
	NestedModules []*Module
}

// AddDependants recursively adds all dependants of this module to the given set.
func (r *Module) AddDependants(s set.Set[*Module]) {
	for _, dependant := range r.DependencyOf {
		s.Add(dependant)
		// XXX: This may not actually be necessary since go.mod files generally list
		// indirect dependencies, but better safe than sorry.
		dependant.AddDependants(s)
	}
}

// ModuleTree represents a tree of Go modules starting from the root directory of the
// repository or mono repo.
type ModuleTree struct {
	// AllModules is a list of all modules found in the repository.
	AllModules []*Module
	// AllModulesMap is a map of all modules found in the repository, keyed by the module package name.
	AllModulesMap map[string]*Module
	// rootDir is the root directory of the repository.
	rootDir string

	// directoryTreeRoot represents the root of the module directory tree.
	// i.e. this is modules as they are laid out in the file system.
	directoryTreeRoot modDirTreeNode
}

// ModuleContainingFile returns the module that contains the given file.
// This is accomplished by walking the tree of modules and finding the nearest parent module.
// It is assumed that the file path is ALREADY relative to the rootDir of the module tree.
// The result *Module may be nil if no module was found.
func (r *ModuleTree) ModuleContainingFile(f string) *Module {
	// Possibly danger: this assumes the file path is relative to the rootDir.
	path := strings.Split(filepath.Dir(f), string(filepath.Separator))
	if len(path) > 0 && path[0] == "." {
		path = path[1:]
	}
	return r.directoryTreeRoot.findModuleContainingFile(path)
}

// AddModule adds a module to the module tree.
func (r *ModuleTree) AddModule(mod *Module) error {
	if _, exists := r.AllModulesMap[mod.ModFile.Module.Mod.Path]; exists {
		return fmt.Errorf("duplicate module path found: %q", mod.ModFile.Module.Mod.Path)
	}
	r.AllModulesMap[mod.ModFile.Module.Mod.Path] = mod

	temp := strings.TrimPrefix(mod.ModRoot, r.rootDir+"/")
	path := strings.Split(temp, string(filepath.Separator))
	r.directoryTreeRoot.addModuleWithPath(mod, path)
	return nil
}

// modDirTreeNode represents modules in the mono repo as they are laid out in the
// file system. E.g.
// .
// ├── foo
// │   └── github.com/foo  (module)
// └── bar
//
//	└── github.com/bar  (module)
type modDirTreeNode struct {
	// directoryName is the literal name of the root directory of the module.
	// i.e. this is not a path, just the name of the directory.
	directoryName string
	// moduleAtPath is the module present at this path, if any.
	moduleAtPath *Module
	// children is a list of modules  nodes in the module tree.
	children []*modDirTreeNode
}

// addModuleWithPath adds a module to the module tree at the given path.
func (r *modDirTreeNode) addModuleWithPath(mod *Module, path []string) {
	if len(path) == 0 {
		r.moduleAtPath = mod
		return
	}

	if r.moduleAtPath != nil {
		r.moduleAtPath.NestedModules = append(r.moduleAtPath.NestedModules, mod)
	}

	node := r.nodeWithDirName(path[0])
	if node == nil {
		node = &modDirTreeNode{directoryName: path[0]}
		r.children = append(r.children, node)
	}
	node.addModuleWithPath(mod, path[1:])
}

// findModuleContainingFile finds the module that contains the given file path,
// starting from the current node, following the directory tree downwards.
// The in put "path" is assumed to be a list of directory names, starting from the current node.
func (r *modDirTreeNode) findModuleContainingFile(path []string) *Module {
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

func (r *modDirTreeNode) nodeWithDirName(dirName string) *modDirTreeNode {
	for _, node := range r.children {
		if node.directoryName == dirName {
			return node
		}
	}
	return nil
}

// listAllModules lists all Go modules in the given directory and its subdirectories.
func listAllModules(ctx context.Context, opts *AppOptions) (*ModuleTree, error) {
	excludePaths, excludeDirs, err := opts.ExcludePathPatterns(ctx)
	if err != nil {
		return nil, err
	}

	executor, done := gsync.NewSliceExecutor[*Module](ctx)
	defer done()

	err = filepath.WalkDir(opts.GetRoot(), func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if matchesAny(excludePaths, excludeDirs, path) {
				if opts.Verbose {
					opts.Infof("Skipping excluded path: %s", path)
				}
				return filepath.SkipDir
			}
			return nil // recurse
		}

		if d.Name() == "go.mod" {
			return executor.AddTask(func(context.Context) (*Module, error) {
				// Sad no streaming parser for go.mod yet...
				data, err := os.ReadFile(path)
				if err != nil {
					return nil, fmt.Errorf("failed to read go.mod file at %s: %w", path, err)
				}
				f, err := modfile.Parse(path, data, nil)
				if err != nil {
					return nil, fmt.Errorf("failed to parse go.mod file at %s: %w", path, err)
				}
				return &Module{
					ModFile:     f,
					ModFilePath: path,
					ModRoot:     filepath.Dir(path),
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

	return buildModuleTree(mods, opts.GetRoot())
}

func buildModuleTree(mods []*Module, rootDir string) (*ModuleTree, error) {
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
		for _, dep := range mod.ModFile.Require {
			if depMod, exists := root.AllModulesMap[dep.Mod.Path]; exists {
				mod.DependsOn = append(mod.DependsOn, depMod)
				depMod.DependencyOf = append(depMod.DependencyOf, mod)
			}
		}
	}

	return root, nil
}
