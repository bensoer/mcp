# HTTP API key auth and rate limits

Status: implemented on `feat/http-api-key-auth`.

Implement this file. It is the task. `AGENTS.md` says `docs/plans/` is not the source of truth. That does not apply while implementing this file. After the code lands, update the living docs listed below so they match the code.

Do not redesign. If a locked rule blocks you, stop and report. Do not invent a workaround.

Confirmed: auth is opt-in via `-auth`; register failures are a uniform `403`; rate limits cover email and remote address; `/mcp` accepts `Authorization: Bearer` only. Branch name is `feat/http-api-key-auth`.

## Where to work

Do not edit the checkout that contains this file if that checkout is `main` or already has other work. Implement, test, commit, and open the pull request from a new worktree on its own branch.

Preferred worktree directory: `~/Documents/PROJECTS/worktrees/` (see `docs/plans/opencode-agent-setup.md`). Branch naming: `standards://git/branches`. Name is already chosen: `feat/http-api-key-auth`. Do not ask again. Do not invent a ticket.

From this repo, after `git fetch origin`:

`git worktree add -b feat/http-api-key-auth ~/Documents/PROJECTS/worktrees/mcp-http-api-key-auth origin/main`

If that path or branch already exists, stop. Do not reuse it.

This plan may not be on `origin/main`. Copy it into the worktree at `docs/plans/http-api-key-auth.md` before editing. All later commands use the worktree as cwd. Leave the original checkout untouched.

## Read first

| Path | Why |
|------|-----|
| `cmd/main.go` | Flag set, stdio vs `-http`, `newMCPHandler`, shutdown. |
| `cmd/http_handler_test.go` | Existing HTTP tests must keep passing with auth off. |
| `docs/DECISIONS.md` | Today: HTTP is stateless, no auth, TLS outside the process. |
| `docs/ARCHITECTURE.md` | Transport wiring. Add the auth path; do not redraw unrelated parts. |
| `README.md` | "Run over HTTP" says no auth. Update that section only. |
| `.agents/rules/validation.md` | Required checks. |
| `.agents/rules/go-conventions.md` | Repo Go layout. Wins over a conflicting standard. |
| `.agents/skills/server-development/SKILL.md` | Build cycle. Do not treat its stdio diagram as the transport spec. |
| `Dockerfile` | `CGO_ENABLED=0`, distroless, `CMD ["-http", "0.0.0.0:8080"]`. Do not change the image default. |
| `charts/mcp/` | Auth is a values setting. Read `values.yaml`, `templates/deployment.yaml`, `templates/hpa.yaml`, `templates/NOTES.txt`, and `tests/`. |
| `test/e2e/http_test.go` | Unauthenticated `/mcp` and `/healthz`. Must keep passing. |

Standards, read from `standards-and-workflows-mcp` rather than copying them here:

- `standards://go/architecture` — `internal/<feature>/`, thin `cmd`, constructors return structs, no I/O in `New`, no new `init()`.
- `standards://go/logging` — zap; log an error once at the edge; never log and return the same error. This repo already uses zap. Pass `*zap.SugaredLogger` into feature constructors. `cmd` may keep `zap.S()` until it has a context.
- `standards://go/testing` — `package` tests, stdlib `testing`, hand-written fakes, no production export added only for a test.
- `standards://go/syntax` — errors and context. Read only if naming or error wrapping is unclear.

Do not add a served `assets/` document. This task is server behavior, not a new `standards://` or `workflows://` resource.

SDK is `github.com/modelcontextprotocol/go-sdk` v1.8.0. Do not bump it. Do not change `Stateless: true` in `newMCPHandler`. Auth is a bearer token on each HTTP request, not an MCP session.

## Goal

HTTP only. Stdio stays local, unauthenticated, and unlimited.

When the operator sets `-auth` to a BoltDB path, HTTP requires an API key on `/mcp`, attributes the caller to the email used at registration, and rate-limits that email. Registration is a one-shot `POST /register`. The plaintext key is returned once and is not stored.

When `-auth` is unset, HTTP behaves as it does today: no auth, no `/register`, no rate limit. The image `CMD` and the chart default stay that way, so `test/e2e` keeps calling `/mcp` with no bearer.

No `standards://` or `workflows://` guidance covers Helm. Follow the chart that is already in the repo.

## Locked behavior

### Flags

Add flags only on the existing `flag.NewFlagSet` in `cmd/main.go`. Do not use the global `flag.CommandLine`. Do not add Cobra or Viper.

| Flag | Meaning |
|------|---------|
| `-auth <path>` | Bolt file. Empty: auth off. Whitespace-only: fatal in `main`. |
| `-rate-limit <n>` | Sustained requests/second per identity. `0`: use the defaults below. Negative: fatal. |

`-auth` or `-rate-limit` with no `-http`: fatal. Do not open the database in stdio mode.

Open the database from `main` after flags parse, not from a constructor. Mode `0600`. Create the file if missing. Do not create missing parent directories. Close it on HTTP shutdown, including the signal path in `cmd/main.go`.

Library: `go.etcd.io/bbolt`. It is the maintained BoltDB API and is pure Go, so `CGO_ENABLED=0` still links. Do not use archived `github.com/boltdb/bolt`. Rate limit: `golang.org/x/time/rate`, in memory, per process. A second replica would get its own counters; do not add Redis.

### Registration

`POST /register` exists only when `-auth` is set. It is not behind the API key. There is no login, recovery, list, rotate, or delete endpoint.

Body: JSON object, `Content-Type: application/json`, at most 1 KiB. Required field: `email` (string). Ignore unknown fields.

Normalize email before lookup: trim space, lowercase. Reject empty after normalize. Do not add a stricter format check.

Success, and only success: `200` and `{"apiKey":"<hex>"}`. Generate 32 bytes from `crypto/rand`, hex-encoded (64 chars). No prefix. This is the only time the plaintext key exists.

Store SHA-256 of the raw token bytes (the 32 bytes, not the hex string) as hex. Never write the plaintext key or the hex key.

One Bolt update transaction, two buckets:

| Bucket | Key | Value |
|--------|-----|-------|
| `emails` | normalized email | hex SHA-256 |
| `keys` | hex SHA-256 | normalized email |

If the email key already exists, write nothing.

Every non-success on this route is `403` with `{"error":"forbidden"}`. That includes bad JSON, wrong content type, missing email, duplicate email, empty body, and any method other than `POST`. Do not tell the client which case it was. The server log may include the normalized email for a duplicate or a successful register. Never log the key, the hash, or the request body.

### Calling `/mcp`

Accept only `Authorization: Bearer <hex key>`. Do not accept `X-API-Key` or a query parameter.

Hash the presented token the same way and read `keys`. Missing header, bad scheme, unknown token: `401` with an empty body. Do not echo the token. Log the remote address, not the header value.

On success, the identity is the stored email. Do not put the token on the context.

`GET /healthz` stays open: `200` and `ok\n`. No auth, no rate limit. Do not change its body.

Any other path keeps today's routing. Do not add routes.

### Rate limit

In-memory token buckets. Burst is `2 * rps`, rounded up, minimum 1.

| Traffic | Key | Default rps |
|---------|-----|-------------|
| Authenticated `/mcp` | normalized email | 10 |
| `POST /register` and unauthenticated `/mcp` | `r.RemoteAddr` only | 1 |

Do not trust `X-Forwarded-For`. `-rate-limit` overrides both rps values; burst stays `2 * rps`.

Over limit: `429`, empty body, `Retry-After` in whole seconds. Log email when the request was authenticated; otherwise log remote address. Never log the token.

Apply the IP bucket before doing register work and before reporting `401`, so a flood cannot skip the limit. `429` wins over `403` and `401`.

### Process shape

Leave `newMCPHandler` as the MCP plus health handler. Wrap it only when a store is open. The wrapper owns `/register`, the bearer check, and the limiter. `BootstrapServer` stays unchanged.

`cmd/main.go` wires flags, logger, store, and the wrap. Key generation, hashing, Bolt updates, and HTTP decisions live under `internal/auth/`. Follow `standards://go/architecture` for file split: one behavior-owning struct per file, interfaces beside the consumer if a test needs a fake. Do not put a Bolt client in `main`.

Logs stay zap on stderr. No request bodies. No stdout.

### Helm

`auth.enabled` defaults to false. Default render must stay what it is today: no PVC, no container `args`, HPA still 1–3.

Several pods cannot use one Bolt file. `Open` takes an exclusive lock. A second `Open`, including read-only, cannot use the file while that lock is held. `ReadWriteMany` does not remove the lock. Do not add a sidecar or a second replica.

When `auth.enabled` is true:

- Render a PersistentVolumeClaim. Access mode `ReadWriteOnce`. Default size `1Gi`. Empty `storageClass` omits `storageClassName` so the cluster default is used. Mount the volume at `/data`. The flag value is `/data/keys.db`. Do not use `emptyDir`. Do not use a `subPath` for the database file.
- Container `args` replace the image `CMD`. Set `-http`, `0.0.0.0:8080`, `-auth`, `/data/keys.db`. When auth is off, leave `args` unset.
- Keep `readOnlyRootFilesystem: true`. The claim is the only writable mount. Set pod `fsGroup: 65532` only when that mount exists, so uid 65532 can create the file. Do not change the auth-off security context.
- Replicas are 1. Do not render the HPA. If `autoscaling.enabled` is true or `replicaCount` is not 1, `fail` the render. The message must say Bolt allows one process.
- Probes stay `GET /healthz` with no bearer. Do not enable ingress. Do not widen `templates/networkpolicy.yaml`; same-namespace traffic to 8080 already includes `/register`.
- The claim is a normal release resource. Do not set `helm.sh/resource-policy: keep`. Uninstall deletes it.

`charts/mcp/templates/NOTES.txt` must say: auth is off unless enabled; when on, there is one replica, keys are on the claim, register is `POST /register`, `/mcp` needs the bearer, and the HPA must stay off.

In-memory rate buckets are the cluster limit only because the chart forces one pod. A restart clears buckets, not keys. Do not store the limiter in Bolt.

## Docs to update

- `README.md` — HTTP section: auth is off unless `-auth` is set; how to register; `Authorization: Bearer`; keys are not recoverable; stdio has neither auth nor limits. Helm section: `auth.enabled` is off by default; turning it on adds the claim and forces one replica.
- `docs/DECISIONS.md` — replace the blanket "no auth" line with: HTTP auth is opt-in via `-auth <bolt path>`; unset `-auth` keeps the open HTTP server; stdio has no auth. Chart auth is the same opt-in. One replica, because Bolt's file lock is exclusive. The default chart still has no auth and still has the HPA.
- `docs/ARCHITECTURE.md` — one short note that an optional `internal/auth` wrap sits in front of `/mcp` when `-auth` is set, and that the chart mounts a claim for that file only when `auth.enabled` is true.

Do not edit `docs/plans/http-streamable-transport.md`.

## Tests

Use a temp Bolt file. Auth-off cases stay on the current handler with a nil store.

Cover: first register returns a key and a second register of that email is a uniform `403`; malformed register is `403` and does not create a row; `/mcp` without a bearer is `401` when the store is open; a valid bearer can list resources the same way `cmd/http_handler_test.go` already does; a wrong bearer is `401`; `/healthz` stays open with the store open; authenticated traffic past the bucket is `429` with `Retry-After`; register flood is `429` before a `403`; stdio flag combination is rejected without opening Bolt (flag parse test is enough); restart against the same file still accepts the original key and still rejects a second register.

Assert the Bolt file does not contain the plaintext hex key.

Do not weaken `Stateless` or turn the default e2e install into an authenticated server.

Chart tests in `charts/mcp/tests/`: auth off still has no claim and no `args`, and the HPA still renders; auth on renders the claim, the four args, `replicas: 1`, `fsGroup: 65532`, and no HPA; auth on with autoscaling or `replicaCount` other than 1 fails the render. `make helm-lint` must pass.

## Checks

Run these in the worktree. Stop on the first failure. This is the local validation `workflows://create-pull-request` Step 1 must use. Discover nothing weaker from the Makefile.

From `.agents/rules/validation.md`: `make fmt`, `make fmt-check`, `make typecheck`, `make vet`, `make test`, `make lint`, `make vuln`, then `pre-commit run --all-files`.

The test list above is the minimum. Add or extend tests for every new flag, status code, store rule, and wrapper path you implement. Do not delete or narrow an existing test to make the suite pass. A new behavior with no test is incomplete.

Update living docs so they match the code: `README.md`, `docs/DECISIONS.md`, `docs/ARCHITECTURE.md`, and this plan file in the worktree. If a Makefile comment, help line, or other repo doc would be false after the change, update that too. Do not add an `assets/` resource.

Also run `make helm-lint`. Do not change the Dockerfile or the default e2e client. `make e2e` installs the default chart, which stays auth-off. Run it when `docker`, `minikube`, `helm`, and `kubectl` are present, because CI runs it. If one is missing, leave that run for CI and do not record it as passed.

## Pull request

Commit in the worktree only after the checks above pass. Follow `standards://git/commit-messages` and `standards://git/commit-staging` (separate source, tests, and docs). Include this plan in the docs commit. Do not commit from the original checkout.

Then follow `workflows://create-pull-request` to the end. Read a related resource only when its `when` matches that step. In particular: `standards://git/pull-requests` for the title and description, `standards://git/operations` for the push, `workflows://investigate-ci-failures` if a check fails. Do not merge, enable auto-merge, or approve. Do not open a draft.

## Not doing

- Implementing in the current checkout, or reusing an existing worktree or branch.
- Auth or limits on stdio.
- OAuth, JWT, sessions, cookies, users, passwords, email delivery.
- Fetching, rotating, or deleting a key through HTTP.
- More than one pod on the Bolt file, including `ReadWriteMany` or a read-only second replica.
- `emptyDir` for the key file, or enabling auth in the image `CMD`.
- Changing MCP protocol version, session mode, or `/healthz`.
- A new `assets/` resource.
- Merging the pull request.
