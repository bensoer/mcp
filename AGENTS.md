# AGENTS.md — mcp

Personal MCP server exposing standards, practices, and preferences as MCP
resources. See `README.md` for the project overview.

---

## Agent Platform Setup

This repo is agent-agnostic — skills and rules live in `.agents/` and are
discoverable by opencode, Cursor, Claude Code, and other agent platforms.
No platform-specific config files (e.g., `opencode.json`, `.cursorrules`) are
committed.

### Connecting this MCP server to your agent

The compiled binary at `./bin/mcp` is an MCP server that exposes all documents
under `assets/` as resources. To connect it:

**opencode** — add to your project or global `opencode.json`:
```json
"mcp": {
  "ben-mcp": {
    "type": "local",
    "command": ["./bin/mcp"],
    "enabled": true
  }
}
```

**Claude Code** — add to `.claude/mcp.json` or `~/.claude.json`:
```json
"mcpServers": {
  "ben-mcp": {
    "command": "./bin/mcp",
    "args": []
  }
}
```

**Cursor** — add to Cursor MCP settings (`.cursor/mcp.json`):
```json
"mcpServers": {
  "ben-mcp": {
    "command": "./bin/mcp",
    "args": []
  }
}
```

The server reads `assets/` from the process CWD, so configure your agent to
set the working directory to the project root (or set `MCP_ASSET_ROOT`).

### Available agent skills

Skills in `.agents/` provide step-by-step guidance for common tasks in this repo:

| Skill | Location | When to use |
|-------|----------|-------------|
| `create-resource-document` | `.agents/create-resource-document/` | Adding a new standard or workflow document to `assets/` |
| `test-resource-document` | `.agents/test-resource-document/` | Verifying a new or edited resource document is discoverable |
| `server-development` | `.agents/server-development/` | Making Go code changes to the MCP server itself |

### Available agent rules

Rules in `.agents/rules/` capture repo-specific conventions:

| Rule | Covers |
|------|--------|
| `go-conventions.md` | Package structure, interface patterns, error handling for this Go server |
| `resource-conventions.md` | Frontmatter format, URI naming, directory conventions for `assets/` |

---

## Prerequisites

- Go 1.24+
- `make` (recommended) or direct `go` commands
- Python 3 (for the test harness — stdlib only, no deps)

## Development Commands

| Target | Command | Purpose |
|--------|---------|---------|
| `make` / `make build` | `go build -o bin/mcp ./cmd` | Build the binary to `bin/mcp` |
| `make vet` | `go vet ./...` | Static analysis |
| `make test` | `go test -v ./...` | Run tests (none yet — see `docs/plans`) |
| `make fmt` | `go fmt ./...` | Format Go source |
| `make run` | build + `./bin/mcp` | Build then run the server |
| `make clean` | `rm -rf bin` | Remove build artifacts |

Linting uses `go vet` as the primary linter. `golangci-lint` is configured in
the `Makefile` `lint` target but is not required for basic development.

## Common Workflows

### Build and verify

```bash
make build && make vet
```

Run after ANY Go source change. The server won't pick up asset changes until
rebuilt.

### Add a new resource document

1. Use the `create-resource-document` skill for step-by-step guidance.
2. Create the `.md` file under `assets/standards/<category>/` or
   `assets/workflows/<category>/` with correct YAML frontmatter.
3. Rebuild and test:
   ```bash
   make build
   python3 .agents/test-resource-document/scripts/test_resource.py
   ```
4. Verify the URI appears in the resource catalog and content is readable:
   ```bash
   python3 .agents/test-resource-document/scripts/test_resource.py --read standards://your/uri
   ```

### Test resource documents

```bash
# List all registered resources (JSON to stdout)
python3 .agents/test-resource-document/scripts/test_resource.py --build

# Read a specific resource (content to stderr)
python3 .agents/test-resource-document/scripts/test_resource.py --read standards://git/commit-messages
```

The `--build` flag rebuilds the binary first. Without it, uses the existing
`bin/mcp`. The script handles the stdio handshake and server lifecycle
automatically.

### Run the server manually

```bash
make run
```

The server runs over stdio (stdin/stdout JSON-RPC). It reads stdin to EOF then
exits. See the test harness for an example of how to drive it programmatically.

---

## Project Layout

```
cmd/main.go               # Server entrypoint — starts stdio MCP server
internal/
  server.go               # BootstrapServer — walks assets/ and registers MCP resources
  logger/                 # zap-based logger setup (development/production modes)
  models/frontmatter.go   # FrontMatter struct (YAML tags for resource metadata)
  utils/assets.go         # AssetsFinder interface + impl — filesystem walker for assets/
assets/                   # Markdown standard/practice documents (auto-discovered)
  standards/              # Coding standards (git, python, comments, skill-authoring)
  workflows/              # Workflow guides (comments, skill-authoring)
.agents/                  # Agent skills and rules (agent-agnostic)
  create-resource-document/  # Skill: create new resource documents
  test-resource-document/    # Skill: verify resource documents + test harness script
  server-development/        # Skill: work on Go server code
  rules/                     # Repo-specific conventions
docs/plans/               # Implementation plans and design decisions
bin/mcp                   # Compiled binary (gitignored)
```

---

## How Resources Work

Every `.md` file under `assets/` becomes an MCP resource automatically. On
startup, `BootstrapServer()` (`internal/server.go`) walks the `assets/`
directory recursively, reads each file, parses its YAML frontmatter, and
registers it as a resource. **No Go code changes are needed to add a new
resource** — just drop in a markdown file with the right frontmatter.

The asset root is configurable via the `MCP_ASSET_ROOT` env var (defaults to
`assets`, resolved against the process CWD).

### Architecture: Data flow

```
assets/**/*.md
  → AssetsFinder.GetAllAssetPaths()     (internal/utils/assets.go)
  → AssetsFinder.GetAssetContents()     (reads file bytes)
  → frontmatter.Parse()                 (parses YAML frontmatter)
  → models.FrontMatter                  (typed struct)
  → server.AddResource()                (registers MCP resource + handler)
  → MCP resources/list, resources/read  (served to clients)
```

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

---

## Go Code Conventions

- All server code lives under `internal/` — no public API surface.
- Interfaces are defined alongside consumers (see `AssetsFinder` in
  `internal/utils/assets.go` — defined where used, not in a separate package).
- New domains get their own package under `internal/` (e.g., `internal/tools/`
  for future MCP tool endpoints).
- Logging uses `zap.L()` (global) and `zap.S()` (sugared global), initialized
  in `cmd/main.go` via `logger.InitLogger()`.
- See `.agents/rules/go-conventions.md` for detailed package and interface
  conventions.

## Commit Messages

Conventional Commits are required. Follow `standards://git/commit-messages`
(served by this server). In short:

```
<type>[optional scope]: <description>

[optional body explaining what and why]
```

Common types: `feat`, `fix`, `chore`, `refactor`, `docs`, `test`, `ci`, `revert`.

See also `standards://git/commit-staging` for rules on separating commits
(tests, docs, and source changes must be in separate commits).

---

## Before You Run Commands

- Build with `make build` (or `go build -o bin/mcp ./cmd`) before testing
  resources — the server reads `assets/` at startup, so a rebuild + restart is
  required to pick up new or changed documents.
- The binary must run from the project root (or with `MCP_ASSET_ROOT` set) so
  `assets/` is found.
- Run `make vet` after Go changes to catch static analysis issues.