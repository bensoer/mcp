.DEFAULT_GOAL := build
.PHONY: all build clean test coverage vet fmt fmt-check typecheck lint vuln check run help

SHELL := /bin/bash

BINARY := mcp
GOFMT := $(shell go env GOROOT)/bin/gofmt

all: build ## Build the binary

clean: ## Remove build artifacts
	rm -rf bin coverage.out

build: clean ## Build the binary
	go build -o bin/$(BINARY) ./cmd

test: ## Run tests
	go test -v ./...

coverage: ## Run tests with coverage report
	go test ./internal/... ./cmd/... -coverprofile=coverage.out
	go tool cover -func=coverage.out

vet: ## Run go vet
	go vet ./...

fmt: ## Format Go source files
	go fmt ./...

fmt-check: ## Fail if Go files are not formatted
	@unformatted=$$($(GOFMT) -l $$(go list -f '{{.Dir}}' ./...)); \
	if [ -n "$$unformatted" ]; then \
		printf '%s\n' "$$unformatted"; \
		echo "run make fmt"; \
		exit 1; \
	fi

typecheck: ## Type-check by compiling
	go build -o bin/$(BINARY) ./cmd

lint: vet ## Run linter
	go tool golangci-lint run

vuln: ## Scan dependencies for known vulnerabilities
	go tool govulncheck ./...

check: fmt-check typecheck vet lint vuln ## Full verification gate
	go test -race ./...

run: build ## Build and run the binary
	./bin/$(BINARY)

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
