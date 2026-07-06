---
uri: standards://git/commit-messages
name: Git Commit Message Standards
description: Standards around commit messages, formats, contents, and verbosity
languages: 
    - all
file_types:
    - "*.*"
priority: required
related_resources:
    - standards://git/commit-staging
    - standards://git/operations
---

# Commit Message Standards
Standards around commit messages, formats, contents, and verbosity


- ALWAYS follow [Conventional Commits specification](https://www.conventionalcommits.org/en/v1.0.0/#specification):

Example:
```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```
Common types: `feat`, `fix`, `chore`, `refactor`, `docs`, `test`, `ci`, `revert`.

- ALWAYS explain *what happened and why*, not just *what changed*. ALWAYS include context for:
    - Iterative attempts and course corrections
    - Reverts and the reason for them
    - Partial progress checkpoints

Examples:

```
# ✅ GOOD
fix: revert accidental deletion of retry logic

The previous commit removed the retry block while refactoring error handling.
This reverts that removal — retry logic is still required for transient network errors.

# ❌ BAD
fix stuff
```