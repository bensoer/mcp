# AGENTS.md — mcp

Personal MCP server exposing standards, practices, and preferences as MCP
resources. See `README.md` for the project overview.

---

## Prerequisites

- Go 1.24+
- `make` (recommended) or direct `go` commands

## Development Commands

| Target | Command | Purpose |
|--------|---------|---------|
| `make` / `make build` | `go build -o bin/mcp ./cmd` | Build the binary to `bin/mcp` |
| `make vet` | `go vet ./...` | Static analysis |
| `make test` | `go test -v ./...` | Run tests (none yet — see `docs/plans`) |
| `make fmt` | `go fmt ./...` | Format Go source |
| `make run` | build + `./bin/mcp` | Build then run the server |
| `make clean` | `rm -rf bin` | Remove build artifacts |

Linting uses `go vet` as the linter target (see `Makefile`). `golangci-lint`
is configured in the `Makefile` `lint` target but is not required for basic
development.

## Project Layout

```
cmd/main.go               # Server entrypoint — starts stdio MCP server
internal/
  server.go               # BootstrapServer — walks assets/ and registers MCP resources
  logger/                 # zap-based logger setup
  models/frontmatter.go   # FrontMatter struct (YAML tags for resource metadata)
  utils/assets.go         # AssetsFinder — filesystem walker for the assets/ dir
assets/                   # Markdown standard/practice documents (auto-discovered)
  standards/              # Coding standards (git, python, comments)
  workflows/              # Workflow guides (comments, refactoring)
docs/plans/               # Multi-phase improvement plan
.agents/                  # Agent skills (how-to guides for working in this repo)
bin/mcp                   # Compiled binary
```

## How Resources Work

Every `.md` file under `assets/` becomes an MCP resource automatically. On
startup, `BootstrapServer()` (`internal/server.go`) walks the `assets/`
directory recursively, reads each file, parses its YAML frontmatter, and
registers it as a resource. **No Go code changes are needed to add a new
resource** — just drop in a markdown file with the right frontmatter.

The asset root is configurable via the `MCP_ASSET_ROOT` env var (defaults to
`assets`, resolved against the process CWD).

### Required frontmatter fields

Every resource document must have a YAML frontmatter block (`--- ... ---`) with
at minimum:

| Field | Type | Notes |
|-------|------|-------|
| `uri` | string | Unique MCP resource URI, e.g. `standards://git/commit-messages`. **Required** — server fails to start if empty. |
| `name` | string | Display name. **Required** — server fails to start if empty. |
| `description` | string | Shown in resource listings. |
| `languages` | `[]string` | e.g. `["all"]`. Documentation-only. |
| `file_types` | `[]string` | e.g. `["*.*"]`. Documentation-only. |
| `priority` | string | e.g. `required`. Documentation-only. |
| `related_resources` | `[]string` | URIs of related standards. **Must be snake_case** — the struct tag is `yaml:"related_resources"`; camelCase keys (`relatedResources`) are silently dropped. |

Follow the pattern in existing files like `assets/standards/git/commit-messages.md`.
See the `create-resource-document` and `test-resource-document` skills in
`.agents/` for step-by-step instructions.

## Commit Messages

Conventional Commits are required. See `standards://git/commit-messages`
(served by this same server) and `ref/commands/commit-message.md`.

## Before You Run Commands

- Build with `make build` (or `go build -o bin/mcp ./cmd`) before testing
  resources — the server reads `assets/` at startup, so a rebuild + restart is
  required to pick up new or changed documents.
- The binary must run from the project root (or with `MCP_ASSET_ROOT` set) so
  `assets/` is found.
