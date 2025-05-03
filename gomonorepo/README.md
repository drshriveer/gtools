gomonorepo
==========

Is a commandline tool that assists in managing a golang monorepo!

### Features

-	invoke test/lint/format on only the modules that have changed (and modules that depend on them).

### Getting started

Install with:

```bash
go install github.com/drshriveer/gtool/gomonorepo/cmd/gomonorepo@latest
```

##### Usage:

Once installed, run `gomonorepo --help`:

```
Usage:
  main [OPTIONS] <command>

Application Options:
  -r, --root=        Root directory of the mono repo. (default: .)
  -v, --verbose      Enable verbose logging.
  -p, --parallelism= Permitted parallelism for tasks that can be parallelized. (default: 4)
  -t, --timeout=     Timeout for the command. (default: 5m)

Help Options:
  -h, --help         Show this help message

Available commands:
  fmt                   Invoke format command in the mono repo.
  generate              Invoke go generate in the mono repo.
  lint                  Invoke lint command in the mono repo.
  list-dependency-tree  List the dependency structure of modules in the monorepo.
  list-modules          Recursively list all go modules, and their dependencies also defined in the mono repo.
  test                  Invoke go tests in the mono repo.

```

**Example: Add to your justfile:**

- TODO!! (link to updated just-file) 

**Example: Add to CI:**

- TODO!! (link to updated ci yaml)

### TODO:

-	improve documentation (after upgrading gtools golint-ci).
-   make the lint pass.
-   tool to update inter dependencies