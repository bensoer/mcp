---
uri: standards://skill-authoring
name: Skill Authoring Standards
description: Standards for creating Agent Skills (SKILL.md files) — format requirements, frontmatter rules, naming conventions, context management, specificity calibration, and script design. Applies when creating, reviewing, or refactoring agent skills. Covers the full agentskills.io specification.
languages: 
    - all
file_types:
    - "*.*"
priority: required
related_resources:
    - workflows://skill-authoring
---

# Skill Authoring Standards
Standards for creating Agent Skills — format requirements, frontmatter rules, naming conventions, context management, specificity calibration, and script design.

These standards are derived from the [Agent Skills specification](https://agentskills.io/specification) and [best practices](https://agentskills.io/skill-creation/best-practices). Agents should consult the source pages for the latest guidance as the specification evolves.

## Source References

- **Specification**: https://agentskills.io/specification
- **Best Practices**: https://agentskills.io/skill-creation/best-practices
- **Optimizing Descriptions**: https://agentskills.io/skill-creation/optimizing-descriptions
- **Using Scripts**: https://agentskills.io/skill-creation/using-scripts
- **Evaluating Skills**: https://agentskills.io/skill-creation/evaluating-skills

## Directory Structure

- ALWAYS place each skill in its own directory. The directory name ALWAYS matches the `name` field in `SKILL.md`.
- ALWAYS include a `SKILL.md` file at the skill directory root. This is the only required file.

```
skill-name/
├── SKILL.md          # Required: metadata + instructions
├── scripts/          # Optional: executable code
├── references/       # Optional: documentation
├── assets/           # Optional: templates, resources
└── ...               # Any additional files or directories
```

## `SKILL.md` Frontmatter Requirements

### `name` field (required)

- ALWAYS include a `name` field. The server MUST fail to load a skill without one.
- ALWAYS use 1-64 characters.
- ALWAYS use lowercase letters, numbers, and hyphens only.
- NEVER start or end the name with a hyphen (`-`).
- NEVER use consecutive hyphens (`--`).
- NEVER use uppercase characters.
- ALWAYS ensure the `name` matches the parent directory name exactly.

```yaml
# ✅ GOOD
name: pdf-processing

# ✅ GOOD
name: data-analysis

# ❌ BAD
name: PDF-Processing  # uppercase not allowed
```

### `description` field (required)

- ALWAYS include a `description` field. The server MUST fail to load a skill without one.
- ALWAYS keep descriptions to 1-1024 characters.
- ALWAYS describe both what the skill does AND when to use it.
- ALWAYS include specific keywords that help agents identify relevant tasks.
- ALWAYS use imperative phrasing — frame the description as an instruction to the agent ("Use this skill when...").
- ALWAYS focus on user intent, not implementation mechanics.
- PREFER being explicit about when the skill applies, including cases where the user doesn't name the domain directly.

```yaml
# ✅ GOOD
description: >
  Extracts text and tables from PDF files, fills PDF forms, and merges
  multiple PDFs. Use when working with PDF documents or when the user
  mentions PDFs, forms, or document extraction.

# ❌ BAD
description: Helps with PDFs.
```

### `license` field (optional)

- PREFER keeping the license field short — either the name of a license or the name of a bundled license file.

```yaml
license: Apache-2.0
```

### `compatibility` field (optional)

- ONLY include if the skill has specific environment requirements.
- ALWAYS keep to 1-500 characters.
- ALWAYS state required system packages, network access needs, or intended products.

```yaml
compatibility: Requires git, docker, jq, and access to the internet
```

### `metadata` field (optional)

- PREFER using reasonably unique key names to avoid accidental conflicts.
- ALWAYS use a map of string keys to string values.

```yaml
metadata:
  author: example-org
  version: "1.0"
```

### `allowed-tools` field (optional)

- ALWAYS use a space-separated string of pre-approved tool names.

```yaml
allowed-tools: Bash(git:*) Bash(jq:*) Read
```

## Body Content

- ALWAYS write instructions the agent follows when the skill activates.
- ALWAYS structure content for progressive disclosure (see below).
- ALWAYS use relative paths from the skill root for file references.

## Progressive Disclosure Design

- ALWAYS keep the main `SKILL.md` body under 500 lines and under 5,000 tokens.
- ALWAYS move detailed reference material to separate files in `references/` or similar directories.
- ALWAYS tell the agent WHEN to load each referenced file — use conditional triggers ("Read references/api-errors.md if the API returns a non-200 status code") rather than generic pointers ("see references/ for details").
- NEVER load all reference files unconditionally at activation time.

## Context Management

- ALWAYS add what the agent lacks; NEVER include content the agent already knows from general training.
- NEVER explain basic concepts (what a PDF is, how HTTP works, what a database migration does).
- ALWAYS cut content if the answer to "Would the agent get this wrong without this instruction?" is "no."
- PREFER concise, stepwise guidance with a working example over exhaustive documentation.
- PREFER coherent units of work — a skill should encapsulate a single cohesive capability that composes well with other skills.
- NEVER make skills too narrow (forcing multiple skills to load for one task) or too broad (hard to activate precisely).

## Specificity Calibration

- PREFER giving the agent freedom when multiple approaches are valid and the task tolerates variation.
- ALWAYS be prescriptive when operations are fragile, consistency matters, or a specific sequence must be followed.
- PREFER providing a clear default approach with brief mention of alternatives over presenting all options as equal.
- ALWAYS favor procedures over declarations — teach the agent HOW to approach a class of problems, not WHAT to produce for one instance.

## Instruction Patterns

### Gotchas Sections

- ALWAYS include a gotchas section for environment-specific facts that defy reasonable assumptions.
- ALWAYS list concrete corrections to mistakes the agent will make without being told otherwise.
- PREFER keeping gotchas in `SKILL.md` where the agent reads them before encountering the situation.
- WHEN an agent makes a mistake you have to correct, ALWAYS add the correction to the gotchas section.

```markdown
## Gotchas

- The `users` table uses soft deletes. Queries must include
  `WHERE deleted_at IS NULL` or results will include deactivated accounts.
- The `/health` endpoint returns 200 as long as the web server is running,
  even if the database connection is down. Use `/ready` to check full
  service health.
```

### Output Templates

- PREFER providing a concrete template over describing the format in prose.
- ALWAYS store short templates inline in `SKILL.md`.
- PREFER storing longer templates or templates only needed in certain cases in `assets/` and referencing them conditionally.

### Checklists

- PREFER using explicit checklists for multi-step workflows with dependencies or validation gates.

### Validation Loops

- PREFER instructing the agent to validate its own work before proceeding: do the work → run validation → fix issues → repeat until pass.

### Plan-Validate-Execute

- ALWAYS use this pattern for batch or destructive operations: create structured plan → validate against source of truth → execute.

## Script Design

- ALWAYS make scripts non-interactive — accept all input via command-line flags, environment variables, or stdin.
- NEVER use interactive prompts (TTY prompts, password dialogs, confirmation menus). Scripts blocking on interactive input will hang indefinitely.
- ALWAYS provide `--help` output with a brief description, available flags, and usage examples.
- ALWAYS write helpful error messages that say what went wrong, what was expected, and what to try.
- PREFER structured output formats (JSON, CSV, TSV) over free-form text.
- ALWAYS send structured data to stdout and diagnostics (progress, warnings) to stderr.
- PREFER idempotent operations — "create if not exists" over "create and fail on duplicate."
- ALWAYS support a `--dry-run` flag for destructive or stateful operations.
- ALWAYS use meaningful, distinct exit codes for different failure types and document them in `--help`.
- ALWAYS pin dependency versions for reproducibility.
- WHEN complex commands are hard to get right on first try, ALWAYS move them into tested scripts in `scripts/`.

### One-Off Commands

- ALWAYS pin versions (e.g., `npx eslint@9.0.0`).
- ALWAYS state prerequisites in `SKILL.md`.
- PREFER `uvx` or `pipx` for Python packages, `npx` for Node.js packages.

### Self-Contained Scripts

- PREFER scripts that declare their own dependencies inline (PEP 723 for Python, `npm:` specifiers for Deno, `bundler/inline` for Ruby).
- ALWAYS design script output to be agent-consumable — the agent reads stdout and stderr to decide next actions.

## Description Optimization

- PREFER testing descriptions with eval queries (8-10 should-trigger and 8-10 should-not-trigger prompts).
- ALWAYS split queries into train (~60%) and validation (~40%) sets to avoid overfitting.
- PREFER running each query multiple times (at least 3) to account for nondeterministic model behavior.
- ALWAYS iterate on the description using the optimization loop: evaluate → identify failures → revise → repeat.
- NEVER add specific keywords from failed queries — find the general category and address that.

## Validation

- ALWAYS validate skills using the [skills-ref](https://github.com/agentskills/agentskills/tree/main/skills-ref) reference library before distributing.