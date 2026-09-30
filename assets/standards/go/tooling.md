---
uri: standards://go/tooling
name: Go (Golang) Tooling Standards
description: "Makefile, bin output, pre-commit, golangci-lint, and govulncheck standards for Go projects."
languages:
    - go
file_types:
    - "*.go"
    - "Makefile"
    - ".pre-commit-config.yaml"
    - ".golangci.yml"
priority: required
related_resources:
    - standards://go/architecture
    - standards://go/testing
    - workflows://go-new-project
---

# Go (Golang) Tooling Standards
Makefile, bin output, pre-commit, golangci-lint, and govulncheck standards for Go projects.

- ALWAYS follow the repository when it already standardizes on different tools. Repository rules win.
- ALWAYS apply these rules to new Go projects. NEVER restyle an existing repository's toolchain unless asked.

## Makefile

- ALWAYS include a `Makefile` at the repository root of every new Go project.
- ALWAYS route build, test, lint, and clean through `make`.
- NEVER tell developers to run raw `go` commands as the primary workflow.
- ALWAYS send build output to `bin/`. Example: `go build -o bin/<name> ./cmd/<name>`.
- ALWAYS list `bin/` in `.gitignore`.

Minimum targets:

| Target | Purpose |
|--------|---------|
| `all` | Default. Calls `build`. |
| `build` | Compile to `bin/`. |
| `test` | `go test ./...` |
| `vet` | `go vet ./...` |
| `lint` | `golangci-lint` |
| `clean` | Remove `bin/` and other build artifacts |
| `check` | build, `go test -race ./...`, `vet`, lint, vuln, and CLI `--help` smoke when the project is a CLI |

- ALWAYS include `fmt`, `tidy`, and `vuln` (`govulncheck ./...`) on a real project.
- `.DEFAULT_GOAL` is `build`.

## Pre-commit

- ALWAYS add a root `.pre-commit-config.yaml` to every new Go project.
- ALWAYS make hooks local. Hooks call Makefile targets.
- NEVER call `go`, `golangci-lint`, or `govulncheck` directly from a hook.
- ALWAYS set `language: system`, `always_run: true`, and `pass_filenames: false` on each hook, so the full suite runs on every commit.
- NEVER lint only the staged filenames.

Required hooks:

| Hook | Entry |
|------|--------|
| `lint` | `make lint` |
| `vet` | `make vet` |
| `test` | `make test` |

- `make vuln`, the race test, the build, and CLI `--help` stay on `make check`. They are NOT separate pre-commit hooks.
- BEFORE creating a commit, install the hooks if needed and run the full suite. Both MUST pass:

```bash
pre-commit install
pre-commit run --all-files
```

## Lint and security

- ALWAYS lint with `golangci-lint` via `.golangci.yml`.
- ALWAYS run lint with `make lint`.
- ALWAYS run `govulncheck` via a `vuln` target. `vuln` is also part of `check`.

`.golangci.yml` for a new project (golangci-lint v2):

```yaml
version: "2"

run:
  timeout: 5m
  tests: true

linters:
  enable:
    - govet
    - errcheck
    - staticcheck
    - ineffassign
    - unused
```
