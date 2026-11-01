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
	for tool in docker minikube helm kubectl python3; do \
		command -v "$$tool" >/dev/null 2>&1 || missing="$$missing $$tool"; \
	done; \
	if [ -n "$$missing" ]; then \
		echo "make e2e requires:$$missing" >&2; \
		exit 1; \
	fi; \
	profile_running() { \
		minikube profile list -o json | python3 -c 'import json, sys; data = json.load(sys.stdin); profile = sys.argv[1]; sys.exit(0 if any(p.get("Name") == profile and p.get("Status") in ("Running", "OK") for p in data.get("valid", [])) else 1)' "$$1"; \
	}; \
	profile_exists() { \
		minikube profile list -o json | python3 -c 'import json, sys; data = json.load(sys.stdin); profile = sys.argv[1]; sys.exit(0 if any(p.get("Name") == profile for p in data.get("valid", [])) else 1)' "$$1"; \
	}; \
	if [ -n "$${MINIKUBE_PROFILE:-}" ]; then \
		PROFILE="$$MINIKUBE_PROFILE"; \
	else \
		PROFILE=$$(minikube profile list -o json | python3 -c 'import json, sys; data = json.load(sys.stdin); print(next((p["Name"] for p in data.get("valid", []) if p.get("Status") in ("Running", "OK")), ""))'); \
		PROFILE="$${PROFILE:-mcp}"; \
	fi; \
	STARTED=0; \
	if ! profile_running "$$PROFILE"; then \
		if profile_exists "$$PROFILE"; then \
			minikube -p "$$PROFILE" start; \
		else \
			minikube -p "$$PROFILE" start --driver=docker; \
		fi; \
		STARTED=1; \
	fi; \
	if [ "$$STARTED" -eq 1 ]; then \
		printf '%s\n' "$$PROFILE" > "$(E2E_PROFILE_FILE)"; \
	fi; \
	minikube -p "$$PROFILE" addons enable metrics-server || echo "warning: metrics-server may be unavailable; continuing"; \
	docker build -t mcp:dev .; \
	docker build -f test/e2e/Dockerfile -t mcp-e2e:dev .; \
	if ! minikube -p "$$PROFILE" image load mcp:dev mcp-e2e:dev; then \
		echo "minikube image load failed; rebuilding images in the profile Docker environment" >&2; \
		(eval "$$(minikube -p "$$PROFILE" docker-env --shell=bash)"; docker build -t mcp:dev .; docker build -f test/e2e/Dockerfile -t mcp-e2e:dev .); \
	fi; \
	minikube -p "$$PROFILE" update-context; \
	kubectl --context "$$PROFILE" create namespace mcp --dry-run=client -o yaml | kubectl --context "$$PROFILE" apply -f -; \
	helm upgrade --install mcp charts/mcp --namespace mcp --wait --timeout 5m \
		--set image.tag=dev --set image.pullPolicy=Never \
		--set tests.image.repository=mcp-e2e --set tests.image.tag=dev --set tests.image.pullPolicy=Never; \
	kubectl --context "$$PROFILE" rollout status deployment/mcp --namespace mcp --timeout=5m; \
	READY=$$(kubectl --context "$$PROFILE" get deployment mcp --namespace mcp -o jsonpath='{.status.readyReplicas}'); \
	[ "$$READY" = "1" ] || { echo "expected 1 ready replica, got $$READY" >&2; exit 1; }; \
	helm test mcp --namespace mcp --logs --timeout 10m; \
	MIN=$$(kubectl --context "$$PROFILE" get hpa mcp --namespace mcp -o jsonpath='{.spec.minReplicas}'); \
	MAX=$$(kubectl --context "$$PROFILE" get hpa mcp --namespace mcp -o jsonpath='{.spec.maxReplicas}'); \
	[ "$$MIN" = "1" ] && [ "$$MAX" = "3" ] || { echo "unexpected HPA bounds: min=$$MIN max=$$MAX" >&2; exit 1; }; \
	AFFINITY=$$(kubectl --context "$$PROFILE" get service mcp --namespace mcp -o jsonpath='{.spec.sessionAffinity}'); \
	[ "$$AFFINITY" = "ClientIP" ] || { echo "unexpected session affinity: $$AFFINITY" >&2; exit 1; }; \
	USER=$$(kubectl --context "$$PROFILE" get deployment mcp --namespace mcp -o jsonpath='{.spec.template.spec.containers[0].securityContext.runAsUser}'); \
	[ "$$USER" = "65532" ] || { echo "unexpected container user: $$USER" >&2; exit 1; }; \
	READONLY=$$(kubectl --context "$$PROFILE" get deployment mcp --namespace mcp -o jsonpath='{.spec.template.spec.containers[0].securityContext.readOnlyRootFilesystem}'); \
	[ "$$READONLY" = "true" ] || { echo "unexpected readOnlyRootFilesystem: $$READONLY" >&2; exit 1; }; \
	ENDPOINT=$$(kubectl --context "$$PROFILE" get endpoints mcp --namespace mcp -o jsonpath='{.subsets[0].addresses[0].ip}'); \
	[ -n "$$ENDPOINT" ] || { echo "service has no ready endpoints" >&2; exit 1; }; \
	RESTARTS=$$(kubectl --context "$$PROFILE" get pods --namespace mcp -l app.kubernetes.io/name=mcp,app.kubernetes.io/instance=mcp -o jsonpath='{range .items[*]}{.status.containerStatuses[0].restartCount}{"\n"}{end}'); \
	echo "$$RESTARTS" | grep -qv '[^0]' || { echo "app container restarted: $$RESTARTS" >&2; exit 1; }; \
	echo "e2e passed for minikube profile $$PROFILE"; \
	if [ "$$STARTED" -eq 1 ]; then \
		echo "This command started profile $$PROFILE; tear it down with: make e2e-down"; \
	elif [ -f "$(E2E_PROFILE_FILE)" ]; then \
		echo "This command reused profile $$PROFILE, which an earlier make e2e started; tear it down with: make e2e-down"; \
	else \
		echo "This command reused profile $$PROFILE; make e2e-down will not delete it"; \
	fi

e2e-down: ## Remove the e2e release and cluster started by make e2e
	@set -euo pipefail; \
	helm uninstall mcp --namespace mcp --ignore-not-found || true; \
	kubectl delete namespace mcp --ignore-not-found || true; \
	if [ -f "$(E2E_PROFILE_FILE)" ]; then \
		PROFILE=$$(cat "$(E2E_PROFILE_FILE)"); \
		minikube delete -p "$$PROFILE"; \
		rm -f "$(E2E_PROFILE_FILE)"; \
	fi

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
