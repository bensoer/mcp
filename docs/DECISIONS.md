# Decisions — mcp

- **Go 1.26.6**: single binary, fast startup. Pinned so `govulncheck` does not fail on known standard-library issues in older 1.26 patches.
- **MCP SDK v1.8.0**: HTTP is stateless so the server can serve the `2026-07-28` protocol version. The legacy `initialize` handshake remains supported through `2025-11-25`; there is no auth.
- **Transport**: stdio is the default and `make run` blocks on stdin; `-http` opts into streamable HTTP; no auth; TLS is outside the process; image is distroless nonroot.
- **`assets/` is payload**: standards for other languages are documents this server serves, not this repo's toolchain.
- **Resource reads include frontmatter**: clients get URI metadata with the markdown body.
- **`related_resources` is snake_case**: camelCase YAML keys are dropped.
- **Interfaces live next to their consumer**: see `AssetsFinder` in `internal/utils/assets.go`.
- **Commit checks go through Make**: pre-commit hooks call `make`, they do not call `go` directly.
