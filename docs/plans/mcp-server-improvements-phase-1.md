# MCP Server Phase 1 Improvement Plan

## Scope

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

### Out of Scope (Future Phases)

- **Phase 2:** Resource annotations (priority, languages, file_types), graceful
  per-file error handling, resource templates for dynamic URI patterns
- **Phase 3:** MCP Tools (`validate_code`, `search_standards`), MCP Prompts
  (`summarize_standards`, `explain_violation`)
- **Phase 4:** Expanded Makefile, tests, README updates
