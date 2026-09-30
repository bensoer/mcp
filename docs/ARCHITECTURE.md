# Architecture — mcp

```mermaid
graph TD
    Client[MCP client] <-->|stdio JSON-RPC| Main[cmd/main.go]
    Main --> Bootstrap[internal.BootstrapServer]
    Bootstrap --> Finder[AssetsFinder]
    Bootstrap --> Tools[internal/tools]
    Finder -->|walk *.md| Assets[(assets/)]
    Bootstrap -->|parse YAML| FM[FrontMatter]
    Tools --> Finder
```

| Component | Path | Role |
|-----------|------|------|
| Entrypoint | `cmd/main.go` | Logger, asset root, stdio server |
| Bootstrap | `internal/server.go` | Register resources and tools |
| Assets | `internal/utils/assets.go` | Walk and read `assets/` or `MCP_ASSET_ROOT` |
| Models | `internal/models/` | YAML frontmatter |
| Tools | `internal/tools/` | discover/search standards and workflows |
| Payload | `assets/` | Markdown registered as resources |

Resource reads return the full file, including frontmatter. New documents do not need Go changes. The server reads `assets/` once at startup.
