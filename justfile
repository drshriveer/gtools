GO_LINT_VERSION := '1.54.2'
PKG_ROOT := `pwd`
MODS := `go list -f '{{.Dir}}' -m`
export PATH := env_var('PATH') + ':' + PKG_ROOT + '/bin'
export GOBIN := PKG_ROOT + "/bin"
CURRENT_DIR := invocation_directory_native()

# Runs `go mod tidy` all modules or a single specified target, then sync go workspaces.
tidy:
    @just _invokeMod "go mod tidy -C {}" "{{ CURRENT_DIR }}"
    go work sync

# Runs `go test --race ` all modules or a single specified target.
test:
    @just _invokeMod "go test --race {}/..." "{{ CURRENT_DIR }}"

# Runs lint and test on all modules.
check: lint test

# Runs lint and format checkers all modules or a single specified target.
lint: _tools-linter
    @just _invokeMod "golangci-lint run {}/..." "{{ CURRENT_DIR }}"

# Fixes all auto-fixable format and lint errors on all modules or a single specified target.
fix: _tools-linter
    just --fmt --unstable
    @just _invokeMod "golangci-lint run --fix {}/..." "{{ CURRENT_DIR }}"

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
generate: _tools-generate
    @just _invokeMod "go generate -C {} ./..." "{{ CURRENT_DIR }}"

# always rebuild the genum executable in this package for testing.
# other packages should use something like:

# @go install github.com/drshriveer/gtools/genum/genum@{{ GENUM_VERSION }}
_tools-generate:
    go build -o bin/genum genum/genum/main.go
    go build -o bin/gsort gsort/cmd/main.go
    go build -o bin/gerror gerror/cmd/main.go

# a the placeholder `{}` which is the path to the correct module.
[macos]
[unix]
[windows]
_invokeMod cmd target='all':
    #!/usr/bin/env bash
    if [ "{{ target }}" = "{{ PKG_ROOT }}" ]; then
      xargs -L1 -P 8 -t -I {} {{ cmd }} <<< "{{ MODS }}"
     else
      xargs -L1 -t -I {} {{ cmd }} <<< "{{ target }}"
    fi

# a the placeholder `{}` which is the path to the correct module.

# The linux has a reduced parallelism as his seems to cause issues with git actions.
[linux]
_invokeMod cmd target='all':
    #!/usr/bin/env bash
    if [ "{{ target }}" = "{{ PKG_ROOT }}" ]; then
      xargs -L1 -P 2 -t -I {} {{ cmd }} <<< "{{ MODS }}"
     else
      xargs -L1 -t -I {} {{ cmd }} <<< "{{ target }}"
    fi
