---
uri: standards://go/testing
name: Go (Golang) Testing Standards
description: "Use when writing or reviewing Go tests. Covers package choice, fakes, tables, and what make test and make check run."
languages:
    - go
file_types:
    - "*_test.go"
priority: required
related_resources:
    - uri: standards://go/architecture
      when: "a test is pressuring production layout or exports"
    - uri: standards://go/tooling
      when: "changing how tests, lint, or make check run"
---

# Go (Golang) Testing Standards
Use when writing or reviewing Go tests. Covers package choice, fakes, tables, and what make test and make check run.

Read this file. Read a related resource only when its `when` condition matches this step.

- ALWAYS follow the repository when it already uses another test framework. Repository rules win.
- PREFER the standard `testing` package unless the repository already uses another framework.

## Test package

- ALWAYS use same-package tests (`package foo`) by default.
- Use `package foo_test` ONLY for a deliberate public-API test.

## How to write cases

- PREFER direct named tests.
- Use a table and `t.Run` ONLY when the cases are genuinely repetitive.
- ALWAYS use hand-written fakes and stubs.
- Use generated mocks ONLY when repetition justifies them.
- NEVER reshape production code for tests. NEVER extract or export a helper only so a test can call it. See `standards://go/architecture`.

## What must pass

- `make test` runs `go test ./...`.
- `make check` adds race, vet, lint, vuln, build, and CLI `--help` smoke when the project is a CLI.
- WHEN the task is a bug fix, a fix without a regression test is incomplete.
- NEVER weaken production design to make a test easier.
