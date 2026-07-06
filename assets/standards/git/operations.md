---
uri: standards://git/operations
name: Git Operations Standards
description: Operational standards for working with Git in regards to managing history, pushing and pulling changes
languages: 
    - all
file_types:
    - "*.*"
priority: required
related_resources:
    - standards://git/commit-staging
    - standards://git/commit-messages
    - standards://git/branches
---

# Git Operations Standards
Operational standards for working with Git in regards to managing history, pushing and pulling changes

- NEVER amend commits. NEVER use the `--amend` flag. IF a commit needs a fix, create a new commit on top.
- NEVER force push: NEVER use `git push --force` or `git push --force-with-lease`.
- AVOID rebasing: AVOID `git rebase`. PREFER `git pull` to integrate upstream changes.
- NEVER edit history: Never use `git rebase -i`, `git reset --hard` to discard commits, or any operation that rewrites existing history.

## Corrections and Reversions

- ALWAYS use `git revert <sha>` to undo a commit. ALWAYS keep the full audit trail
- ALWAYS Apply fixes, typo corrections, and mistake resolutions as new commits. NEVER edit prior commits.