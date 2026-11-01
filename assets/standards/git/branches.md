---
uri: standards://git/branches
name: Git Branching Standards
description: "Use when creating or naming a Git branch. Does not cover pushing, history rewrites, or pull requests."
languages:
    - all
file_types:
    - "*.*"
priority: required
related_resources:
    - uri: standards://git/commit-messages
      when: "writing the commit message on that branch"
    - uri: standards://git/commit-staging
      when: "deciding how to split commits on that branch"
---

# Branching Standards
Use when creating or naming a Git branch. Does not cover pushing, history rewrites, or pull requests.

Read this file. Read a related resource only when its `when` condition matches this step.

- ALWAYS branch off from main/base branch before making any changes
- ALWAYS pull latest remote changes from main/base branch before branching

- ALWAYS name branches after ticket numbers. IF a ticket number is not provided, ask the user for the ticket number or the name to use

Examples:
```
username/eng-123-short-description - GOOD
vague-branch-name - BAD
eng-123 - GOOD
username/eng-123 - NOT IDEAL: Shorten to eng-123 instead
```
