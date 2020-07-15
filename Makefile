BIN_DIR         ?= $(shell go env GOPATH)/bin

default: test

deps: ## download go modules
	go mod download

fmt: lint/install # ensure consistent code style
	gofmt -s -w .
	golangci-lint run --fix > /dev/null 2>&1 || true

lint/install:
	@if ! golangci-lint --version > /dev/null 2>&1; then \
	  echo "Installing golangci-lint"; \
	  curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(BIN_DIR) v1.28.3; \
	fi

lint: lint/install ## run golangci-lint
	@if ! golangci-lint run; then \
  		echo "\033[0;33mdetected fmt problems: run \`\033[0;32mmake fmt\033[0m\`"; \
  		exit 1; \
  	fi

test: lint ## run go tests
	go vet ./...
	go test ./... -race

help: ## displays this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_\/-]+:.*?## / {printf "\033[34m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | \
		sort | \
		grep -v '#'
