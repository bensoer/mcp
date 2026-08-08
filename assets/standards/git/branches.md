---
uri: standards://git/branches
name: Git Branching Standards
description: Rules on branch creation, branch naming conventions, pushing and PR practices
languages: 
    - all
file_types:
    - "*.*"
priority: required
related_resources:
    - standards://git/commit-messages
    - standards://git/commit-staging
---

# Branching Standards
Rules on branch creation, branch naming conventions, pushing and PR practices

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