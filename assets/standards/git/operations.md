---
uri: standards://git/operations
name: Git Operations Standards
description: "Use when pushing, pulling, rebasing, or rewriting Git history. Does not cover branch names or pull request text."
languages:
    - all
file_types:
    - "*.*"
priority: required
related_resources:
    - uri: standards://git/commit-staging
      when: "splitting commits before a history change"
    - uri: standards://git/commit-messages
      when: "rewriting a message during amend or rebase"
    - uri: standards://git/branches
      when: "creating a branch before the history change"
---

# Git Operations Standards
Use when pushing, pulling, rebasing, or rewriting Git history. Does not cover branch names or pull request text.

Read this file. Read a related resource only when its `when` condition matches this step.

- NEVER amend commits. NEVER use the `--amend` flag. IF a commit needs a fix, create a new commit on top.
- NEVER force push: NEVER use `git push --force` or `git push --force-with-lease`.
- AVOID rebasing: AVOID `git rebase`. PREFER `git pull` to integrate upstream changes.
- NEVER edit history: Never use `git rebase -i`, `git reset --hard` to discard commits, or any operation that rewrites existing history.

## Corrections and Reversions

- ALWAYS use `git revert <sha>` to undo a commit. ALWAYS keep the full audit trail
- ALWAYS Apply fixes, typo corrections, and mistake resolutions as new commits. NEVER edit prior commits.
