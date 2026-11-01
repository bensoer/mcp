# Decisions — mcp

- **Go 1.27.3**: single binary, fast startup. Pinned to the patched native Go 1.27.3 toolchain.
- **MCP SDK v1.8.0**: HTTP is stateless so the server can serve the `2026-07-28` protocol version. The legacy `initialize` handshake remains supported through `2025-11-25`; there is no auth.
- **Transport**: stdio is the default and `make run` blocks on stdin; `-http` opts into streamable HTTP; no auth; TLS is outside the process; image is distroless nonroot.
- **`assets/` is payload**: standards for other languages are documents this server serves, not this repo's toolchain.
- **Resource reads include frontmatter**: clients get URI metadata with the markdown body.
- **`related_resources` is snake_case**: camelCase YAML keys are dropped.
- **`related_resources` entries are `{uri, when}`**: a bare URI string does not parse. `when` is required and is the only condition for reading that resource.
- **Interfaces live next to their consumer**: see `AssetsFinder` in `internal/utils/assets.go`.
- **Commit checks go through Make**: pre-commit hooks call `make`, they do not call `go` directly.
- **Kubernetes deployment uses `charts/mcp`**: the HPA ranges from 1 to 3 replicas and the Service uses `ClientIP` affinity because sessions are in-memory; probes use `/healthz`, with no in-process TLS or auth; release images publish to GitHub Container Registry and packaged charts attach to GitHub Releases.
