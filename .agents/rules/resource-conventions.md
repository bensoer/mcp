# Resource Document Conventions

Conventions for creating and maintaining resource documents in the `assets/`
directory. Every `.md` file under `assets/` becomes an MCP resource at server
startup.

---

## Directory Structure

```
assets/
  standards/          # Rules — ALWAYS/NEVER/PREFER/AVOID
    <category>/       # e.g., git/, python/, comments/, skill-authoring/
      <topic>.md      # e.g., commit-messages.md, syntax.md
  workflows/          # Processes — step-by-step guides
    <category>/       # e.g., comments/, skill-authoring/
      <topic>.md      # e.g., refactoring-comments.md
```

- Standards documents define rules agents MUST follow.
- Workflow documents define processes agents follow step-by-step.
- New categories get their own subdirectory under `standards/` or `workflows/`.
- Filenames are kebab-case and end in `.md`.

## URI Naming

URIs follow the pattern `{type}://{category}/{topic}`:

| Asset path | URI |
|------------|-----|
| `assets/standards/git/commit-messages.md` | `standards://git/commit-messages` |
| `assets/standards/python/syntax.md` | `standards://python/syntax` |
| `assets/workflows/comments/refactoring-comments.md` | `workflows://refactoring-comments` |
| `assets/standards/comments/comments.md` | `standards://comments` |

- URIs are independent of file paths but conventionally mirror them.
- Single-topic categories (like `comments`) omit the topic segment.
- Multi-topic categories (like `git`) include the topic segment.
- Each URI must be unique across all assets.

## Frontmatter

### Required fields

| Field | Notes |
|-------|-------|
| `uri` | Must be non-empty. Server fails to start if empty. |
| `name` | Display name. Must be non-empty. Server fails to start if empty. |

### Optional fields

| Field | Type | Notes |
|-------|------|-------|
| `description` | string | Short summary shown in resource listings. |
| `languages` | `[]string` | e.g. `["all"]`, `["python"]`. |
| `file_types` | `[]string` | e.g. `["*.*"]`, `["*.py"]`. |
| `priority` | string | e.g. `required`. |
| `related_resources` | `[]string` | URIs of related standards. |

### Critical: `related_resources` naming

The Go struct tag is `yaml:"related_resources"` (snake_case). Using
`relatedResources:` (camelCase) in YAML will SILENTLY drop the field — the
server will start but the related resources won't be populated.

```yaml
# ✅ CORRECT
related_resources:
    - standards://git/commit-messages

# ❌ WRONG — silently dropped
relatedResources:
    - standards://git/commit-messages
```

## Content Style

- Standards use imperative directives: `ALWAYS`, `NEVER`, `PREFER`, `AVOID`.
- Workflows use numbered steps with clear actions per step.
- Include concrete examples with `# ✅ GOOD` and `# ❌ BAD` markers.
- Keep documents focused — one standard per file.

## Testing New Documents

1. Rebuild the binary: `make build`
2. List all resources and verify the URI appears:
   ```bash
   python3 .agents/test-resource-document/scripts/test_resource.py
   ```
3. Read the content and verify frontmatter and body are intact:
   ```bash
   python3 .agents/test-resource-document/scripts/test_resource.py --read standards://your/uri
   ```
4. Check stderr for any `FATAL` messages (should be none).

See the `test-resource-document` skill for detailed verification steps.