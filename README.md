# mcp

Personal [Model Context Protocol](https://modelcontextprotocol.io) server. It serves engineering standards and workflows to AI agents over stdio JSON-RPC.

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
