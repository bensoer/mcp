---
uri: standards://comments
name: Code Comment Preservation Standards
description: Standards around treatment of comments in code - whether that is human or machine authored
languages: 
    - all
file_types:
    - "*.*"
priority: required
related_resources:
    - workflows://refactoring-comments
---

# Code Comment Preservation Standards
Standards around treatment of comments in code - whether that is human or machine authored

## Scope

These rules apply to all comment forms across languages: `#` (Python, Shell, Ruby), `//` (JS, TS, Go, Rust, Java, C/C++), `/* … */` (JS, TS, Java, C/C++, CSS), `<!-- … -->` (HTML, XML), `--` (SQL, Lua), `%` (LaTeX), docstrings (`"""…"""`, `/** … */`, `///`), and any other inline or block comment syntax. All are treated identically — they are human-authored documentation and must be preserved with the same care as the code they annotate.

These rules apply equally to inline comments on the same line as code and to docstring/block comments on classes and methods

## Standards

- NEVER delete existing code comments for any reason
- NEVER reword existing code comments for any reason
- ALWAYS copy code comments verbatim. NEVER rephrase, shorten, clarify or "improve" the wording, EVEN if there is spelling or grammatical errors

- NEVER change, summarise or remove the following comment markers:
    - `CONTEXT:` — explains the non-obvious historical or systemic reason a piece of code exists
    - `NOTE:` — records a constraint or non-obvious behaviour that future developers must know
    - `TODO:` — tracks outstanding work; removal requires the work to be done, not just the comment deleted
    - `FIXME:` — marks a known defect that has not yet been resolved
    - `HACK:` — marks a deliberate workaround with known technical debt

Examples:

```python
# CONTEXT: AWS SDK resets the session on every call; we cache it here to avoid re-auth
```

```typescript
// NOTE: this endpoint is rate-limited to 10 req/s; do not batch above that threshold
```

```go
/* TODO: replace with structured error once Go 1.23 error groups land */
```

- NEVER remove comments of the form `# 1)`, `// 2)`, `-- 3)`
    - IF code is reordered or split up, you MAY adjust the NUMBER to be correct. NEVER edit or the remove the text that follows the number

Examples:

```
# ✅ Allowed — renumbering when steps move to a different method
# 1) Pull chart    →    # 2) Pull chart  (if a new step 1 was added before it)

# ❌ Not allowed — rewording the step text
# 1) Initialise context and set failure defaults so that exceptions produce a clean state
#    changed to:
# 1) Init context    ← do NOT do this
```


- WHEN refactoring splits a method into multiple methods or classes, ALWAYS move each comment to the new location where the code it annotates lives.
    - NEVER leave comments behind at the old call site 
    - NEVER discard comments because the refactor "reorganised" things

    - IF a comment annotates a block that is now spread across multiple locations, ALWAYS place the comment at the most relevant location and do not duplicate or drop it.

