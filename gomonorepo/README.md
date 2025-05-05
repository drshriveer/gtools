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
  -r, --root=        Root directory of the mono repo. (default: .)
  -v, --verbose      Enable verbose logging.
  -p, --parallelism= Permitted parallelism for tasks that can be parallelized. (default: 4)
  -t, --timeout=     Timeout for the command. (default: 5m)
  -x, --excludePath= Paths to to exclude from searches (these may be regex). (default: node_modules, vendor)

Help Options:
  -h, --help         Show this help message

Available commands:
  detect-changes  detect changed modules and their dependencies.
  fmt             Invoke format command in the mono repo.
  generate        Invoke go generate in the mono repo.
  lint            Invoke lint command in the mono repo.
  list-modules    List all go modules, and their dependencies also defined in the mono repo.
  test            Invoke go tests in the mono repo.

```

**Example: Add to your justfile:**

-	TODO!! (link to updated just-file)

**Example: Add to CI:**

-	TODO!! (link to updated ci yaml)

### TODO:

What remains / ideas of the future, time permitting.

-	improve documentation (after upgrading gtools golint-ci).
-	smarter dependency resolution and test/lint planning
	-	Option to run in dependency order (with parallelism where permitted)
	-	sub-package dependency graph
	-	test-package change detection and planning
-	parallelize sub-package generate commands
	-	Run in dependency order, parallelize what can be.
-	module release & upgrade flow
	-	semver change classification detection
	-	tool to update inter dependencies
-	config.yaml support (if/when configurations are complex enough to desire it)
-	skip directories (e.g. node_modules)
