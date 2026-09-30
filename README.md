# mcp

Personal [Model Context Protocol](https://modelcontextprotocol.io) server. It serves engineering standards and workflows to AI agents over stdio JSON-RPC, or streamable HTTP with `-http`.

## Prerequisites

- Go 1.26.6+
- GNU Make
- Python 3 (resource harness, stdlib only)
- pre-commit

## Quickstart

```bash
make build
python3 .agents/skills/test-resource-document/scripts/test_resource.py
```

The binary is `./bin/mcp`. Run it from the repository root, or set `MCP_ASSET_ROOT`.

## Connect a client

```json
{
  "mcpServers": {
    "ben-mcp": {
      "command": "/absolute/path/to/mcp/bin/mcp",
      "args": []
    }
  }
}
```

### Run over HTTP

Streamable HTTP (MCP spec `2025-11-25`). No auth; TLS is outside the process.

```bash
make run-http            # builds, then serves on http://127.0.0.1:8080/mcp
docker run --rm -p 8080:8080 mcp
```

Connect a client to `http://127.0.0.1:8080/mcp`:

```bash
pi mcp add ben-mcp --url http://127.0.0.1:8080/mcp
claude mcp add --transport http ben-mcp http://127.0.0.1:8080/mcp
```

Pi and Claude (and any MCP host that takes a URL):

```json
{ "mcpServers": { "ben-mcp": { "url": "http://127.0.0.1:8080/mcp" } } }
```

Codex uses TOML:

```toml
[mcp_servers.ben-mcp]
url = "http://127.0.0.1:8080/mcp"
```

Pi project prompts in `.agents/commands/` are loaded through `.pi/settings.json`.

## Layout

```
cmd/main.go          entrypoint
internal/            server, tools, asset walker
assets/              markdown payload, not repo instructions
.agents/             rules, skills, and commands for this repo
docs/                context, architecture, decisions, plans
```

Agents changing this repository should start at [AGENTS.md](AGENTS.md).

## Checks

| Target | Purpose |
|--------|---------|
| `make build` | Compile `bin/mcp` |
| `make test` | `go test -v ./...` |
| `make vet` | `go vet` |
| `make lint` | vet + golangci-lint |
| `make vuln` | govulncheck |
| `make fmt` | gofmt |
| `make check` | fmt, typecheck, vet, lint, vuln, race tests |

```bash
pre-commit install
pre-commit run --all-files
```
