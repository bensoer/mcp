---
uri: standards://git/commit-messages
name: Git Commit Message Standards
description: "Use when writing, amending, or reviewing a commit message. Covers Conventional Commits, the body, and verbosity."
languages:
    - all
file_types:
    - "*.*"
priority: required
related_resources:
    - uri: standards://git/commit-staging
      when: "splitting or grouping the changes the message describes"
    - uri: standards://git/operations
      when: "amending, rebasing, or pushing that commit"
---

# Commit Message Standards
Use when writing, amending, or reviewing a commit message. Covers Conventional Commits, the body, and verbosity.

Read this file. Read a related resource only when its `when` condition matches this step.


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
