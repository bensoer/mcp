# Go Conventions — mcp Server

Repo-specific Go conventions for the mcp MCP server. These apply to Go source
in `cmd/` and `internal/` and supplement general Go best practices with
conventions established in this codebase.

---

## Package Structure

- All server code lives under `internal/` — no exported API surface.
- `cmd/main.go` is the sole entrypoint. Its only job: init logger, bootstrap,
  run.
- New capability domains (tools, prompts, validation) get their own package
  under `internal/` rather than being added to `server.go`.

## Interface Conventions

- Define interfaces in the same file as the primary consumer, NOT in a
  separate `interfaces/` or `ports/` package.
- The pattern is established by `internal/utils/assets.go`:
  ```go
  type AssetsFinder interface {
      GetAssetPath(assetName string) string
      GetAssetFolderRoot() string
      GetAssetContents(assetPathInsideFolder string) ([]byte, error)
      GetAllAssetPaths() ([]string, error)
  }

  type AssetsFinderImpl struct { ... }

  func NewAssetsFinder(assetFolderRoot *string) AssetsFinder { ... }
  ```
- Interface names describe the role: `AssetsFinder`, not `IAssetsFinder`.
- Implementation structs are unexported (`assetsFinderImpl` is acceptable but
  `AssetsFinderImpl` is the existing convention in this repo).

## Error Handling

- Return errors; never panic in library code.
- Wrap errors with context: `fmt.Errorf("failed to read asset %s: %w", path, err)`.
- `main()` is the only place `Fatal`/`Fatalf` is called (via `zap.S().Fatalf`).

## Logging

- Logger is initialized once in `cmd/main.go` via `logger.InitLogger()`.
- Use sugared logger: `zap.S().Infof(...)`, `zap.S().Errorf(...)`.
- `DEVELOPMENT` mode for local work; `PRODUCTION` mode for deployed use.
- Loggers in packages follow the pattern of accepting an already-configured
  logger (or using the global `zap.L()` / `zap.S()`).

## Types

- Data structures with YAML/JSON serialization live in `internal/models/`.
- Use struct tags explicitly: `yaml:"field_name"` for YAML, `json:"field_name"`
  for JSON.
- YAML tags use snake_case to match frontmatter conventions.

## Constructor Pattern

- Use `New<Type>` functions that accept dependencies and return the interface
  type (not the concrete implementation):
  ```go
  func NewAssetsFinder(assetFolderRoot *string) AssetsFinder {
      return &AssetsFinderImpl{assetFolderRoot: assetFolderRoot}
  }
  ```

## Dependencies

- Keep dependencies minimal. Current core deps:
  - `github.com/modelcontextprotocol/go-sdk` — MCP protocol
  - `github.com/adrg/frontmatter` — YAML frontmatter parsing
  - `go.uber.org/zap` — structured logging
- Run `go mod tidy` after adding or removing imports.