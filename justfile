PKG_ROOT := `pwd`
INSTALLED_TOOLS := PKG_ROOT / "bin" / ".installed_tools"
export PATH := env_var('PATH') + ':' + PKG_ROOT + '/bin'
export GOBIN := PKG_ROOT + "/bin"
CURRENT_DIR := invocation_directory_native()
COVERAGE_DIR := PKG_ROOT / "go-coverage"

# Runs `go mod tidy` for all modules in the current directory, then sync go workspaces.
tidy: _tools-monorepo
    gomonorepo tidy --invocationDir={{ CURRENT_DIR }}

# Runs `go test --race ` for all modules in the current directory.
test: _tools-monorepo
    gomonorepo test --parent main --invocationDir={{ CURRENT_DIR }}

# Runs lint and test for all modules in the current directory.
check: lint test

# Runs lint/format for all modules in the current directory.
lint: _tools-monorepo _tools-linter
    gomonorepo lint --parent main --invocationDir={{ CURRENT_DIR }}

# Fixes all auto-fixable format and lint errors for all modules in the current directory.
fix: _tools-monorepo _tools-linter format-md
    gomonorepo lint --parent main -f="--fix" --invocationDir={{ CURRENT_DIR }}
    # just --fmt --unstable - Disabled due to combining single lines.


# Updates interdependent modules of gtools. TODO: could make this wayyy smarter.
update-interdependencies: _tools-monorepo && tidy
    gomonorepo update-pkgs \
        --pkg github.com/drshriveer/gtools/gencommon \
        --pkg github.com/drshriveer/gtools/genum \
        --pkg github.com/drshriveer/gtools/set \
        --pkg github.com/drshriveer/gtools/gomonorepo

# updates a single package across all go modules.
update-pkg pkgName: _tools-monorepo && tidy
    gomonorepo update-pkgs --pkg {{ pkgName }}

# Formats markdown.
format-md: (_install-go-pkg "github.com/moorereason/mdfmt")
    @mdfmt -w -l ./**/*.md

# Runs `go generate` on all modules in the current directory.
generate: _tools-monorepo _tools-generate
    gomonorepo generate --parent main --invocationDir={{ CURRENT_DIR }}

[no-cd]
coverage-go: _tools-monorepo
    mkdir -p {{COVERAGE_DIR}}
    rm -rf {{COVERAGE_DIR}}/*
    env CGO_ENABLED=1 gomonorepo test --invocationDir={{ CURRENT_DIR }} \
      -f="-race" \
      -f="-count=1" \
      -f="-cover" \
      --args="-test.gocoverdir {{COVERAGE_DIR}}"
    go tool covdata textfmt -i={{COVERAGE_DIR}} -o={{COVERAGE_DIR}}/compiled_coverage.out
    go tool cover -html {{COVERAGE_DIR}}/compiled_coverage.out

_tools-linter: (_tools-install "golangci-lint" "curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.62.2")

# Always rebuild the genum, gsort, and gerror executibles from this package direclty for testing.

# Other packages should use the _install-go-pkg script.
_tools-generate: (_install-go-pkg "google.golang.org/protobuf" "cmd/protoc-gen-go")
    go build -o bin/genum genum/cmd/genum/main.go
    go build -o bin/gsort gsort/cmd/gsort/main.go
    go build -o bin/gerror gerror/cmd/gerror/main.go
    go build -o bin/gogenproto gogenproto/cmd/gogenproto/main.go

# Installs the latest monorepo tool.
_tools-monorepo:
    go build -o bin/gomonorepo gomonorepo/cmd/gomonorepo/main.go

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
