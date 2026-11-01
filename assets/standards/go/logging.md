---
uri: standards://go/logging
name: Go (Golang) Logging Standards
description: "Use when adding or configuring logging in a Go project. Covers logger setup, levels, and where an error is logged."
languages:
    - go
file_types:
    - "*.go"
priority: required
related_resources:
    - uri: standards://go/syntax
      when: "wrapping or returning the error being logged"
    - uri: standards://go/architecture
      when: "placing the logger in the container or a feature package"
    - uri: standards://go/cli
      when: "the CLI owns process exit or stdout for that log"
---

# Go (Golang) Logging Standards
Use when adding or configuring logging in a Go project. Covers logger setup, levels, and where an error is logged.

Read this file. Read a related resource only when its `when` condition matches this step.

- ALWAYS follow the repository when it already standardizes on a different logger. Repository rules win.
- PREFER `go.uber.org/zap` sugared logger for new logging.
- Use Go `slog` ONLY when explicitly asked, or when it is absolutely necessary.
- NEVER add `uber/zap` if the codebase already uses another logger.

## Where the logger lives

- ALWAYS carry a contextual logger on `context.Context`.
- NEVER store `context.Context` on a struct.
- Use `zap.S()` ONLY when there is no request or command context.
- `cmd/` glue MAY use `zap.S()` until it has a context.
- `LoggerFromContext` falls back to `zap.S()` when the context has no logger.
- ALWAYS put the base logger on `context.Context` at the command edge.
- Constructors receive `*zap.SugaredLogger` as an explicit parameter when the struct needs a logger. NEVER read a global logger from inside a method. See `standards://go/architecture`.

## What gets logged

- ALWAYS log an error once, at the application boundary.
- NEVER log and return the same error.
- Libraries NEVER call `log.Fatal`. The command layer owns process exit. See `standards://go/cli`.
