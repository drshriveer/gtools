GO_LINT_VERSION := '1.54.2'
MDFMT_VERSION := 'latest'
PKG_ROOT := `pwd`
MODS := `go list -f '{{.Dir}}' -m`
export PATH := env_var('PATH') + ':' + PKG_ROOT + '/bin'
export GOBIN := PKG_ROOT + "/bin"
export GOEXPERIMENT := "loopvar"
CURRENT_DIR := invocation_directory_native()

# Github Actions doesn't appreciate high parallelism... The rest of us develop on macos.

PARALLEL := if os() == "macos" { '8' } else { '1' }

# Runs `go mod tidy` for all modules in the current directory, then sync go workspaces.
tidy:
    @just _invokeMod "go mod tidy -C {}" "{{ CURRENT_DIR }}"
    go work sync

# Runs `go test --race ` for all modules in the current directory.
test:
    @just _invokeMod "go test --race {}/..." "{{ CURRENT_DIR }}"

# Runs lint and test for all modules in the current directory.
check: lint test

# Runs lint/format for all modules in the current directory.
lint: _tools-linter
    @just _invokeMod "golangci-lint run {}/..." "{{ CURRENT_DIR }}"

# Fixes all auto-fixable format and lint errors for all modules in the current directory.
fix: _tools-linter format-md
    just --fmt --unstable
    @just _invokeMod "golangci-lint run --fix {}/..." "{{ CURRENT_DIR }}"

# Formats markdown.
format-md:
    @go install github.com/moorereason/mdfmt@{{ MDFMT_VERSION }}
    mdfmt -w -l ./**/*.md

# Runs `go generate` on all modules in the current directory.
generate: _tools-generate
    @just _invokeMod "go generate -C {} ./..." "{{ CURRENT_DIR }}"

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

# Always rebuild the genum, gsort, and gerror executibles from this package direclty for testing.

# Other packages should use the _install-go-pkg script.
_tools-generate:
    go build -o bin/genum genum/cmd/genum/main.go
    go build -o bin/gsort gsort/cmd/gsort/main.go
    go build -o bin/gerror gerror/cmd/gerror/main.go

# installs a go package at the version indicaed in go.work / go.mod.
# This may break if we're using inconsistent versions across projects, but I don't think it will.

# If it does, we might consider picking the latest version, or maybe we just want it to break.
_install-go-pkg package cmdpath:
    #!/usr/bin/env bash
    pkgVersion=`go list -f '{{{{.Version}}' -m {{ package }}`
    echo "installing {{ package / cmdpath }}@$pkgVersion"
    go install {{ package / cmdpath }}@$pkgVersion

# a the placeholder `{}` which is the path to the correct module.
_invokeMod cmd target='all':
    #!/usr/bin/env bash
    if [ "{{ target }}" = "{{ PKG_ROOT }}" ]; then
      xargs -L1 -P {{ PARALLEL }} -t -I {} {{ cmd }} <<< "{{ MODS }}"
     else
      xargs -L1 -t -I {} {{ cmd }} <<< "{{ target }}"
    fi
