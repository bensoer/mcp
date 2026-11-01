# Upgrade go-sdk to v1.8.0

Implement this file. It is the task, not a historical note. `AGENTS.md` says `docs/plans/` is not the source of truth. That does not apply here.

Do the work in its own git worktree. Do not edit the checkout that currently has this file untracked on `main`. When the checks below pass, commit, push, and open a PR.

Do not redesign the server. Do not add auth, a second transport, a second MCP path, or new tools. If a step is blocked, stop and report. Do not invent a workaround.

`docs/plans/http-streamable-transport.md` pinned `go-sdk` v1.6.1 and stateful HTTP. Ignore those two constraints. Everything else that plan locked (flag, paths, healthz, no auth, DNS rebinding, logging) still applies unless this file overrides it.

## Decision already made

Operator accepted losing stateful HTTP sessions. This server has no per-session application state. Resources and tools are registered once at startup and never change. Do not keep `SessionTimeout` to preserve sessions.

Do not drop legacy clients. Pi 1.4.0 (`@earendil-works/pi-mcp`) only speaks through `2025-11-25` and still sends `initialize`. The Python harness sends `initialize` at `2024-11-05` over stdio. Both must keep working.

## Target

| | Now | After |
|---|---|---|
| Module | `github.com/modelcontextprotocol/go-sdk` v1.6.1 | **v1.8.0** (latest stable on proxy.golang.org as of this plan; do not use a `-pre` tag) |
| Newest spec the SDK negotiates | `2025-11-25` | `2026-07-28`, plus `2025-11-25`, `2025-06-18`, `2025-03-26`, `2024-11-05` |
| HTTP | stateful streamable, `SessionTimeout: 30m` | **stateless** streamable. Required or HTTP cannot serve `2026-07-28` |
| Go | `go 1.26.6` | unchanged |

SDK facts, do not re-derive:

- Release notes: https://github.com/modelcontextprotocol/go-sdk/releases/tag/v1.7.0 and https://github.com/modelcontextprotocol/go-sdk/releases/tag/v1.8.0
- Protocol model: `$(go env GOMODCACHE)/github.com/modelcontextprotocol/go-sdk@v1.8.0/docs/protocol.md`
- `StreamableServerTransport.SupportsProtocolVersion`: `2026-07-28` only when `Stateless` is true. Same file: `mcp/streamable.go`.
- Stateless handler: no `Mcp-Session-Id`, GET and DELETE return 405, each request is a temporary session. Legacy `initialize` is still accepted. Pi treats GET 405 as "no server-to-client stream" and continues (`pi-mcp` `openSseStream`).
- Do not set `ServerOptions.SupportedProtocolVersions`. Empty means advertise every version the SDK implements. Narrowing would reject a client we still need.
- Do not set `MCPGODEBUG`. Do not add OAuth.

## Preserve

These are in use. A compile-green bump that drops one of them is a failure.

- Default transport is stdio. `-http host:port` opts into streamable HTTP. Flag lives on `flag.NewFlagSet`, not `flag.CommandLine`. Whitespace-only `-http` is fatal. Unset means stdio.
- MCP path is `/mcp` only. `GET /healthz` returns `200` and `ok\n`. No `/sse`, no auth, no in-process TLS.
- One `*mcp.Server` from `BootstrapServer`, reused for every HTTP request.
- DNS rebinding stays on. Leave `DisableLocalhostProtection` false and `CrossOriginProtection` nil. Leave `JSONResponse` false. Do not set `EventStore`.
- `MCP_ASSET_ROOT`, default `assets`. Resource read returns the full file, including frontmatter. MIME `text/markdown`. URI, name, title, description, size stay as registered in `internal/server.go`.
- Tools, names and contracts unchanged: `discover_standards`, `discover_workflows`, `search_standards`, `search_workflows`. Registration still loops `tools.RegisterAll`. Search still matches URI, name, description, and body; tools still return metadata only.
- Logger: stdio `DEVELOPMENT`, `-http` `PRODUCTION`, zap on stderr, no request-body logs, no stdout logs. `main` owns `zap.S().Fatalf`. Leave the commented `zap.S().Info` line in `cmd/main.go` alone.
- Python harness (`.agents/skills/test-resource-document/scripts/test_resource.py`) stays on stdio `initialize` / `2024-11-05`. Do not rewrite it onto `2026-07-28`.
- `go 1.26.6` directive. The `tool (` block in `go.mod` (golangci-lint, govulncheck) must survive `go mod tidy`.

Wire changes that are not feature loss, do not opt out:

- Unknown resource error code moves from `-32002` to `-32602` (SEP-2164). Repo tests only require a non-nil error.
- v1.8 inbound limits (4 MiB HTTP body, 16 MiB stdio line, JSON nesting cap) stay at SDK defaults. This server's requests are tiny; largest asset is ~10 KiB.

## Worktree

No ticket number was given. Do not block on one. This repo's branches are descriptive (`fix/go1.27-lint`, `http-streamable-transport`). Use `chore/go-sdk-v1.8`.

From the repo that contains this plan (`/home/bsoer/Documents/PROJECTS/VSCodeProjects/mcp`):

```bash
git fetch origin
git pull origin main
git worktree add -b chore/go-sdk-v1.8 ../mcp-go-sdk-v1.8 origin/main
```

If `../mcp-go-sdk-v1.8` already exists, stop and report. Do not reuse a dirty tree.

Copy this plan into the worktree (`docs/plans/go-sdk-v1.8-upgrade.md`). It is not on `main`. All edits, tests, commits, and the PR happen with cwd = the worktree. Leave the original checkout on `main`.

Branch off updated `origin/main` only. Do not rebase, amend, or force-push (`standards://git/operations`, `standards://git/branches`).

## Code changes

1. `go get github.com/modelcontextprotocol/go-sdk@v1.8.0` then `go mod tidy`. Confirm the require line is v1.8.0 and the `tool (` block is intact.
2. `cmd/main.go` `newMCPHandler`: set `Stateless: true`. Remove `SessionTimeout`. Do not set other `StreamableHTTPOptions` fields. `StreamableHTTPHandler` still has no public `Close`; do not call `closeAll`.
3. Touch Go call sites only if v1.8.0 does not compile or a test fails. Current API usage (`NewServer`, `AddResource`, `AddTool`, `Run` + `StdioTransport`, `NewStreamableHTTPHandler`, in-memory and streamable clients) matches v1.8.0 signatures. Do not restyle handlers.
4. If `AddTool` panics on schema generation, fix the tool input/output types so the same four tools still register. Do not delete a tool or hand-write a weaker schema that drops fields.
5. Docs, only these:
   - `README.md`: HTTP line that says spec `2025-11-25`. Say the server negotiates; newest is `2026-07-28`; older clients still use `initialize` through `2025-11-25`. HTTP is stateless (no session id; GET `/mcp` is 405).
   - `docs/DECISIONS.md`: SDK v1.8.0; HTTP stateless so `2026-07-28` can be served; legacy handshake kept; no auth.
   - `docs/PROJECT_CONTEXT.md`: the "stdio JSON-RPC" sentence is already incomplete. Mention streamable HTTP as well. One sentence.
   - Do not edit other plans, skills, or `assets/`.

## Tests

Review every test against v1.8.0 behavior. Update assertions that assume v1.6.1, a session id, `SessionTimeout`, GET `/mcp` as a live SSE stream, error code `-32002`, or a negotiated version of `2025-11-25` when the peer is a v1.8 client. Do not delete a test to go green. Do not weaken a check that still describes current behavior (resource body includes frontmatter, tool names, search matching, healthz, flag parsing, DNS rebinding on loopback).

Files:

| File | What to re-check |
|------|------------------|
| `cmd/http_handler_test.go` | Loopback client must not 403. v1.8 `mcp.Client` on a stateless handler negotiates `2026-07-28`. That is fine. Do not require `Mcp-Session-Id`. |
| `cmd/main_test.go` | Flag parsing, whitespace `-http`, no global `flag.CommandLine` flag. |
| `internal/server_test.go` | In-memory client will negotiate `2026-07-28`. List, read, unknown resource still work. Unknown resource is a non-nil error; do not pin `-32002`. |
| `internal/tools/tools_test.go` | All four tools still callable. Output shape unchanged. |
| `internal/tools/helper_test.go`, `internal/tools/search_terms_test.go` | Search matching unchanged. |
| `internal/utils/assets_test.go`, `internal/models/frontmatter_test.go`, `internal/logger/logging_test.go` | Unrelated to the SDK. Run them. Change only if the bump breaks them. |
| `.agents/skills/test-resource-document/scripts/test_resource.py` | Stays on stdio `initialize` / `2024-11-05`. Do not rewrite it. If it fails, fix the server. |

Add one HTTP test that does **not** use `mcp.Client`. POST a legacy `initialize` (`protocolVersion` `2025-11-25` or `2024-11-05`) to `/mcp` on the handler from `newMCPHandler`, then `resources/list` or `tools/list`, and assert the fixture resource or the four tool names come back. A v1.8 `mcp.Client` will not catch a legacy regression.

Also assert GET `/mcp` is 405 and GET `/healthz` is still `200` / `ok\n`. Do not require GET `/mcp` to return 200.

## Verify

Cwd is the worktree. Do not use raw `go` as the workflow (`standards://go/tooling`). `make help` prints the targets. Use those names. Do not run `make run` or `make run-http`; they block on stdio or a listen socket.

Stop on the first failure. See `.agents/rules/validation.md`.

```bash
make help
make fmt
make fmt-check
make typecheck
make vet
make test
make lint
make vuln
make check
pre-commit run --all-files
make build
python3 .agents/skills/test-resource-document/scripts/test_resource.py
python3 .agents/skills/test-resource-document/scripts/test_resource.py --read standards://go/tooling
```

`make check` repeats several of those targets and adds `go test -race`. Run it anyway. A green `make test` is not a substitute for `make lint`, `make vuln`, or `make check`.

## Commits, push, PR

After the commands above pass, commit in the worktree. Conventional Commits (`standards://git/commit-messages`). Separate commits (`standards://git/commit-staging`):

1. Source: `go.mod`, `go.sum`, `cmd/main.go`, any non-test Go fix. `go.mod` and `go.sum` stay together. If that commit exceeds 500 lines because of `go.sum`, do not split the sum file; note the count in the PR body.
2. Tests: every `*_test.go` change, and the Python harness only if it had to change.
3. Docs: this plan, `README.md`, `docs/DECISIONS.md`, `docs/PROJECT_CONTEXT.md`.

Do not amend. Do not mix those three groups.

```bash
git push -u origin chore/go-sdk-v1.8
gh pr create --base main --head chore/go-sdk-v1.8
```

PR title is one Conventional Commits subject and must list every commit type present, priority order `feat`, `fix`, `refactor`, `build`, `docs`, `test`, `chore` (`standards://git/pull-requests`). Expected shape if the three commits above are the only ones: `build, test, docs: upgrade MCP go-sdk to v1.8.0`.

PR body must have all four sections, using `- None` when empty:

```markdown
### Added
### Removed
### Updated
### Deleted
```

Say what was verified (`make help` targets that were run, Python harness). Link this plan. Do not merge the PR.

## Stop and report

- Legacy `initialize` fails on stdio (Python harness) or on the new HTTP test.
- Serving `2026-07-28` requires something other than `Stateless: true` (second port, auth, dropping a tool or resource field).
- `go 1.26.6` cannot build v1.8.0.
- `go mod tidy` drops the `tool (` block.
- A Makefile verification target fails and the fix is not a straight test or call-site update.
- The worktree path already exists.

Do not push or open a PR if you stopped.
