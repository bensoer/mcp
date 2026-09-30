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
