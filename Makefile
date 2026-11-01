.DEFAULT_GOAL := build
.PHONY: all build clean test coverage vet fmt fmt-check typecheck lint vuln check run run-http helm-lint e2e e2e-down help

SHELL := /bin/bash

BINARY := mcp
GOFMT := $(shell go env GOROOT)/bin/gofmt
E2E_PROFILE_FILE := .minikube-e2e-profile

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
	go vet -tags e2e ./test/e2e

fmt: ## Format Go source files
	go fmt ./...
	$(GOFMT) -w $$(find "$$(go list -tags e2e -f '{{.Dir}}' ./test/e2e)" -name '*.go')

fmt-check: ## Fail if Go files are not formatted
	@unformatted=$$($(GOFMT) -l $$(find . -name '*.go' -not -path './.git/*')); \
	if [ -n "$$unformatted" ]; then \
		printf '%s\n' "$$unformatted"; \
		echo "run make fmt"; \
		exit 1; \
	fi

typecheck: ## Type-check by compiling
	go build -o bin/$(BINARY) ./cmd

lint: vet ## Run linter
	go tool golangci-lint run
	go tool golangci-lint run --build-tags e2e ./test/e2e

vuln: ## Scan dependencies for known vulnerabilities
	go tool govulncheck ./...

check: fmt-check typecheck vet lint vuln ## Full verification gate
	go test -race ./...

run: build ## Build and run the binary
	./bin/$(BINARY)

run-http: build ## Build and serve HTTP on 127.0.0.1:8080
	./bin/$(BINARY) -http 127.0.0.1:8080

helm-lint: ## Lint and render the Helm chart
	helm lint --strict charts/mcp
	@tmp=$$(mktemp); \
	trap 'rm -f "$$tmp"' EXIT; \
	helm template mcp charts/mcp > "$$tmp"

e2e: ## Run the minikube Helm and in-cluster e2e suite
	@set -euo pipefail; \
	missing=; \
	for tool in docker minikube helm kubectl; do \
		command -v "$$tool" >/dev/null 2>&1 || missing="$$missing $$tool"; \
	done; \
	if [ -n "$$missing" ]; then \
		echo "make e2e requires:$$missing" >&2; \
		exit 1; \
	fi; \
	PROFILE="$${MINIKUBE_PROFILE:-mcp}"; \
	kube() { minikube -p "$$PROFILE" kubectl -- "$$@"; }; \
	minikube delete -p "$$PROFILE"; \
	minikube -p "$$PROFILE" start --driver=docker; \
	printf '%s\n' "$$PROFILE" > "$(E2E_PROFILE_FILE)"; \
	minikube -p "$$PROFILE" addons enable metrics-server || echo "warning: metrics-server may be unavailable; continuing"; \
	docker build -t mcp:dev .; \
	docker build --target e2e -t mcp-e2e:dev .; \
	if ! minikube -p "$$PROFILE" image load mcp:dev mcp-e2e:dev; then \
		echo "minikube image load failed; rebuilding images in the profile Docker environment" >&2; \
		(eval "$$(minikube -p "$$PROFILE" docker-env --shell=bash)"; docker build -t mcp:dev .; docker build --target e2e -t mcp-e2e:dev .); \
	fi; \
	minikube -p "$$PROFILE" update-context; \
	kube create namespace mcp --dry-run=client -o yaml | kube apply -f -; \
	helm --kube-context "$$PROFILE" upgrade --install mcp charts/mcp --namespace mcp --wait --timeout 5m \
		--set image.repository=mcp --set image.tag=dev --set image.pullPolicy=Never \
		--set tests.image.repository=mcp-e2e --set tests.image.tag=dev --set tests.image.pullPolicy=Never; \
	kube rollout status deployment/mcp --namespace mcp --timeout=5m; \
	READY=$$(kube get deployment mcp --namespace mcp -o jsonpath='{.status.readyReplicas}'); \
	[ "$$READY" = "1" ] || { echo "expected 1 ready replica, got $$READY" >&2; exit 1; }; \
	ENDPOINT=""; \
	for attempt in {1..30}; do \
		ENDPOINT=$$(kube get endpoints mcp --namespace mcp -o jsonpath='{.subsets[0].addresses[0].ip}'); \
		[ -n "$$ENDPOINT" ] && break; \
		sleep 1; \
	done; \
	[ -n "$$ENDPOINT" ] || { echo "service has no ready endpoints" >&2; exit 1; }; \
	helm --kube-context "$$PROFILE" test mcp --namespace mcp --logs --timeout 10m; \
	MIN=$$(kube get hpa mcp --namespace mcp -o jsonpath='{.spec.minReplicas}'); \
	MAX=$$(kube get hpa mcp --namespace mcp -o jsonpath='{.spec.maxReplicas}'); \
	[ "$$MIN" = "1" ] && [ "$$MAX" = "3" ] || { echo "unexpected HPA bounds: min=$$MIN max=$$MAX" >&2; exit 1; }; \
	AFFINITY=$$(kube get service mcp --namespace mcp -o jsonpath='{.spec.sessionAffinity}'); \
	[ "$$AFFINITY" = "ClientIP" ] || { echo "unexpected session affinity: $$AFFINITY" >&2; exit 1; }; \
	USER=$$(kube get deployment mcp --namespace mcp -o jsonpath='{.spec.template.spec.containers[0].securityContext.runAsUser}'); \
	[ "$$USER" = "65532" ] || { echo "unexpected container user: $$USER" >&2; exit 1; }; \
	READONLY=$$(kube get deployment mcp --namespace mcp -o jsonpath='{.spec.template.spec.containers[0].securityContext.readOnlyRootFilesystem}'); \
	[ "$$READONLY" = "true" ] || { echo "unexpected readOnlyRootFilesystem: $$READONLY" >&2; exit 1; }; \
	ENDPOINT=$$(kube get endpoints mcp --namespace mcp -o jsonpath='{.subsets[0].addresses[0].ip}'); \
	[ -n "$$ENDPOINT" ] || { echo "service has no ready endpoints" >&2; exit 1; }; \
	RESTARTS=$$(kube get pods --namespace mcp -l app.kubernetes.io/name=mcp,app.kubernetes.io/instance=mcp -o jsonpath='{range .items[*]}{.status.containerStatuses[0].restartCount}{"\n"}{end}'); \
	echo "$$RESTARTS" | grep -qv '[^0]' || { echo "app container restarted: $$RESTARTS" >&2; exit 1; }; \
	echo "e2e passed for minikube profile $$PROFILE"; \
	echo "This command created a clean profile $$PROFILE; tear it down with: make e2e-down"

e2e-down: ## Remove the e2e release and cluster started by make e2e
	@set -euo pipefail; \
	PROFILE=""; \
	if [ -f "$(E2E_PROFILE_FILE)" ]; then \
		PROFILE=$$(cat "$(E2E_PROFILE_FILE)"); \
	elif [ -n "$${MINIKUBE_PROFILE:-}" ]; then \
		PROFILE="$$MINIKUBE_PROFILE"; \
	fi; \
	if [ -n "$$PROFILE" ]; then \
		helm --kube-context "$$PROFILE" uninstall mcp --namespace mcp --ignore-not-found || true; \
		minikube -p "$$PROFILE" kubectl -- delete namespace mcp --ignore-not-found || true; \
	else \
		helm uninstall mcp --namespace mcp --ignore-not-found || true; \
		kubectl delete namespace mcp --ignore-not-found || true; \
	fi; \
	if [ -f "$(E2E_PROFILE_FILE)" ]; then \
		minikube delete -p "$$PROFILE"; \
		rm -f "$(E2E_PROFILE_FILE)"; \
	fi

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
