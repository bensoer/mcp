# MCP Server Testing Infrastructure (IMPLEMENTED)

Adds unit tests, integration tests, and CI to the MCP server. Resolves the
"Add Tests" bullet in Phase 4 of `mcp-server-improvements-phase-1.md`.

---

## Motivation

The MCP server had zero test coverage. `BootstrapServer()` created its own
`AssetsFinderImpl` internally (hardcoding the `assets/` path), making it
untestable without real asset files on disk. There was also no CI to enforce
tests on PRs.

## Approach

### 1. Refactor `BootstrapServer` for Dependency Injection

**File:** `internal/server.go`

Changed `BootstrapServer()` → `BootstrapServer(finder utils.AssetsFinder)`.
The function now receives the `AssetsFinder` interface instead of constructing
one internally. The existing `AssetsFinder` interface was already perfectly
suited for this — no new abstraction needed.

The `os.Getenv("MCP_ASSET_ROOT")` logic moved to `cmd/main.go`, which is
the only caller that should know about environment configuration.

**Rationale:** An interface parameter is simpler than adding a new public
function variant. It's a one-line signature change. Tests inject a mock;
production injects the real `AssetsFinderImpl`.

### 2. Unit Tests for `AssetsFinderImpl`

**File:** `internal/utils/assets_test.go` (new)

7 tests covering the real filesystem implementation:

| Test | What's verified |
|------|----------------|
| `GetAssetFolderRoot` | Returns the configured root path |
| `GetAssetPath` | Joins root + relative path correctly |
| `GetAssetContents` | Reads file contents (via `t.TempDir()`) |
| `GetAssetContents_NotFound` | Error on nonexistent file |
| `GetAllAssetPaths` | Recursively discovers files in subdirs |
| `GetAllAssetPaths_EmptyDir` | Returns empty slice for empty dir |
| `GetAllAssetPaths_NonexistentRoot` | Error for nonexistent root |

### 3. Unit + Integration Tests for `BootstrapServer`

**File:** `internal/server_test.go` (new)

Uses a `mockAssetsFinder` that satisfies the `AssetsFinder` interface with
in-memory strings — no disk I/O.

**Unit tests (6):**

| Test | What's verified |
|------|----------------|
| `ValidAssets` | Server starts with valid frontmatter |
| `MissingURI` | Error when `uri` frontmatter field is absent |
| `MissingName` | Error when `name` frontmatter field is absent |
| `InvalidYAML` | Error on malformed YAML frontmatter |
| `GetAllAssetPathsError` | Error propagated from finder |
| `GetAssetContentsError` | Error propagated from finder |

**Integration tests (2):**

Uses `mcp.NewInMemoryTransports()` to test the full MCP protocol stack:

| Test | What's verified |
|------|----------------|
| `Integration` | Bootstrap with mock → connect MCP client → list resources → read resource → verify content and MIME type match |
| `Integration_UnknownResource` | Reading a non-existent URI returns an error |

### 4. GitHub Actions CI

**File:** `.github/workflows/ci.yml` (new)

Runs on every PR and push to `main`/`master`:

```yaml
jobs:
  test:
    steps:
      - checkout
      - setup-go (1.25, with caching)
      - make build
      - make vet
      - make test
```

---

## Results

```
15 tests, 0 failures
  internal:         8 tests (6 unit + 2 integration)
  internal/utils:   7 tests
  internal/logger:  no test files
  internal/models:  no test files
  cmd:              no test files
```

Existing Python test harness (`test_resource.py --build`) continues to pass
with all 11 resources discoverable.