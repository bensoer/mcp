# OpenCode & Agent Developer Experience Plan

This document captures the decisions, assumptions, and implementation plan for
improving developer and agent onboarding in the mcp repository.

---

## Goal

Add the necessary skills, rules, and AGENTS.md improvements so developers and
agents can easily navigate and work within the repository — without relying on
any gitignored content (e.g., `ref/`).

---

## Design Decisions

### Decision 1: Skip `opencode.json`

`opencode.json` is opencode-specific project-level configuration. Including it
in the repo would couple the repo to one agent platform.

**Decision:** Do NOT create `opencode.json`. Instead, document setup
instructions in `AGENTS.md` so each developer configures their preferred agent
platform (opencode, Cursor, Claude Code, etc.) themselves.

### Decision 2: Keep `.agents/` for skills

The existing skills (`create-resource-document`, `test-resource-document`) live
in `.agents/`, which is discoverable by multiple agent platforms (opencode,
Cline, Cursor via external skill paths).

**Decision:** Keep `.agents/` as the skill directory. Add a new
`server-development` skill for Go development workflows specific to this repo.

### Decision 3: Rules in `.agents/rules/`

Rules specific to THIS repository (not general coding standards served by the
MCP server) need a home. `.agents/rules/` is agent-agnostic and follows the
existing `.agents/` convention.

**Decision:** Create `.agents/rules/` with repo-specific conventions:
- Go internal package structure
- Resource document frontmatter conventions
- Development workflow (build, vet, test cycle)

### Decision 4: Commands documented in AGENTS.md, not `.opencode/command/`

Agent-specific slash-commands (`.opencode/command/`) couple to one platform.

**Decision:** Document common workflows as shell commands in `AGENTS.md` rather
than as platform-specific command files. These are usable by any developer
regardless of their agent or editor.

### Decision 5: No rules in `.opencode/rules/`

Platform-specific rules folders couple to one agent.

**Decision:** All rules go in `.agents/rules/` for agent-agnostic discoverability.

---

## Assumptions

- `ref/` is treated as absent — no references to its content in any new file.
- The binary path is `./bin/mcp` and must run from the project root.
- Skills in `.agents/` are auto-discovered by the agent platform (standard path).
- Worktrees at `~/Documents/PROJECTS/worktrees/` are the preferred branch workflow.
- No python standards are applied to the Go server code — only Go conventions.

---

## Files to Create/Modify

| File | Action | Purpose |
|------|--------|---------|
| `AGENTS.md` | Rewrite | Comprehensive dev/agent guide with setup, workflows, skill reference, architecture overview |
| `.agents/server-development/SKILL.md` | Create | Go development workflow for working on the MCP server code |
| `.agents/rules/go-conventions.md` | Create | Go conventions specific to this repo's internal package structure |
| `.agents/rules/resource-conventions.md` | Create | Rules for creating/maintaining resource documents in `assets/` |
| `docs/plans/opencode-agent-setup.md` | Create | This plan document |

---

## Implementation Summary

1. Create this plan document in `docs/plans/`.
2. Rewrite `AGENTS.md` with:
   - Agent platform setup instructions (opencode, Cursor, Claude Code)
   - Documented skills reference (when to use each)
   - Common workflow commands (build, test resources, add resource)
   - Architecture overview (data flow from assets to MCP resources)
   - Debugging/troubleshooting guide
3. Create `.agents/server-development/SKILL.md` covering:
   - Build/vet cycle (`make build && make vet`)
   - Internal package conventions (`internal/` layout)
   - How `BootstrapServer` works
   - Logger usage
   - Testing with `test_resource.py`
4. Create `.agents/rules/go-conventions.md` covering:
   - Package naming and placement in `internal/`
   - Interface conventions (used in `utils/assets.go`)
   - Error handling patterns
5. Create `.agents/rules/resource-conventions.md` covering:
   - Frontmatter required/optional fields
   - URI naming conventions
   - Directory structure under `assets/`
   - How to test a new document
6. Verify: `make build && make vet` passes.
7. Commit, push, create PR.