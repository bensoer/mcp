---
uri: standards://git/pull-requests
name: Pull Request Standards
description: Standards for pull request titles, descriptions, and review practices. Consult when creating, updating, reviewing and maintaining pull requests.
languages: 
    - all
file_types:
    - "*.*"
priority: required
related_resources:
    - standards://git/commit-messages
    - standards://git/commit-staging
    - standards://git/branches
    - standards://git/operations
---

# Pull Request Standards
Standards for pull request titles, descriptions, and review practices. Refer to these standards when creating, updating, reviewing and maintaining a pull request.

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
# ✅ GOOD - single type, reflects all commits
feat: add password reset flow

# ✅ GOOD - multiple types listed, comma-separated
feat, test: add password reset flow and unit tests

# ✅ GOOD - title updated when commits changed
feat, fix, refactor: add password reset flow and harden lookups
# (originally just "feat, test: add password reset", but later commits added
#  a bug fix and a refactor — title was reviewed and broadened accordingly)

# ❌ BAD - not conventional commit format
Added some stuff

# ❌ BAD - omits types present in the PR
feat: add password reset modal
# (PR also contains a `fix` and a `refactor` commit — title should be
#  "feat, fix, refactor: add password reset modal and harden lookups")
```

## Description

- ALWAYS include a bullet list in the description covering CRUD-style changes
- ALWAYS include all four sections — `Added`, `Removed`, `Updated`, `Deleted` — even if a section has no changes (mark it `- None`)

Structure:
```
### Added
- ...

### Removed
- ...

### Updated
- ...

### Deleted
- ...
```

Example:
```
### Added
- New `POST /api/reset-password` endpoint
- Password reset email template
- `ResetPasswordRequest` DTO with validation

### Removed
- Legacy `sendResetEmail` utility function

### Updated
- Updated `AuthController` to route password reset requests to the new endpoint
- Updated user guide with password reset instructions

### Deleted
- `docs/legacy/reset-flow.md`
```
