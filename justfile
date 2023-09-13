GO_LINT_VERSION := '1.53.3'
PKG_ROOT := `pwd`
MODS := `go list -f '{{.Dir}}' -m`
export PATH := env_var('PATH') + ':' + PKG_ROOT + '/bin'
export GOBIN := PKG_ROOT + "/bin"

tidy target='all':
    @just _invokeMod "go mod tidy -C {}" "{{ target }}"

test target='all':
    @just _invokeMod "go test --race {}" "{{ target }}"

check: lint test

# Lint
lint target='all': _tools-linter
    @just _invokeMod "golangci-lint run {}" "{{ target }}"

fix target='all': _tools-linter
    @just _invokeMod "golangci-lint run --fix {}" "{{ target }}"

_tools-linter:
    #!/usr/bin/env sh
    if command -v golangci-lint && golangci-lint --version | grep -q '{{ GO_LINT_VERSION }}'; then
      echo 'golangci-lint v{{ GO_LINT_VERSION }} already installed!'
    else
      echo "installing golangci-lint at version v{{ GO_LINT_VERSION }}"
      if test -e ./bin/golangci-lint; then
        rm ./bin/golangci-lint
      fi
      curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v{{ GO_LINT_VERSION }}
    fi

# Generate
generate: _tools-generate
    @go generate ./...

_tools-generate:
    @go install github.com/drshriveer/gtools/pkg/genum/genum

# invokeMod invokes a command on a module target or all the input command must include

# a the placeholder `{}` which is the path to the correct module.
_invokeMod cmd target='all':
    #!/usr/bin/env sh
    if [ "{{ target }}" = "all" ]; then
      xargs -L1 -t -I {} {{ cmd }} <<< "{{ MODS }}"
     else
      xargs -L1 -t -I {} {{ cmd }} <<< "{{ target }}"
    fi
