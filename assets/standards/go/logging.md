---
uri: standards://go/logging
name: Go (Golang) Logging Standards
description: "Standards, settings and practices for setting up and configuring logging in any go project."
languages:
    - go
file_types:
    - "*.go"
priority: required
related_resources:
    - standards://go/syntax
    - standards://go/architecture
    - standards://go/cli
---

# Go (Golang) Logging Standards
Standards, settings and practices for setting up and configuring logging in any go project.

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
