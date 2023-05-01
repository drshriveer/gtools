PKG := $(shell go list -m)

build: tools.build
	@go build pkg

test: build
	@go test --race ./pkg/...

format:
	@gofmt -l -w -s ./pkg/...

lint: tools.lint
	@#golangci-lint run ./pkg/...

check.format:
	@gofmt -l -d ./pkg/...

generate: tools.genum
	@go generate ./...

tools.build:
	@go mod download

tools.lint:
	@echo "currently disabled, should download golangci-lint"

tools.genum:
	go install github.com/drshriveer/gcommon/cmd/genum