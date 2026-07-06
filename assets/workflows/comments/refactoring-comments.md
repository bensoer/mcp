---
uri: workflows://refactoring-comments
name: Refactoring Comments Workflow
description: Workflow instructions for refactoring code that contains comments. Instructions and rules to follow on how to effectively preserve comments through code changes or refactors - whether that planned or actual.
related_resources:
    - standards://comments
---

# Step 1: Review All Code For Comments

## Pre-edit Comment Checks

**Before writing any `new_string` for a StrReplace (or equivalent) call, scan the `old_string` region for every comment — inline, block, or docstring.** Enumerate them. Verify each one appears verbatim in `new_string` before submitting the edit. Do not skip this step even for small or "obvious" changes.

```python
# ❌ BAD — old_string contains comments; new_string silently drops them
old: |
    # Replace underscores with hyphens
    cleaned = re.sub(r'[_]', '-', name)
new: |
    cleaned = re.sub(r'[_]', '-', name.lower())

# ✅ GOOD — comments carried forward verbatim
old: |
    # Replace underscores with hyphens
    cleaned = re.sub(r'[_]', '-', name)
new: |
    cleaned = name.lower()
    # Replace underscores with hyphens
    cleaned = re.sub(r'[_]', '-', cleaned)
```

## Pre-delete Comment Checks

Before removing any block of code, check whether it has comments, inline comments, or docstrings attached. If it does, **do not delete** — use a targeted in-place edit to modify the code, keeping the comments intact. Only proceed with deletion if the block has no comments of any kind attached to it.

# Step 2: Apply Changes

- NEVER use delete-then-rewrite on existing files
- ALWAYS use targets in-place edits (StrReplace or equivalent patch operations) on existing files
- NEVER delete a file or method and then rewrite it from scratch


IF a change is so large that targeted edits are genuinely impractical THEN:
1. Read the entire file and explicitly enumerate every comment, inline comment, and docstring it contains.
2. Carry all of them forward verbatim into the new version.
3. Do not proceed with the rewrite until step 1 is complete.

# Step 3: Post Change Review

## Factually stale comments after refactoring

If a comment appears to describe something the code no longer does after a refactor (e.g. a numbered step whose step was removed), do not silently leave a misleading comment and do not silently delete it. Flag it to the user with a note explaining why the comment may now be inaccurate, and ask how to handle it before proceeding.
