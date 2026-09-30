---
uri: workflows://go-new-project
name: Go New Project Workflow
description: "Step-by-step workflow for scaffolding a new Go module with cmd/<binary>/main.go, a composition container, Makefile output in bin/, and golangci-lint."
languages:
    - go
file_types:
    - "*.go"
    - "Makefile"
priority: required
related_resources:
    - standards://go/architecture
    - standards://go/syntax
    - standards://go/logging
    - standards://go/tooling
    - standards://go/cli
    - standards://go/testing
    - workflows://go-cobra-subcommand
---

# Go New Project Workflow
Step-by-step workflow for scaffolding a new Go module.

Use this for a **new** Go module. Existing repositories keep their layout. Follow `standards://go/architecture`, `standards://go/syntax`, `standards://go/logging`, `standards://go/tooling`, `standards://go/cli`, and `standards://go/testing`.

- ALWAYS confirm the module path and binary name before writing files.
- IF a slash command invoked this workflow, stay in Plan mode until the user approves.
- NEVER use `cobra-cli init`. It emits a root `main.go` plus `cmd/root.go`, which is not this layout.

## Step 1: Create the layout

```
cmd/<binary>/main.go       # compose the container and execute
cmd/<binary>/root.go       # Cobra root only, when the program is a CLI
internal/app/container.go  # composition root
internal/<feature>/        # behavior; one owning struct per file
Makefile
.golangci.yml
.pre-commit-config.yaml    # local hooks call make lint, make vet, make test
.gitignore                 # includes bin/
go.mod
```

- ALWAYS run `go mod init <module>`.
- `main` builds the container and calls `Execute`.
- NEVER put business logic in `main`.
- ALWAYS put domain types in `internal/<feature>/`, not under `cmd/`.

## Step 2: Add the container

- ALWAYS add a container for anything beyond a tiny script.
- The container constructs dependencies and is the only registration site.
- NEVER add factory switches, functional options, or application `init()`.
- WHEN several implementations are chosen by name, register them on the container and look them up by key.
- NEVER grow a `switch` per implementation.
- Constructors take explicit parameters, or one typed config struct when there are more than four parameters.
- NEVER write `NewFoo(WithX(...))`.

```go
func NewContainer() *Container {
    return &Container{}
}

func (c *Container) Execute(ctx context.Context) error {
    logger := LoggerFromContext(ctx) // contextual logger, else zap.S()
    _ = logger
    return nil
}
```

- ALWAYS put the base logger on `context.Context` at the command edge.
- NEVER store `context.Context` on a struct.
- `LoggerFromContext` falls back to `zap.S()` when the context has no logger.

## Step 3: Add the CLI root, if this is a CLI

- WHEN the program is a CLI, `cmd/<binary>/root.go` is package `main` and holds the Cobra root.
- ALWAYS fill `Use`, `Short`, `Long`, and `Example`.
- ALWAYS validate flags before I/O.
- Libraries NEVER print or exit.
- Later subcommands follow `workflows://go-cobra-subcommand`.
- ALWAYS place generated command files in `cmd/<binary>/`, not a second `cmd/` package.
- Cobra's own `init()` registration is allowed.
- NEVER add `init()` for tools or services.

## Step 4: Add the Makefile

- `.DEFAULT_GOAL` is `build`.
- Output is `bin/<binary>`.

| Target | Command |
|--------|---------|
| `all` | `build` |
| `build` | `mkdir -p bin` then `go build -o bin/<binary> ./cmd/<binary>` |
| `test` | `go test ./...` |
| `lint` | `golangci-lint run ./...` |
| `vet` | `go vet ./...` |
| `vuln` | `govulncheck ./...` |
| `check` | `build`, `vet`, `go test -race ./...`, `lint`, `vuln`, and `./bin/<binary> --help` when it is a CLI |
| `clean` | `rm -rf bin` |

- ALWAYS list `bin/` in `.gitignore`.
- `make build`, `make test`, and `make check` are the documented commands.

## Step 5: Add pre-commit

- ALWAYS add `.pre-commit-config.yaml` at the repo root.
- Each hook is `language: system`, `always_run: true`, `pass_filenames: false`.
- Each hook `entry` is a Makefile target: `make lint`, `make vet`, `make test`.
- NEVER point hooks at raw `go` commands.
- BEFORE a commit: `pre-commit install` and `pre-commit run --all-files`.

## Step 6: Add lint config

`.golangci.yml` at the repo root (golangci-lint v2):

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

## Checklist

- [ ] `go mod init <module>`
- [ ] `cmd/<binary>/main.go` only wires and executes
- [ ] Container in `internal/app`
- [ ] No functional options, no generics, no application `init()`
- [ ] `bin/` gitignored
- [ ] `.pre-commit-config.yaml` runs `make lint`, `make vet`, and `make test` on every commit
- [ ] `make build`, `make test`, and `make check` are the documented commands
