gomonorepo
==========

Is a commandline tool that assists in managing a golang monorepo!

### Features

-	invoke test/lint/format on only the modules that have changed (and modules that depend on them).
-	skip directories based on `.gitignore`

### Getting started

##### Requires:

-	`git`
-	`golanglint-ci`
-	`go` (of course)

##### Install:

```bash
go install github.com/drshriveer/gtool/gomonorepo/cmd/gomonorepo@latest
```

##### Usage:

Once installed, run `gomonorepo --help`:

```
Usage:
  gomonorepo [OPTIONS] <command>

Application Options:
  -r, --root=          Root directory of the mono repo. (default: .)
      --invocationDir= If invocationDir is not the root directory, invocation will be limited to the invocationDir and its subdirectories.
  -v, --verbose        Enable verbose logging.
  -p, --parallelism=   Permitted parallelism for tasks that can be parallelized. (default: 4)
  -t, --timeout=       Timeout for the command. (default: 5m)
  -x, --excludePath=   Paths to to exclude from searches (these may be regex). Note: Anything excluded by git is ignored by default. (default: node_modules, vendor)

Help Options:
  -h, --help           Show this help message

Available commands:
  detect-changes  detect changed modules and their dependencies.
  fmt             Invoke format command in the mono repo.
  generate        Invoke go generate in the mono repo.
  lint            Invoke lint command in the mono repo.
  list-modules    List all go modules, and their dependencies also defined in the mono repo.
  test            Invoke go tests in the mono repo.
  tidy            Run 'go mod tidy' on all go modules in the monorepo. If a go.work file is found, this will also be tidied.
  update-pkgs     Update all modules containing the packages to the version specified.
```

**Example: Add to your justfile:**

-	See the [justfile](../justfile) in the current repo or:

```justfile
CURRENT_DIR := invocation_directory_native()

# Runs `go mod tidy` for all modules in the current directory, then sync go workspaces.
tidy: _tools-monorepo
    gomonorepo tidy --invocationDir={{ CURRENT_DIR }}

# Runs `go test --race ` for all modules in the current directory.
test: _tools-monorepo
    gomonorepo test --parent=origin/main --invocationDir={{ CURRENT_DIR }}

# Runs lint/format for all modules in the current directory.
lint: _tools-monorepo _tools-linter
    gomonorepo lint --parent=origin/main --invocationDir={{ CURRENT_DIR }}

# Fixes all auto-fixable format and lint errors for all modules in the current directory. 
fix: _tools-monorepo _tools-linter
    gomonorepo lint --parent=origin/main -f="--fix" --invocationDir={{ CURRENT_DIR }}
```

**Example: Add to CI:**

-	TODO!! (link to updated ci yaml)

### TODO:

What remains / ideas of the future, time permitting.

-	improve documentation (after upgrading gtools golint-ci).
-	smarter dependency resolution and test/lint planning
	-	Option to run in dependency order (with parallelism where permitted)
	-	sub-package dependency graph
	-	Planning should account for test-package (only) changes
	-	Planning should account for comment (only) changes
	-	Planning should account for non-go files (markdown, yaml, etc)
		-	account for 'embed' directives / embedded files must trigger update
	-	Planning should account for links that go off the dependency graph (e.g. http/grpc/etc)
		-	maybe define our own directives for this? (e.g. `//go:gomonorepo run-with-changes-to --package=${} --module=${} --scopeTo="./sub-pkg"`\)
-	parallelize sub-package generate commands
	-	Run in dependency order, parallelize what can be.
-	module release & upgrade flow
	-	semver change classification detection
	-	tool to update inter dependencies
-	config.yaml support (if/when configurations are complex enough to desire it)-
