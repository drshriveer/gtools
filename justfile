PKG_ROOT := `pwd`
INSTALLED_TOOLS := PKG_ROOT / "bin" / ".installed_tools"
export PATH := env_var('PATH') + ':' + PKG_ROOT + '/bin'
export GOBIN := PKG_ROOT + "/bin"
export GOEXPERIMENT := "loopvar"
CURRENT_DIR := invocation_directory_native()

# Runs `go mod tidy` for all modules in the current directory, then sync go workspaces.
tidy: (_invokeMod "go mod tidy -C {}")
    go work sync

# Runs `go test --race ` for all modules in the current directory.
test: (_invokeMod "go test --race {}/...")

# Runs lint and test for all modules in the current directory.
check: lint test

# Runs lint/format for all modules in the current directory.
lint: _tools-linter (_invokeMod "golangci-lint run {}/...")

# Fixes all auto-fixable format and lint errors for all modules in the current directory.
fix: _tools-linter format-md && (_invokeMod "golangci-lint run --fix {}/...")

#    just --fmt --unstable <- disabled because it is causing issues.

# Updates interdependent modules of gtools. TODO: could make this wayyy smarter.
update-interdependencies: \
    (_invokeMod "go get -C {} -u github.com/drshriveer/gtools/gencommon") \
    (_invokeMod "go get -C {} -u github.com/drshriveer/gtools/genum") && tidy

# updates a single package across all go modules.
update-pkg pkgName: && tidy
    _invokeMod "go get -C {} -u  {{ pkgName }}"

# Formats markdown.
format-md: (_install-go-pkg "github.com/moorereason/mdfmt")
    @mdfmt -w -l ./**/*.md

# Runs `go generate` on all modules in the current directory.
generate: _tools-generate (_invokeMod "go generate -C {} ./...")

_tools-linter: (_tools-install "golangci-lint" "curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.54.2")

# Always rebuild the genum, gsort, and gerror executibles from this package direclty for testing.

# Other packages should use the _install-go-pkg script.
_tools-generate: (_install-go-pkg "google.golang.org/protobuf" "cmd/protoc-gen-go")
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
_invokeMod cmd:
    #!/usr/bin/env bash
    set -euo pipefail # makes scripts act like justfiles (https://github.com/casey/just#safer-bash-shebang-recipes)
    if [ "{{ CURRENT_DIR }}" = "{{ PKG_ROOT }}" ]; then
      go list -f '{{{{.Dir}}' -m | xargs -L1 -P 8 -t -I {} {{ cmd }}
    else
      xargs -L1 -P 8 -t -I {} {{ cmd }} <<< "{{ CURRENT_DIR }}"
    fi
