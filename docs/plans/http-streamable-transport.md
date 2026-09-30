# HTTP Streamable Transport

Implement this file. It is the task, not a historical note. `AGENTS.md` says `docs/plans/` is not the source of truth. That does not apply here.

Do not redesign. Do not add packages. If a step is blocked, stop and report. Do not invent a workaround.

## Read first

| Path | Why |
|------|-----|
| `cmd/main.go` | Only Go file to change. Stdio today. |
| `internal/server.go` | `BootstrapServer(finder) (*mcp.Server, error)`. Do not edit. |
| `Makefile` | Add `run-http` only. |
| `.agents/rules/validation.md` | Checks. |
| `.agents/skills/server-development/SKILL.md` | Build cycle only. Its "stdio only" diagram is stale for this task. Do not update the skill. Do not follow it for transport. |

SDK examples are not in this repo. They are in the module cache:

`$(go env GOMODCACHE)/github.com/modelcontextprotocol/go-sdk@v1.6.1/examples/server/everything/main.go`

`$(go env GOMODCACHE)/github.com/modelcontextprotocol/go-sdk@v1.6.1/examples/http/main.go`

`$(go env GOMODCACHE)/github.com/modelcontextprotocol/go-sdk@v1.6.1/mcp/streamable.go` (`StreamableHTTPOptions`)

Pinned module: `github.com/modelcontextprotocol/go-sdk` v1.6.1. Do not bump it. `go.mod` says `go 1.26.6`. Local toolchain may be newer. Do not change `go.mod` or `go.sum`.

## What to build

Same process, two transports.

- No args: stdio, as today. Python harness (`test_resource.py`) runs `bin/mcp` with no args and talks stdin/stdout. That must keep working.
- `-http <host:port>`: streamable HTTP from the SDK. Not a hand-rolled JSON-RPC server. Not SSE. Not WebSocket.

Auth is optional in the MCP spec (`2025-11-25`). This server has none. Pi, Claude, Cursor, and Codex all connect to a streamable HTTP URL with no token. Do not add bearer auth, API keys, or OAuth.

## Locked behavior

- Flag name is `-http`. Value is `host:port`. Example: `-http 0.0.0.0:8080`. Whitespace-only is a fatal error in `main`. Unset means stdio.
- No Cobra, Viper, or mode env var. This repo is not a Cobra CLI. SDK examples use `flag`.
- Declare the flag in `main` with `flag.NewFlagSet(os.Args[0], flag.ContinueOnError)`, not the global `flag.CommandLine`. `go test` parses the global set. A global `-http` flag breaks `go test`.
- MCP path is `/mcp` only. No `/mcp/` redirect. No `/sse`.
- Also `GET /healthz` → `200` and body `ok\n`. Nothing else.
- One `*mcp.Server` from `BootstrapServer`, reused for every HTTP request.
- Options: stateful. Leave `Stateless`, `JSONResponse`, and `DisableLocalhostProtection` at false. Leave `CrossOriginProtection` nil. Do not set `EventStore`. Set `SessionTimeout: 30 * time.Minute` (abandoned SDK sessions leak otherwise).
- v1.6.1 `StreamableHTTPHandler` has no `Close` method. Do not add one. Do not call an unexported shutdown.
- No TLS in process. No embedded assets. Keep reading `MCP_ASSET_ROOT`, default `assets`.
- Logs stay zap on stderr. `-http` → `logger.PRODUCTION` before bootstrap. Stdio → `logger.DEVELOPMENT`, as today. Do not log request bodies. Do not write stdout.
- Leave the existing commented `zap.S().Info` line in `main.go` alone.
- `main` owns `zap.S().Fatalf`. No new `os.Exit` or `log.Fatal`.

DNS rebinding: the SDK returns 403 when the listener is loopback and `Host` is not. Do not disable that. Local flag value is `127.0.0.1:8080`. Client URL is `http://127.0.0.1:8080/mcp`. Container listens on `0.0.0.0:8080`, which is not loopback, so a normal `Host` is allowed.

Anyone who can open the port can read every resource and call every tool. That is accepted. Do not code auth to "fix" it.

## `cmd/main.go`

Keep bootstrap as it is. Add an unexported helper used by `main` and the test:

```go
func newMCPHandler(server *mcp.Server) http.Handler {
    mcpHandler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
        return server
    }, &mcp.StreamableHTTPOptions{
        SessionTimeout: 30 * time.Minute,
    })
    mux := http.NewServeMux()
    mux.Handle("/mcp", mcpHandler)
    mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        _, _ = w.Write([]byte("ok\n"))
    })
    return mux
}
```

Go 1.22+ method patterns are valid here (`go 1.26.6`). Do not add a router library.

HTTP branch, after bootstrap:

```go
ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
defer stop()

srv := &http.Server{
    Addr:              addr,
    Handler:           newMCPHandler(server),
    ReadHeaderTimeout: 10 * time.Second,
}
errCh := make(chan error, 1)
go func() {
    errCh <- srv.ListenAndServe()
}()

select {
case <-ctx.Done():
    shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := srv.Shutdown(shutdownCtx); err != nil {
        zap.S().Fatalf("http shutdown: %v", err)
    }
case err := <-errCh:
    if err != nil && err != http.ErrServerClosed {
        zap.S().Fatal(err)
    }
}
```

Stdio branch stays `server.Run(context.Background(), &mcp.StdioTransport{})`. Do not wrap stdio in the signal handler.

`Makefile`: add

```make
run-http: build ## Build and serve HTTP on 127.0.0.1:8080
	./bin/$(BINARY) -http 127.0.0.1:8080
```

Do not change `run` or other targets.

## Test

File: `cmd/http_handler_test.go`, `package main`. Keep `cmd/main_test.go`.

Do not exec the binary. Do not bind a fixed port. Do not point at repo `assets/`.

1. Temp dir. Write one file, `std.md`:

```markdown
---
uri: standards://plan/http
name: HTTP Plan Fixture
description: fixture
---

fixture
```

`BootstrapServer` rejects empty `uri` or `name`, and a `uri` with no `://`.

2. `server, err := internal.BootstrapServer(utils.NewAssetsFinder(dir))`
3. `httptest.NewServer(newMCPHandler(server))` then `defer Close()`.
4. `mcp.NewClient(&mcp.Implementation{Name: "test", Version: "0"}, nil)`
5. `Connect` with `&mcp.StreamableClientTransport{Endpoint: ts.URL + "/mcp"}` and a 10s timeout context.
6. `ListResources`. Fail unless `standards://plan/http` is present.
7. `GET ts.URL+"/healthz"` is 200 and body `ok\n`.

`httptest` is loopback and the client `Host` is loopback. If this 403s, stop. Do not set `DisableLocalhostProtection`.

Do not change `test_resource.py`.

## Docker

Root `Dockerfile` and `.dockerignore`.

`.dockerignore` must ignore `.git`, `bin`, `ref`, `coverage.out`, `.pi`. It must not ignore `go.mod`, `go.sum`, `cmd/`, `internal/`, or `assets/`.

```dockerfile
FROM golang:1.26.6-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/mcp ./cmd

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build --chown=nonroot:nonroot --chmod=755 /out/mcp /mcp
COPY --chown=nonroot:nonroot assets /app/assets
ENV MCP_ASSET_ROOT=/app/assets
EXPOSE 8080
ENTRYPOINT ["/mcp"]
CMD ["-http", "0.0.0.0:8080"]
```

If `golang:1.26.6-bookworm` does not exist, stop. Do not pick another Go version. The final image has no shell. Do not add `USER`, `HEALTHCHECK`, a package manager, or compose. The base image already runs as `nonroot`.

`docker build -t mcp . && docker run --rm -p 8080:8080 mcp`

## Docs

Separate from code. Short edits only.

- `README.md`: keep the stdio snippet. Add `make run-http`, the docker run line, and the URL `http://127.0.0.1:8080/mcp`.
- `docs/ARCHITECTURE.md`: client reaches the same bootstrap over stdio or streamable HTTP `/mcp`.
- `docs/DECISIONS.md`: replace the "stdio only" bullet with: stdio is the default; `-http` opts into streamable HTTP; no auth; TLS is outside the process; image is distroless nonroot.
- `AGENTS.md`: change "Speaks stdio JSON-RPC" to "Speaks stdio JSON-RPC unless `-http` is set." Do not rewrite the skill table.

README snippets, no tokens:

```bash
pi mcp add ben-mcp --url http://127.0.0.1:8080/mcp
claude mcp add --transport http ben-mcp http://127.0.0.1:8080/mcp
```

```json
{ "mcpServers": { "ben-mcp": { "url": "http://127.0.0.1:8080/mcp" } } }
```

```toml
[mcp_servers.ben-mcp]
url = "http://127.0.0.1:8080/mcp"
```

Cursor uses the JSON. Codex uses the TOML. Pi also accepts that JSON. Do not set `type: sse`.

## Checks

From the repo root, stop on the first failure:

```bash
make fmt && make fmt-check && make typecheck && make vet && make test && make lint && make vuln
pre-commit run --all-files
python3 .agents/skills/test-resource-document/scripts/test_resource.py
docker build -t mcp .
```

`make` is the Go workflow. The Python harness calls `go build` itself. That is existing. Do not switch the harness to `make`.

## Commits

Do not commit unless the user asks. If asked, three commits, Conventional Commits, source then tests then docs:

1. `cmd/main.go`, `Makefile`, `Dockerfile`, `.dockerignore`
2. `cmd/http_handler_test.go` only
3. `README.md`, `AGENTS.md`, `docs/` including this plan

## Do not

- Custom HTTP, JSON-RPC, SSE, session store, or handler `Close`.
- `Stateless`, `JSONResponse`, or `DisableLocalhostProtection`.
- Auth, TLS, Cobra, slog, embedded assets, dependency bumps.
- Edits under `internal/` except none. No tool, resource, or frontmatter changes.
- Stdout writes. Stdio behavior changes.
