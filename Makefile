.DEFAULT_GOAL := build
.PHONY: all build clean test vet fmt lint run help

BINARY := mcp

all: build ## Build the binary

clean: ## Remove build artifacts
	rm -rf bin

build: clean ## Build the binary
	go build -o bin/$(BINARY) ./cmd

test: ## Run tests
	go test -v ./...

vet: ## Run go vet
	go vet ./...

fmt: ## Format Go source files
	go fmt ./...

lint: vet ## Run linter
	golangci-lint run

run: build ## Build and run the binary
	./bin/$(BINARY)

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
