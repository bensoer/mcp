# Validation

Run from the repo root. Stop on the first failure.

## Go, tests, or `go.mod`

```bash
make fmt
make fmt-check
make typecheck
make vet
make test
make lint
make vuln
pre-commit run --all-files
```

## Asset documents (`assets/**/*.md`)

```bash
make build
python3 .agents/skills/test-resource-document/scripts/test_resource.py
python3 .agents/skills/test-resource-document/scripts/test_resource.py --read <uri>
pre-commit run --all-files
```

## Makefile targets

| Target | What it runs |
|--------|----------------|
| `make build` | clean `bin/`, then `go build -o bin/mcp ./cmd` |
| `make typecheck` | compile to `bin/mcp` without cleaning first |
| `make test` | `go test -v ./...` |
| `make vet` | `go vet ./...` |
| `make lint` | vet, then `go tool golangci-lint run` |
| `make fmt` | `go fmt ./...` |
| `make fmt-check` | fail if `gofmt -l` finds differences |
| `make vuln` | `go tool govulncheck ./...` |
| `make check` | fmt-check, typecheck, vet, lint, vuln, then `go test -race` |
| `make helm-unittest` | install the pinned helm-unittest plugin and validate rendered manifests |
| `make helm-lint` | helm-unittest, strict Helm lint, and a default-value template render |
| `make e2e` | build/load both images, deploy the chart to minikube, and run the Helm test pod against the deployment |
| `make e2e-down` | uninstall the release, delete the `mcp` namespace, and delete only a minikube profile started by `make e2e` |

`make e2e` requires Docker, minikube, Helm, and kubectl. It is not part of
`make check` or pre-commit, so the default validation gate remains runnable
without Docker. The e2e target does run in CI. The Helm unittest plugin is
installed at a pinned version by `make helm-unittest`, which `make helm-lint`
also runs.

There is no CLI `--help` check. The server reads stdin until EOF.

## Pre-commit

Every Go hook is local, calls `make`, and runs on the whole module (`always_run: true`):

| Hook | Target |
|------|--------|
| `fmt` | `make fmt-check` |
| `typecheck` | `make typecheck` |
| `vet` | `make vet` |
| `lint` | `make lint` |
| `test` | `make test` |
| `vuln` | `make vuln` |

Hygiene hooks also run: merge conflicts, newline, trailing whitespace, YAML, large files.

Install once: `pre-commit install`.
