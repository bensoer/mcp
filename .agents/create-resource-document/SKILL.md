---
name: create-resource-document
description: >-
  Creates a new resource document (a markdown standard) in the assets directory of
  the mcp server, with the correct frontmatter and URI so it is automatically
  discovered and registered as an MCP resource. Use when adding or expanding
  standards, practices, or preferences documentation that should be exposed as a
  resource endpoint.
---

# Create a Resource Document

Use this skill when you need to add a new standards / practices / preferences
document ("resource document") to the MCP server. Each markdown file in the
`assets/` directory is automatically discovered at server startup and registered
as an MCP resource — no Go code changes are required.

---

## When to use

- The user asks to "add a standard", "add a resource", "document a practice", or
  similar.
- You need to expose new guidance as an MCP resource with a `standards://` or
  `workflows://` URI.

## Step 1 — Pick a URI and file path

Resources live under `assets/` and are namespaced by their directory:

| Directory | URI prefix |
|-----------|------------|
| `assets/standards/comments/` | `standards://comments` |
| `assets/standards/git/` | `standards://git/<topic>` |
| `assets/standards/python/` | `standards://python/<topic>` |
| `assets/workflows/<category>/` | `workflows://<category>/<topic>` |

1. Decide the category (`standards` or `workflows`) and the topic.
2. The filename must be **kebab-case** and end in `.md` (e.g. `commit-messages.md`).
3. The `uri` frontmatter field is the canonical identifier — it is **independent**
   of the file path, but conventionally mirrors it. Example:
   `assets/standards/git/pull-requests.md` → `standards://git/pull-requests`.

> If no existing category fits, create a new directory under `assets/`. The
> server walks the tree recursively, so any new directory is picked up
> automatically.

## Step 2 — Add the frontmatter

Every resource document starts with YAML frontmatter delimited by `---`.
**Only `uri` and `name` are required** — the server fails fast at startup if
either is empty (see `internal/server.go:45-51`).

```yaml
---
uri: standards://git/pull-requests
name: Pull Request Standards
description: Standards for pull request titles, descriptions, and review practices.
languages:
    - all
file_types:
    - "*.*"
priority: required
related_resources:
    - standards://git/commit-messages
    - standards://git/commit-staging
---
```

| Field | Required | Notes |
|-------|----------|-------|
| `uri` | **Yes** | Canonical MCP resource URI. Must be unique across all assets. |
| `name` | **Yes** | Display name shown to clients. |
| `description` | No — but recommended | Short summary shown in resource listings. |
| `languages` | No | e.g. `all`, `python`, `go`. (Currently documentation only.) |
| `file_types` | No | e.g. `"*.*"`. (Currently documentation only.) |
| `priority` | No | e.g. `required`. (Currently documentation only.) |
| `related_resources` | No | List of URIs referencing other standards. **Must use snake_case** — the struct tag is `yaml:"related_resources"` and camelCase keys are silently dropped. |

## Step 3 — Write the body

Follow the conventions established by the existing standards:

- Start with an `#` H1 title that mirrors the `name` frontmatter field, followed
  by the `description` as a sub-heading line.
- Use `#` H1 for the title, `##` H2 for major sections, `###` H3 for sub-sections.
- Use **sentence-case** headings.
- Use `- ALWAYS` / `- NEVER` / `- PREFER` for rule statements (imperative mood).
- Include ✅ / ❌ examples where ambiguity is likely.
- Content is rendered as `text/markdown` (MIME type is hardcoded).

## Step 4 — Rebuild and restart

The resource is picked up on the **next server start** — the asset directory is
walked once during `BootstrapServer()`. After creating or editing a file:

1. Rebuild: `go build -o bin/mcp ./cmd` (or `make build`)
2. Restart any running `bin/mcp` instance so it re-walks `assets/`.

## Checklist before finishing

- [ ] URI is in the right namespace and is unique (no other file claims it)
- [ ] `uri` and `name` frontmatter fields are present and non-empty
- [ ] `related_resources` uses snake_case (not camelCase)
- [ ] Filename is kebab-case and ends in `.md`
- [ ] Body H1 mirrors the `name` field
- [ ] Rule statements use ALWAYS/NEVER/PREFER in imperative mood
- [ ] At least one related existing standard is cross-referenced
- [ ] `go vet` and `go build` still pass
- [ ] Server restarts and the new URI appears in `resources/list` (see
      `test-resource-document` skill)
