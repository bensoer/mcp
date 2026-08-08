# Duplicate URI Detection at Bootstrap

Detect and reject duplicate URIs during `BootstrapServer` to prevent silent
resource overwrites.

## Problem

`mcp.Server.AddResource()` silently overwrites any previously registered
resource with the same URI. If two asset files declare the same `uri` in
their frontmatter, the first one's content is lost with no warning.

This already happened: `assets/standards/python/architecture.md` declared
`uri: standards://python/syntax` — the same URI as `syntax.md` — making the
architecture document invisible to all clients (fixed in a prior change).

## Proposed Change

Maintain a `map[string]string` (URI → file path) during bootstrap. Before
calling `server.AddResource()`, check if the URI is already in the map. If
so, fail with an error identifying both the current and conflicting file:

```
asset standards/python/architecture.md has duplicate uri "standards://python/syntax"
  (already used by standards/python/syntax.md)
```

### Files to change

- `internal/server.go` — add URI tracking map and duplicate check in the
  BootstrapServer loop

### Edge cases

- Two files in entirely different directories with the same URI — caught
- Same URI used across different resource types (standards vs workflows) — caught
- File that references itself (same path, somehow) — caught by the map check

### Risks

None. This is a startup-time validation. If all URIs are unique (as they
should be), the map is just O(n) memory overhead proportional to the number
of assets.

## Not Doing

- Case-insensitive URI comparison — MCP URIs are case-sensitive
- Normalizing trailing slashes — assets use consistent URI patterns