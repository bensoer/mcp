# mcp

Personal [Model Context Protocol](https://modelcontextprotocol.io) server. It serves engineering standards and workflows to AI agents over stdio JSON-RPC, or streamable HTTP with `-http`.

## Prerequisites

- Go 1.27.3+
- GNU Make
- Python 3 (resource harness, stdlib only)
- pre-commit
- Docker
- minikube
- Helm
- kubectl

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

Streamable HTTP negotiates the newest supported MCP spec, `2026-07-28`; older clients still use `initialize` through `2025-11-25`. HTTP is stateless (no session id; `GET /mcp` is 405). No auth; TLS is outside the process.

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

### Run the Kubernetes chart

The Helm chart in `charts/mcp` deploys the same HTTP server image. It enables an
HPA with one to three replicas and uses `ClientIP` session affinity because
MCP sessions are held in memory. Affinity is not a session guarantee: a client
IP change or pod replacement can still interrupt a session.

```bash
make helm-lint  # lint and render the chart without a cluster
make e2e        # recreate minikube and run the in-cluster acceptance suite
make e2e-down   # remove the release and delete only a cluster started by make e2e
```

`make e2e` deletes and recreates the selected minikube profile (`mcp` by
default, or `MINIKUBE_PROFILE` when set) so every run starts from a clean
cluster. It leaves the cluster running and reports how to tear it down. The
same target runs in GitHub Actions CI.

Published GitHub Releases carry both the container image and the packaged
Helm chart `.tgz`.

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
| `make helm-unittest` | install the pinned helm-unittest plugin and validate rendered manifests |
| `make helm-lint` | helm-unittest, strict Helm lint, and a default-value template render |

```bash
pre-commit install
pre-commit run --all-files
```
