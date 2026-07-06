---
uri: standards://git/commit-staging
name: Git Commit Staging Standards
description: Standards on git commits. How they should be created, separated, commented and organised
languages: 
    - all
file_types:
    - "*.*"
priority: required
relatedResources:
    - standards://git/commit-messages
---

# Git Commit Staging Standards
Standards on git commits. How they should be created, separated, commented and organised

- NEVER exceed 500 lines per commit. ALWAYS keep total insertions + deletions at or under 500 lines per commit. Split changes as necessary.
- PREFER splitting commits into logical groupings of changes

- ALWAYS commit tests seperatly. NEVER mix test changes with source or doc changes in the same commit
    - ALL files matching the name `*.test.*`, `*.spec.*`, `*_test.*` are ALWAYS considered test files
    - ALL files under the folders `tests/`, `__tests__/` or similar are ALWAYS considered test files
- PREFER batching multiple test files within the same commit. ALWAYS follow 500 line change limit for tests.

- ALWAYS commit docs seperatly. NEVER mix doc changes with source or test changes in the same commit
    - ALL files within the `docs/` or `doc/` folder is ALWAYS considered a doc file
    - ALL file names matching `README*`, `AGENTS*` are ALWAYS considered doc files
    - ALL diagram files such as `.mmd`, `.drawio`, `.puml` are ALWAYS considered doc files
    - ALL image files such as `.png`, `.jpg`, `.jpeg`, `.svg`, `.gif`, `.webp` are ALWAYS considered doc files
- PREFER batching multiple doc files within the same commit. ALWAYS follow 500 line change limit for docs.


See also [git-workflow](./git-workflow.mdc) for branch strategy and commit message conventions.
