# Decisions — mcp

- **Go 1.26.6**: single binary, fast startup. Pinned so `govulncheck` does not fail on known standard-library issues in older 1.26 patches.
- **stdio transport**: no listening port. `make run` blocks on stdin. There is no `--help` smoke test.
- **`assets/` is payload**: standards for other languages are documents this server serves, not this repo's toolchain.
- **Resource reads include frontmatter**: clients get URI metadata with the markdown body.
- **`related_resources` is snake_case**: camelCase YAML keys are dropped.
- **Interfaces live next to their consumer**: see `AssetsFinder` in `internal/utils/assets.go`.
- **Commit checks go through Make**: pre-commit hooks call `make`, they do not call `go` directly.
