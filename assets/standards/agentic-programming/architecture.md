---
uri: standards://agentic-programming/architecture
name: Agentic Programming Practices
description: "Use when creating or reviewing a repository's agent-facing layout. Covers .agents/, docs/, AGENTS.md, and what each file is for."
languages:
    - all
file_types:
    - "*"
priority: required
related_resources: []
---

# Agentic Programming Best Practices
Use when creating or reviewing a repository's agent-facing layout. Covers .agents/, docs/, AGENTS.md, and what each file is for.

```
.agents/
├─ rules/
├─ skills/
├─ commands/
docs/
├─ investigations/
├─ plans/
├─ PROJECT_CONTEXT.md
├─ ARCHITECTURE.md
├─ DECISIONS.md
README.md
AGENTS.md
```

| Folder / File       | Type   | Description |
| ------------------- | ------ | ----------- |
| .agents/            | Folder | Contains agent rules, commands and skills. Standard folder structure for pi, cursor, claude, codex, etc |
| docs/               | Folder | Contains documentation about the project. Both for humans and agents |
| docs/investigations | Folder | Folder containing markdown (*.md) files of investigations and/or conversations had with an LLM. These may be the source of future plans for the project. These may be raw transcripts or more formally structured documents. It does not necessarily cover the current state of the project, but rather potential futures or past decisions that were made off it
| docs/plans          | Folder | Folder containing markdown (*.md) file copies of plans. These are raw copies of plans for changes and features implemented. These tend to not be completely authoritative of changes within the repository. Changes often happen before, during and after this plan is applied.
| docs/PROJECT_CONTEXT.md | File | Document containing high level outline of what the project is about. What is the end goal and desired outcome of the entire repo and code |
| docs/DECISIONS.md | File | Bullet list of important decisions, tradeoffs, strategies to the project. These are likely important for agents and humans to know that can not be easily perceived or understood from the code or elsewhere in documentation |
| docs/ARCHITECTURE.md | File | Documentation and mermaid diagrams of how the application currently works. ARCHITECTURE.md is ephemeral and will change as the project evolves and moves towards its PROJECT_CONTEXT.md. This provides a snapshot of how the code is organised today. |
| README.md | File | Getting started/introduction information for humans about the repository. Compilation, testing, local environment setup, high-level about instructions, etc |
| AGENTS.md | File | Index / Starting point file for Agents on how to navigate the project. Should contain links and meta information about the folders and documents in the docs/ folder along with any contextual information on how to work with the repository and its documentation. Intended to be the starting point for an Agent to get a grasp of the project
