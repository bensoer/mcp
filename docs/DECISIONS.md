# Decisions — mcp

- **Go 1.27.3**: single binary, fast startup. Pinned to the patched native Go 1.27.3 toolchain.
- **MCP SDK v1.8.0**: HTTP is stateless so the server can serve the `2026-07-28` protocol version. The legacy `initialize` handshake remains supported through `2025-11-25`.
- **Transport**: stdio is the default and `make run` blocks on stdin; `-http` opts into streamable HTTP; HTTP auth is opt-in via `-auth <bolt-path>`, and leaving `-auth` unset keeps the open HTTP server; stdio has no auth; TLS is outside the process; image is distroless nonroot.
- **`assets/` is payload**: standards for other languages are documents this server serves, not this repo's toolchain.
- **Resource reads include frontmatter**: clients get URI metadata with the markdown body.
- **`related_resources` is snake_case**: camelCase YAML keys are dropped.
- **`related_resources` entries are `{uri, when}`**: a bare URI string does not parse. `when` is required and is the only condition for reading that resource.
- **Interfaces live next to their consumer**: see `AssetsFinder` in `internal/utils/assets.go`.
- **Commit checks go through Make**: pre-commit hooks call `make`, they do not call `go` directly.
- **Kubernetes deployment uses `charts/mcp`**: the HPA ranges from 1 to 3 replicas and the Service uses `ClientIP` affinity because sessions are in-memory; probes use `/healthz`, with no in-process TLS; release images publish to GitHub Container Registry and packaged charts attach to GitHub Releases.
- **Chart authentication is opt-in**: `auth.enabled` uses the same `-auth <bolt-path>` path and defaults to off. When enabled, the chart runs one replica because Bolt's file lock is exclusive, stores the database on a claim, and disables the HPA. The default chart remains unauthenticated with the HPA enabled.
