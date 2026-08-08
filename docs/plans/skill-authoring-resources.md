# Skill Authoring Resources Plan

This document captures the plan, questions, assumptions, and decisions that
drove the creation of the `standards://skill-authoring` and
`workflows://skill-authoring` MCP resources.

---

## Goal

Create a new standard and workflow resource for "skill authoring" — how Agent
Skills (`SKILL.md` files) are expected to be created and formatted, following
the [agentskills.io](https://agentskills.io) specification and best practices.

---

## Research Sources

Consulted the following upstream pages (append `.md` for markdown version):

- https://agentskills.io/home — Overview of the Agent Skills format
- https://agentskills.io/specification — Complete format specification
- https://agentskills.io/skill-creation/quickstart — Tutorial walkthrough
- https://agentskills.io/skill-creation/best-practices — Quality guidelines
- https://agentskills.io/skill-creation/optimizing-descriptions — Description eval
- https://agentskills.io/skill-creation/using-scripts — Bundled script design
- https://agentskills.io/skill-creation/evaluating-skills — Output quality evals

---

## Design Decisions

### Question 1: Standards vs. Workflow split?

**Decision:** Split into two resources.

- **`standards://skill-authoring`** — Rules. ALWAYS/NEVER/PREFER/AVOID for
  SKILL.md format, frontmatter fields, naming, context management, specificity
  calibration, instruction patterns, script design, and description optimization.
- **`workflows://skill-authoring`** — Process. Step-by-step workflow for creating
  a skill: extract from real task → draft → refine with execution → optimize
  description → evaluate with evals → iterate.

### Question 2: How opinionated should the standard be?

**Decision:** Enforce the full agentskills.io specification. Include source
references to the upstream pages so agents can consult them as the spec evolves.

### Question 3: Branch name?

**Decision:** `docs/skill-authoring-standards` in a separate worktree at
`~/Documents/PROJECTS/worktrees/skill-authoring`.

---

## Assumptions

- Frontmatter format follows existing project conventions: `uri`, `name`,
  `description`, `languages`, `file_types`, `priority`, `related_resources`.
- Content style matches existing resources: imperative ALWAYS/NEVER/PREFER/AVOID
  for standards, numbered step-by-step for workflows.
- Two resource documents (one standard, one workflow).
- Git remote is `origin`; PR created via `gh`.
- Both documents are doc files (under the commit staging standards) and can be
  batched in one commit (< 500 lines total).

---

## Implementation Summary

1. Created worktree at `~/Documents/PROJECTS/worktrees/skill-authoring` on
   branch `docs/skill-authoring-standards`.
2. Created `assets/standards/skill-authoring/skill-authoring.md` (222 lines)
   covering: directory structure, SKILL.md frontmatter requirements (`name`,
   `description`, `license`, `compatibility`, `metadata`, `allowed-tools`),
   progressive disclosure, context management, specificity calibration,
   instruction patterns (gotchas, templates, checklists, validation loops,
   plan-validate-execute), script design, and description optimization.
3. Created `assets/workflows/skill-authoring/skill-authoring.md` (155 lines)
   with a 5-step workflow: Extract Real Expertise → Draft the Skill → Refine
   with Real Execution → Optimize the Description → Evaluate Output Quality.
4. Built and tested the server — both resources registered without errors.
5. Committed, pushed, and created PR #4.