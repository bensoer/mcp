---
name: server-development
description: >-
  Go development workflow for working on the mcp MCP server source code.
  Use when modifying Go files in cmd/, internal/, or when changing the build,
  dependencies, or server behavior. Covers the build/vet cycle, internal
  package conventions, BootstrapServer architecture, and the test harness.
---

# Server Development

Use this skill when making Go code changes to the MCP server — modifying files
in `cmd/`, `internal/`, `go.mod`, `go.sum`, or the `Makefile`.

---

## When to use

- Adding a new internal package (e.g., `internal/tools/` for MCP tools).
- Modifying `BootstrapServer()` or the resource registration logic.
- Changing the logger, frontmatter model, or asset finder.
- Adding tests under `internal/`.
- Updating dependencies (`go get`, `go mod tidy`).

## Step 1 — Understand the architecture

The server is straightforward:

```
cmd/main.go
  → logger.InitLogger()         # init zap global logger
  → internal.BootstrapServer()  # build server with all resources
  → server.Run(stdin_stdout)    # start MCP stdio transport

BootstrapServer() (internal/server.go):
  → utils.NewAssetsFinder()     # point at assets/
  → aff.GetAllAssetPaths()      # walk tree for *.md files
  → for each file:
      → aff.GetAssetContents()  # read bytes
      → frontmatter.Parse()     # extract YAML frontmatter
      → validate uri + name     # fail fast on missing required fields
      → server.AddResource()    # register MCP resource + read handler
```

Key points:
- `BootstrapServer()` is the single place resources are registered.
- The `AssetsFinder` interface (`internal/utils/assets.go`) abstracts file I/O
  — defined alongside the implementation, not in a separate interfaces package.
- Frontmatter parsing uses `github.com/adrg/frontmatter`.
- The MCP SDK is `github.com/modelcontextprotocol/go-sdk`.

## Step 2 — Development cycle

```bash
# Edit Go source...

# Format
make fmt

# Static analysis
make vet

# Build
make build

# If tests exist:
make test
```

Always run `make vet` before committing. The `go vet` tool catches common bugs
(unreachable code, format string mismatches, etc.).

## Step 3 — Internal package conventions

All server code lives under `internal/`. Follow these conventions when adding
new packages:

### Package placement

| Pattern | Example | When |
|---------|---------|------|
| New domain concept | `internal/tools/` | Adding MCP tool endpoints |
| New model | `internal/models/<name>.go` | New data structures with YAML/JSON tags |
| Utility | `internal/utils/<name>.go` | Standalone helpers (file I/O, string manipulation) |

### Interface conventions

- Define interfaces in the same file as the primary consumer, NOT in a
  separate package.
- The `AssetsFinder` interface in `internal/utils/assets.go` is the pattern:
  interface defined at the top, implementation (`AssetsFinderImpl`) below.
- Use `*string` for optional constructor params (see `NewAssetsFinder`).

### Error handling

- Return descriptive errors with `fmt.Errorf("context: %w", err)` to wrap.
- `BootstrapServer()` returns `(*mcp.Server, error)` — errors propagate to
  `main()` where `zap.S().Fatalf()` logs and exits.
- Do not panic in library code — return errors.

### Logger

- The global logger is initialized in `cmd/main.go` via
  `logger.InitLogger(logger.DEVELOPMENT)`.
- Use `zap.S().Infof()`, `zap.S().Errorf()`, etc. (sugared) in server code.
- Use `zap.L()` for structured logging when needed.
- The logger supports two modes: `DEVELOPMENT` (colored, console) and
  `PRODUCTION` (JSON).

## Step 4 — Testing

The test harness at `.agents/test-resource-document/scripts/test_resource.py`
is the primary way to verify the server works end-to-end:

```bash
# From the project root:
python3 .agents/test-resource-document/scripts/test_resource.py --build
```

This rebuilds the binary, starts the server over stdio, sends MCP JSON-RPC
requests, and validates responses. Use `--read <URI>` to verify a specific
resource.

When adding Go tests:
- Test files go next to the code they test (e.g., `internal/utils/assets_test.go`).
- Standard `go test` conventions apply.

## Step 5 — Before committing

- [ ] `make fmt` passes (no formatting changes)
- [ ] `make vet` passes (no static analysis issues)
- [ ] `make build` succeeds
- [ ] Test harness lists all resources without errors (if resource-related changes)
- [ ] Follow `standards://git/commit-staging` — separate source, test, and doc commits
- [ ] Follow `standards://git/commit-messages` — conventional commit format