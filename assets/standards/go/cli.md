---
uri: standards://go/cli
name: Go (Golang) CLI Standards
description: "Cobra and Viper standards for Go CLIs. Thin command edge, help text, flag validation, and who owns stdout and process exit."
languages:
    - go
file_types:
    - "*.go"
priority: required
related_resources:
    - standards://go/architecture
    - standards://go/logging
    - workflows://go-cobra-subcommand
    - workflows://go-new-project
---

# Go (Golang) CLI Standards
Cobra and Viper standards for Go CLIs. Thin command edge, help text, flag validation, and who owns stdout and process exit.

- ALWAYS follow the repository when it already uses another CLI stack. Repository rules win.
- ALWAYS use Cobra for a new CLI unless the project already uses another stack.
- ALWAYS use Viper when the project uses Cobra, so flags, env, and config files resolve in one place.
- ALWAYS scaffold a new command with `cobra-cli`. See `workflows://go-cobra-subcommand`.
- NEVER use `cobra-cli init`. It emits a root `main.go` plus `cmd/root.go`, which is not this layout. See `workflows://go-new-project`.

## Command edge

- `cmd/` is a thin edge: flags, validation, container lookup, and printing.
- ALWAYS keep business logic in `internal/`.
- IF a type is not a `cobra.Command`, or a flag or output helper used only by that command, move it out of `cmd/`.
- Command code parses flags, validates them BEFORE any I/O, asks the container for a service, and prints results.
- `Run` or `RunE` validates flags, loads the service from the container, calls one method, and prints the result.
- NEVER construct clients, dial networks, or decode feature config inside Cobra `init()`.
- Cobra's own command-registration `init()` is allowed. NEVER add any other `init()`. See `standards://go/architecture`.
- IF the parent command already attaches children in an explicit constructor, register the new command there. NEVER introduce a second registration style.

## Process output

- Libraries NEVER `fmt.Print`.
- Libraries NEVER `os.Exit`.
- Libraries NEVER `log.Fatal`.
- The command layer owns stdout, stderr, and process exit.
- ALWAYS send machine output to stdout.
- ALWAYS send diagnostics to stderr.

## Help text

- ALWAYS populate every help field. Sparse `Short` values and empty `Long` fields produce unhelpful `--help` output.

| Field | Rule |
|-------|------|
| `Use` | Command name and argument signature. Example: `"enable-subscription <subscriptionId>"`. |
| `Short` | One complete, descriptive sentence. Shown in the parent's command list. |
| `Long` | At least 2–3 sentences. What it does, when to use it, and important caveats. |
| `Example` | At least one realistic invocation. Use the full binary name. |
| `Args` | ALWAYS set a validator when the command takes positional arguments. Example: `cobra.ExactArgs(1)`. NEVER leave `Args` unset in that case. |

```go
// ✅ GOOD — complete help
var enableSubscriptionCmd = &cobra.Command{
    Use:   "enable-subscription <subscriptionId>",
    Short: "Re-enable a previously disabled logical replication subscription.",
    Long: `Re-enables a logical replication subscription that was previously disabled.
The subscription resumes replication from the last confirmed LSN position.

Requires a database connection. Pass --config or environment variables.`,
    Example: `  stroom postgres enable-subscription sub_keela_main
  stroom postgres enable-subscription sub_keela_main --config ~/.stroom/prod.yaml`,
    Args: cobra.ExactArgs(1),
}

// ❌ BAD — empty or vague help
var enableSubscriptionCmd = &cobra.Command{
    Use:   "enableSubscription",
    Short: "enable subscription",
    Long:  ``,
}
```

## Flags

- ALWAYS declare flags as typed variables at the top of the file, before the command var, matching the rest of `cmd/`.
- ALWAYS bind flags in `init()`, not in `Run`.
- ALWAYS mark a flag required with `cmd.MarkFlagRequired("flagName")` when the command cannot run without it.
- Use `PersistentFlags()` on the parent ONLY if the flag applies to all subcommands of that parent.
- Use `Flags()` for flags local to this command.
- ALWAYS name new commands verb + noun, camelCase, matching the project. Example: `enableSubscription`.

## Validate before I/O

- ALWAYS validate every flag, argument, and combination at the top of `Run` or `RunE`, before any API call, network request, file I/O, or other side effect.
- IF input is illegal, or flags contradict each other, log the problem and `Fatalf` or return an error immediately. NEVER proceed.
- IF a flag means only a subset of work is needed, evaluate that flag first and return as soon as that work is done.
- NEVER execute code a short-circuit flag says not to run.
- NEVER run unrelated work before flag evaluation. Eager execution pollutes stdout, which breaks shell composition.

```go
// ✅ GOOD — flags checked first
if flagA && flagB {
    log.Fatalf("--flag-a and --flag-b are mutually exclusive")
}
if flagExternalOnly {
    fmt.Fprintln(os.Stdout, fetchExternal())
    return
}

// ❌ BAD — full work runs even when --external-only is set
ext := fetchExternal()
internal := fetchInternal()
if flagExternalOnly {
    fmt.Fprintln(os.Stdout, ext)
    return
}
```
