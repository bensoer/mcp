---
uri: standards://go/syntax
name: Go (Golang) Syntax Standards
description: "Standards for how go (golang) code should look. Naming conventions, language structure preferences"
languages:
    - go
file_types:
    - "*.go"
priority: required
related_resources:
    - standards://go/architecture
    - standards://go/logging
    - standards://comments
---

# Go (Golang) Syntax Standards
Standards for how go (golang) code should look. Naming conventions, language structure preferences

- ALWAYS follow the repository when it already standardizes on different tools. Repository rules win.
- ALWAYS favor legibility over performance or speed. ALWAYS choose common, popular, and well-documented implementation strategies.
- NEVER write novel or unique implementations.
- NEVER write clever one-liners. NEVER write ternary-style expressions.

## Function layout

- ALWAYS place the function body below the function signature.
- NEVER place the function body on the same line as the signature.

```go
// ✅ GOOD — body is below the signature
func SayHello() {
    fmt.Println("Hello")
}

// ❌ BAD — body is on the same line as the signature
func SayHello() { fmt.Println("Hello") }
```

## File and package naming

- ALWAYS name `.go` files and folders in snake_case.
- ALWAYS group folders and packages by feature. The package name is the feature. File names describe the type or role inside that package.
- ALWAYS place DTO, request, and response structs in `types.go` for the package.
- ALWAYS place interfaces in `protocols.go` for the package.
- ALWAYS place sentinel and typed errors in `errors.go` for the package.
- ALWAYS place internal Go code inside `internal/`.
- ALWAYS place application entrypoints (`main.go`, CLI or API wiring) inside `cmd/`.
- ALWAYS place compiled output inside `bin/`.
- ALWAYS place vendored 3rd-party library source inside `lib/`.
- ALWAYS place vendored 3rd-party binaries inside `dist/`.
- PREFER not to alias types. AVOID adding indirection.
- ALWAYS keep names globally explicit and include the role, even if the package name repeats a word. Example: `UserRepository`, `NotificationFacade`.
- See `standards://go/architecture` for one-struct-per-file and constructor rules.

## Generics

- NEVER use generics.
- IF a design seems to need generics, the architecture is wrong. Use concrete types instead.

## Context

- ALWAYS make `context.Context` the first parameter.
- ALWAYS propagate `context.Context` through cancellable and I/O calls.
- NEVER store `context.Context` on a struct.

## Errors

- ALWAYS wrap errors with context: `fmt.Errorf("load user %s: %w", id, err)`.
- ALWAYS return the error.
- ALWAYS log an error once, at the application boundary.
- NEVER log and return the same error.
- Libraries NEVER panic. Libraries NEVER call `Fatal`.
- See `standards://go/cli` for who owns process exit.

## Data and concurrency

- PREFER explicit structs for requests, responses, config, and JSON.
- Use `map[string]interface{}` ONLY when the schema is truly open-ended.
- ALWAYS stay synchronous unless concurrency is required.
- WHEN concurrency is required, ALWAYS bound it, ALWAYS honor context cancellation, and NEVER leak goroutines.

## Comments

- ALWAYS document every exported type, function, and method.
- Internal comments explain non-obvious reasons only.
- NEVER delete or reword human comments. Follow `standards://comments`.
