# Repository boundaries

| Tree | Meaning |
|------|---------|
| `cmd/`, `internal/`, `Makefile`, `go.mod` | This product. Go only. |
| `assets/` | Markdown the server serves. Not build instructions. |

`assets/standards/typescript/` does not make this a TypeScript repo. Do not run npm, pip, or pytest here.

Repo rules live in `AGENTS.md` and `.agents/`. A new `assets/**/*.md` file does not need a Go change. Rebuild so the next process start picks it up.

This server is often the `ben-mcp-server` used by the agent editing it. Validate the build and resource harness after asset edits.
