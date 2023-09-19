GO_LINT_VERSION := '1.54.1'
PKG_ROOT := `pwd`
MODS := `go list -f '{{.Dir}}' -m`
export PATH := env_var('PATH') + ':' + PKG_ROOT + '/bin'
export GOBIN := PKG_ROOT + "/bin"

# Runs `go mod tidy` all modules or a single specified target, then sync go workspaces.
tidy target='all':
    @just _invokeMod "go mod tidy -C {}" "{{ target }}"
    go work sync

# Runs `go test --race ` all modules or a single specified target.
test target='all':
    @just _invokeMod "go test --race {}/..." "{{ target }}"

# Runs lint and test on all modules.
check: lint test

# Runs lint and format checkers all modules or a single specified target.
lint target='all': _tools-linter
    @just _invokeMod "golangci-lint run {}/..." "{{ target }}"

# Fixes all auto-fixable format and lint errors on all modules or a single specified target.
fix target='all': _tools-linter
    just --fmt --unstable
    @just _invokeMod "golangci-lint run --fix {}/..." "{{ target }}"

_tools-linter:
    #!/usr/bin/env bash
    if command -v golangci-lint && golangci-lint --version | grep -q '{{ GO_LINT_VERSION }}'; then
      echo 'golangci-lint v{{ GO_LINT_VERSION }} already installed!'
    else
      echo "installing golangci-lint at version v{{ GO_LINT_VERSION }}"
      if test -e ./bin/golangci-lint; then
        rm ./bin/golangci-lint
      fi
      curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v{{ GO_LINT_VERSION }}
    fi

# Runs `go generate` on all go modules or a single specified target.
generate target='all': _tools-generate
    @just _invokeMod "go generate -C {} ./..." "{{ target }}"

# always rebuild the genum executable in this package for testing.
# other packages should use something like:

# @go install github.com/drshriveer/gtools/genum/genum@{{ GENUM_VERSION }}
_tools-generate:
    go build -o bin/genum genum/genum/main.go
    go build -o bin/gsort gsort/cmd/main.go

# a the placeholder `{}` which is the path to the correct module.
_invokeMod cmd target='all':
    #!/usr/bin/env bash
    if [ "{{ target }}" = "all" ]; then
      xargs -L1 -t -I {} {{ cmd }} <<< "{{ MODS }}"
     else
      xargs -L1 -t -I {} {{ cmd }} <<< "{{ target }}"
    fi
