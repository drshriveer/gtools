GO_LINT_VERSION := '1.54.2'
PKG_ROOT := `pwd`
INSTALLED_TOOLS := PKG_ROOT / "bin" / ".installed_tools"
MODS := `go list -f '{{.Dir}}' -m`
export PATH := env_var('PATH') + ':' + PKG_ROOT + '/bin'
export GOBIN := PKG_ROOT + "/bin"
export GOEXPERIMENT := "loopvar"
CURRENT_DIR := invocation_directory_native()

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
    @just _install-go-pkg "github.com/moorereason/mdfmt"
    @mdfmt -w -l ./**/*.md

# Runs `go generate` on all modules in the current directory.
generate: _tools-generate
    @just _invokeMod "go generate -C {} ./..." "{{ CURRENT_DIR }}"

_tools-linter:
    @just _tools-install "golangci-lint" "curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v{{ GO_LINT_VERSION }}"

# Always rebuild the genum, gsort, and gerror executibles from this package direclty for testing.

# Other packages should use the _install-go-pkg script.
_tools-generate:
    go build -o bin/genum genum/cmd/genum/main.go
    go build -o bin/gsort gsort/cmd/gsort/main.go
    go build -o bin/gerror gerror/cmd/gerror/main.go
    go build -o bin/gogenproto gogenproto/cmd/gogenproto/main.go

# installs a go package at the version indicated in go.work / go.mod.
_install-go-pkg package cmdpath="":
    #!/usr/bin/env bash
    set -euo pipefail # makes scripts act like justfiles (https://github.com/casey/just#safer-bash-shebang-recipes)
    pkgVersion=`go list -f '{{{{.Version}}' -m {{ package }}`
    pkgPath="{{ trim_end_match(package / cmdpath, '/') }}"
    just _tools-install {{ package }} "go install $pkgPath@$pkgVersion"

# Installs a given "tool" with command "cmd" provided, if it isn't already installed.
# If the command changes in any way the tool will be re-installed.
# The tool's name "tool" should be unique and is used to keep the dependency list clear

# of previous installs.
_tools-install tool cmd:
    #!/usr/bin/env bash
    set -euo pipefail # makes scripts act like justfiles (https://github.com/casey/just#safer-bash-shebang-recipes)
    mkdir -p {{ parent_directory(INSTALLED_TOOLS) }}
    touch {{ INSTALLED_TOOLS }}
    if grep -Fxq "{{ tool }} # {{ cmd }}" {{ INSTALLED_TOOLS }}
    then
      echo "[tool_install]: {{ tool }} already installed"
    else
      echo "[tool_install]: installing {{ tool }} with command {{ cmd }}"
      {{ cmd }}
    fi
    # Always refresh references to ensure the tools file is clean.
    echo "$(grep -v '{{ tool }}' {{ INSTALLED_TOOLS }})" > {{ INSTALLED_TOOLS }}
    echo "{{ tool }} # {{ cmd }}" >> {{ INSTALLED_TOOLS }}
    sort {{ INSTALLED_TOOLS }} -o {{ INSTALLED_TOOLS }}

# a the placeholder `{}` which is the path to the correct module.
_invokeMod cmd target='all':
    #!/usr/bin/env bash
    set -euo pipefail # makes scripts act like justfiles (https://github.com/casey/just#safer-bash-shebang-recipes)
    if [ "{{ target }}" = "{{ PKG_ROOT }}" ]; then
      xargs -L1 -P 8 -t -I {} {{ cmd }} <<< "{{ MODS }}"
     else
      xargs -L1 -P 8 -t -I {} {{ cmd }} <<< "{{ target }}"
    fi
