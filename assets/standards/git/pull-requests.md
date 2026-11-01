---
uri: standards://git/pull-requests
name: Pull Request Standards
description: "Use when writing or reviewing a pull request title or description. Covers Conventional Commit titles and the required description sections."
languages:
    - all
file_types:
    - "*.*"
priority: required
related_resources:
    - uri: standards://git/commit-messages
      when: "the title must follow commit-message syntax, or a commit message is being written"
    - uri: standards://git/commit-staging
      when: "commits in the pull request are being split or regrouped"
    - uri: standards://git/branches
      when: "the pull request branch is being created or renamed"
    - uri: standards://git/operations
      when: "the branch is being pushed or its history rewritten"
---

# Pull Request Standards
Use when writing or reviewing a pull request title or description. Covers Conventional Commit titles and the required description sections.

Read this file. Read a related resource only when its `when` condition matches this step.

## Title

- ALWAYS follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/#specification) syntax and structure for the pull request title
- ALWAYS ensure the title reflects ALL commits contained within the pull request. If additional commits are added, ALWAYS review and update the title if necessary
- ALWAYS keep the title to a single concise line (subject line only — no body)

Syntax:
```
<type>[, <type>...] [optional scope(s)]: <description>
```
Common types: `feat`, `fix`, `refactor`, `docs`, `test`, `ci`, `chore`, `revert`.

When commits span multiple types:
- ALWAYS include every type present in the pull request as a comma-separated prefix
- List types in priority order: `feat`, `fix`, `refactor`, `perf`, `build`, `ci`, `docs`, `test`, `chore`, `revert`

Examples:
```
# ✅ GOOD — every type present, comma-separated, updated when later commits change the set
feat, fix, refactor: add password reset flow and harden lookups

# ❌ BAD — omits types present in the pull request
feat: add password reset modal
```

## Description

- ALWAYS include every section below, in this order: `Summary`, `Added`, `Removed`, `Updated`, `Deleted`, `QA and testing`, `Operational risk`
- ALWAYS include `Added`, `Removed`, `Updated`, and `Deleted` even when a section has no changes. Mark an empty change section `- None`
- NEVER omit `Summary`, `QA and testing`, or `Operational risk`. IF nothing applies, say so in that section

### Summary

- ALWAYS describe what was created, the intended outcome, what changed, and why
- Write it as a product description, not as a file list
- IF the conversation provides plan files or tickets, ALWAYS use them as the source for the intended outcome and the why
- ALWAYS link each known plan file or ticket. Use the path or URL that was provided. Do not assume a repository layout
- NEVER invent a plan, a ticket, or a link
- IF none were provided, say that, and write the summary from the commits and the change itself

### QA and testing

- ALWAYS summarize what was tested, how it was tested, and what was created or done to verify the changes work as expected
- NEVER claim a check ran if it did not
- IF nothing was tested, say what was not verified

### Operational risk

- ALWAYS summarize the impact on the application and its users
- ALWAYS say whether the changes are breaking
- ALWAYS say what should be verified before deploy, and any compatibility concern
- IF there is no breaking change, no pre-deploy check, or no compatibility concern, say so. Do not leave that question unanswered

Use this template. Do not add a second copy of these sections.

```
### Summary
[What was created, the intended outcome, what changed, and why. Product description, not a file list. Link each plan or ticket that was provided. If none were provided, say so. Do not invent a link.]

### Added
- [change, or None]

### Removed
- [change, or None]

### Updated
- [change, or None]

### Deleted
- [change, or None]

### QA and testing
- [What was tested, how, and what was not verified. Do not claim a check that did not run.]

### Operational risk
- [User impact, whether it is breaking, what to verify before deploy, and compatibility. If none, say so.]
```
