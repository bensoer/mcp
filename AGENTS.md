# AGENTS.md — mcp

Index for agents working in this repository. Do not load `assets/` to learn how to change this repo.

## What this repo is

Go 1.26 MCP server. One binary: `bin/mcp`. Speaks stdio JSON-RPC.

| Path | What it is | Touch it when |
|------|------------|---------------|
| `cmd/`, `internal/` | Server code | Changing behavior, tools, or tests |
| `assets/` | Payload served as `standards://` and `workflows://` | Adding a standard or workflow |
| `.agents/rules/` | Rules for this repo | Checking conventions |
| `.agents/skills/` | Procedures | Server, resource, or harness work |
| `.agents/commands/` | Slash prompts (`/validate`, `/build`, …) | Running a named check |
| `docs/` | Context, architecture, decisions, plans | Needing why, not how |
| `Makefile` | Build, test, lint, vuln | Always. Do not use raw `go` as the main workflow |

`assets/standards/python/` and similar files are documents this server serves. They are not this repo's language or build instructions.

This server is often connected as `ben-mcp-server`. A broken build or bad frontmatter can break later lookups.

## Which skill

| Task | Skill |
|------|--------|
| Go changes | `.agents/skills/server-development` |
| New or edited `assets/` file | `.agents/skills/create-resource-document`, then `.agents/skills/test-resource-document` |
| Unsure what a path is | `.agents/rules/repo-boundaries.md` |

## Required checks

Details: `.agents/rules/validation.md`.

```bash
make fmt && make fmt-check && make typecheck && make vet && make test && make lint && make vuln
pre-commit run --all-files
```

After `assets/` edits, also run:

```bash
python3 .agents/skills/test-resource-document/scripts/test_resource.py
python3 .agents/skills/test-resource-document/scripts/test_resource.py --read <uri>
```

## Docs

| File | Use |
|------|-----|
| `docs/PROJECT_CONTEXT.md` | What the product is for |
| `docs/ARCHITECTURE.md` | How the process is wired |
| `docs/DECISIONS.md` | Trade-offs not obvious from code |
| `docs/plans/` | Historical plans, not the current source of truth |

## Commits

Conventional Commits. Separate source, tests, and docs. See `standards://git/commit-messages` and `standards://git/commit-staging`.
