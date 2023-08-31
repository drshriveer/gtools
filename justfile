GO_LINT_VERSION := '1.53.3'
PKG := `go list -m`
PKG_ROOT := `pwd`
export PATH := env_var('PATH') + ':' + PKG_ROOT + '/bin'
export GOBIN := PKG_ROOT + "/bin"

test:
    @go test --race ./...

check: check-format lint

# Format
format: _tools-format
    @just --fmt --unstable
    @goimports -l -w -local {{ PKG }} {{ PKG_ROOT }}/pkg

check-format: _tools-format
    @just --check --fmt --unstable
    @goimports -l -d -local {{ PKG }} {{ PKG_ROOT }}/pkg
    @echo "Format Check Successful!"

_tools-format:
    @go install golang.org/x/tools/cmd/goimports@latest

# Lint
lint: _tools-linter
    @golangci-lint run ./...

lint-fix: _tools-linter
    @golangci-lint run -fix ./...

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
    @go install github.com/drshriveer/gcommon/pkg/genum/genum
