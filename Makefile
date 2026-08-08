.DEFAULT_GOAL := build
.PHONY: all build clean test coverage vet fmt lint run help

SHELL := /bin/bash

BINARY := mcp

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

lint: vet ## Run linter
	go tool golangci-lint run

run: build ## Build and run the binary
	./bin/$(BINARY)

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
