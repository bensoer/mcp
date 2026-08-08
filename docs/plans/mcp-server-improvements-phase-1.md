# MCP Server Multi-Phase Improvement Plan

This document covers the full multi-phase plan for improving the MCP server code.
Phase 1 has been implemented. Phases 2–4 are planned for future sessions.

---

## Phase 1 — Triage & Cleanup (IMPLEMENTED)

Phase 1 focuses on triage, correctness, and cleanup — low-risk changes that remove
dead code, fix bugs, and make the server more robust without adding new features.

---

### 1. Fix Frontmatter Key Inconsistency

**Files:** `assets/standards/comments/comments.md`, `assets/standards/git/commit-staging.md`

The `FrontMatter` struct in `internal/models/frontmatter.go` tags the
`RelatedResources` field as `yaml:"related_resources"` (snake_case). However,
`comments.md` and `commit-staging.md` use `relatedResources:` (camelCase) in
their frontmatter. Since YAML is case-sensitive when matching struct field tags,
`relatedResources` will not map to the `RelatedResources` field — those two
files silently lose their related resources data.

**Fix:** Normalize `relatedResources:` → `related_resources:` in both files to
match the struct tag and all other asset files.

---

### 2. Add Validation That Throws When URI or Name Are Empty

**File:** `internal/server.go`

Currently, if a markdown file in `assets/` lacks frontmatter (or has an empty
`uri` field), the code registers a resource with `URI: ""`. The MCP protocol
requires URIs to be valid — an empty URI is invalid and may cause client errors.

**Fix:** Before registering each resource, check that `metaData.URI` and
`metaData.Name` are non-empty. If either is empty, return a descriptive error
identifying the offending file, causing the server to fail fast at startup.

---

### 3. Make Asset Path Configurable

**Files:** `internal/server.go`, `cmd/main.go`

The asset folder root is hardcoded as `assetFolderRoot := "assets"`, a relative
path resolved against the process CWD. This only works if the binary is run from
the project root.

**Fix:** Read the asset folder root from the `MCP_ASSET_ROOT` environment
variable at startup, with a fallback to `"assets"`.

---

### 4. Remove Dead Code

**Files to delete or clean:**

- Delete `internal/resources/python.go` — the entire `PythonResourceService`
  interface and implementation is dead code. The old approach (static per-domain
  services with hardcoded paths) was replaced by the dynamic asset-walking
  approach in `server.go`, but this file was never deleted.
- Remove commented-out blocks in `internal/models/frontmatter.go` (lines 3-24) —
  example frontmatter data and example Go source that serve no documentation
  purpose.
- Remove commented-out blocks in `internal/server.go` (lines 64-99) — old
  commented-out implementation of per-domain resource registration using
  `PythonResourceService`.

---

### 5. Fix Typo: `BoostrapServer` → `BootstrapServer`

**Files:** `internal/server.go`, `cmd/main.go`

The function `BoostrapServer` (missing an "r") was renamed to `BootstrapServer`
in both the definition and the call site in `cmd/main.go`.

---

### Execution Order

1. Normalize frontmatter keys in `comments.md` and `commit-staging.md`
2. Rewrite `server.go` with validation, env var support, and typo fix
3. Strip dead code from `frontmatter.go`
4. Delete `internal/resources/python.go`
5. Update `cmd/main.go` call site
6. Run `go mod tidy` to clean up dependency declarations
7. Verify build and vet pass

---

## Phase 2 — Robustness (Planned)

Medium-risk changes that improve error handling and resource metadata.

### 1. Graceful Error Handling for Individual Asset Failures

**File:** `internal/server.go`

Currently, if any single asset file has malformed frontmatter or can't be read,
the entire server fails to start. Instead, log the error for that specific file
and continue registering the rest.

### 2. Surface Frontmatter Metadata as Resource Annotations

**File:** `internal/server.go`

The `FrontMatter` struct has `Languages`, `FileTypes`, `Priority`, and
`RelatedResources` fields, but only `URI`, `Name`, `Description`, and `MIMEType`
are passed to the `mcp.Resource` registration. The MCP protocol supports
resource annotations (priority, audience/role, etc.). These should be surfaced in
the resource definition so clients can filter and present them effectively.

### 3. Convert to Resource Templates (Optional)

Replace the eager-read loop with `server.AddResourceTemplate` using URIs like
`standards://{category}/{topic}`. Defers file I/O to request time, scales better.

---

## Phase 3 — New Capabilities (Planned)

Higher-effort changes that fundamentally transform the server's effectiveness.

### 1. Add MCP Tools

The README's core idea: tool endpoints that allow LLMs to validate code against
standards. A new service would hold the logic, registered in `BootstrapServer`.

- **`internal/tools/code_validator.go`** — implements
  `NewCodeValidatorTool(assetsFinder utils.AssetsFinder, logger *zap.SugaredLogger)`
- Tool `validate_code`: accepts `code` (string) and `language` (string), reads the
  relevant standards from assets, returns text describing violations
- Tool `search_standards`: accept `language`, `file_type`, `priority` params,
  return matching standards
- Tool `get_related_standards`: resolve `related_resources` links from
  frontmatter given a URI

### 2. Add MCP Prompts

Register prompt templates in `BootstrapServer`:

- `summarize_standards` — "Summarize these coding standards for language: X"
- `explain_violation` — "Explain why this code violates the standard at: <URI>"
- `generate_compliant_code` — "Write {language} code following these standards"

---

## Phase 4 — Developer Experience (Planned)

### 1. Expand Makefile

Add standard targets per the repo's Go practices rules:

| Target | Command |
|--------|---------|
| `all` | calls `build` |
| `build` | `go build -o bin/mcp ./cmd/main.go` |
| `test` | `go test ./...` |
| `lint` | `go vet ./...` |
| `fmt` | `go fmt ./...` |
| `tidy` | `go mod tidy` |
| `run` | build and execute |
| `clean` | remove `bin/` |

### 2. Add Tests

- `internal/utils/assets_test.go` — test `GetAllAssetPaths` and `GetAssetContents`
- `internal/server_test.go` — test that `BootstrapServer` registers expected
  resources and fails on invalid frontmatter

### 3. Update README

Document actual resources, URI scheme, available tools/prompts, and the
`MCP_ASSET_ROOT` env var.
