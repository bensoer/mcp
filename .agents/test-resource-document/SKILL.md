---
name: test-resource-document
description: >-
  Verifies that a newly created or edited resource document is correctly
  discovered, registered, and served by the MCP server. Use after creating a
  resource document (see create-resource-document skill) to confirm the file
  parses, the URI appears in resources/list, and the content is readable via
  resources/read.
---

# Test a Resource Document

Use this skill after creating or editing a resource document to confirm it is
correctly registered and served by the MCP server. The server walks the `assets/`
directory once at startup (`BootstrapServer`), so a **rebuild + restart** is
required before the new/updated document is visible.

---

## When to use

- Immediately after creating a new resource document.
- After editing an existing resource's frontmatter or content (to confirm the
  update is reflected).

## Step 1 — Verify the build compiles

Syntax errors in a markdown file won't break the build, but confirming a clean
build ensures nothing else is broken:

```bash
go vet ./...   # catches Go-side issues
go build ./... # ensures the binary compiles
```

## Step 2 — Rebuild the binary

```bash
make build
# or: go build -o bin/mcp ./cmd
```

This ensures `bin/mcp` is current. The server reads `assets/` from the process
CWD (or `MCP_ASSET_ROOT` if set), so the binary picks up the new file on startup.

## Available scripts

- **`scripts/test_resource.py`** — Self-contained Python harness that rebuilds the
  binary (optional), starts the mcp server over stdio, sends the MCP handshake,
  and lists resources or reads a specific URI. JSON goes to stdout; diagnostics
  go to stderr. No dependencies required (stdlib only).

## Step 3 — List resources

Start the server over stdio and send an MCP JSON-RPC `resources/list` request.
Use the bundled test script — it handles the stdio handshake and EOF lifecycle
for you:

```bash
# From the skill directory:
python3 scripts/test_resource.py --build        # rebuilds binary, then lists
python3 scripts/test_resource.py                 # lists using existing binary
```

**Expect:** the new URI appears in the JSON `resources` array with the correct
`name` and `description`.

### How it works (for transparency)

The server reads stdin to EOF then exits. The script:
1. Writes `initialize` + `notifications/initialized` + `resources/list` to the
   server's stdin.
2. Waits briefly so the server processes the buffered requests.
3. Closes stdin (EOF → server shuts down).
4. Reads and prints the JSON-RPC responses from stdout.

## Step 4 — Read the resource content

Use the `--read` flag to fetch a specific resource by URI:

```bash
python3 scripts/test_resource.py --read standards://git/pull-requests
```

This sends `resources/read` and prints JSON metadata (size, MIME type) to stdout,
plus the full document text to stderr.

Example JSON output:
```json
{
  "uri": "standards://git/pull-requests",
  "mime_type": "text/markdown",
  "size": 2878,
  "frontmatter_present": true
}
```

**Expect:** the `contents[].text` contains the full markdown body. Confirm the
frontmatter was parsed (no startup error in stderr) and the body content is
intact.

## Step 5 — Check for startup errors

Inspect stderr from the server process. A healthy startup shows no `FATAL`
messages. Common failure modes if a document is malformed:

| Symptom | Cause | Fix |
|---------|-------|-----|
| `server is closing: EOF` only, no resources | Normal — stdin closed | Handled gracefully by the script; send all requests before EOF |
| Server exits non-zero with `Failed to bootstrap` in stderr | `BootstrapServer` returned an error | Check frontmatter: empty `uri`/`name`, invalid YAML, missing `assets/` |
| URI missing from `resources/list` | File not in `assets/`, wrong CWD, or `MCP_ASSET_ROOT` mismatch | Confirm `make build` ran from project root |
| `related_resources` not populated | Used `relatedResources:` (camelCase) instead of `related_resources:` | Fix YAML key casing |

## Checklist before finishing

- [ ] `go vet` and `go build` pass with no errors
- [ ] Binary rebuilt (`make build`)
- [ ] Server starts without `FATAL` errors on stderr
- [ ] New URI appears in `resources/list` with correct metadata
- [ ] `resources/read` returns the full document content
- [ ] Frontmatter fields (`name`, `description`) match expectations
