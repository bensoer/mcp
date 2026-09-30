---
uri: workflows://go-cobra-subcommand
name: Go Cobra Subcommand Workflow
description: "Step-by-step workflow for adding a Cobra subcommand with cobra-cli, complete help text, and registration under the correct parent."
languages:
    - go
file_types:
    - "*.go"
priority: required
related_resources:
    - standards://go/cli
    - standards://go/architecture
    - workflows://go-new-project
---

# Go Cobra Subcommand Workflow
Step-by-step workflow for adding a Cobra subcommand.

Use this when adding a command or subcommand to an existing Cobra CLI. Follow `standards://go/cli` and `standards://go/architecture`.

- `cmd/` is a thin edge: flags, validation, container lookup, and printing.
- ALWAYS put business types and logic in `internal/<feature>/`, not in `cmd/`.
- IF a type is not a `cobra.Command`, or a flag or output helper used only by that command, move it out of `cmd/`.

## Step 1: Read the command tree

1. Open `cmd/root.go` and note the root command name and any persistent flags already set up.
2. Read `cmd/*.go` and match the naming convention and how existing commands register themselves.
3. Identify the parent command the new command nests under. Example: `rootCmd`, `k8sCmd`, `postgresCmd`.
4. IF the user has not specified a parent, ask before proceeding.

## Step 2: Generate with cobra-cli

- ALWAYS use `cobra-cli add` to generate the initial boilerplate.
- NEVER hand-write the skeleton from scratch.
- Cobra already registers commands from `init()`. That `init()` is allowed.
- NEVER add any other `init()`: no tool registration, no Viper defaults for library config, and no self-registration of services. Those belong on the container, called from the composition entrypoint.
- IF the parent command already attaches children in an explicit constructor, register the new command there. NEVER introduce a second registration style.

```bash
# Top-level command
cobra-cli add <commandName>

# Nested subcommand. --parent takes the parent var name.
cobra-cli add <commandName> --parent <parentCmdVar>
```

- ALWAYS name commands verb + noun, camelCase, matching the project.

```bash
cobra-cli add enableSubscription --parent postgresCmd
cobra-cli add dumpDatabaseSchema --parent postgresCmd
cobra-cli add deleteAbandonedPods --parent k8sCmd
```

- IF `cobra-cli` is not installed: `go install github.com/spf13/cobra-cli@latest`

## Step 3: Strip the copyright boilerplate

The generated file contains a comment block like:

```go
/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
```

- ALWAYS remove that block entirely.
- NEVER replace it with a copyright or author line.

## Step 4: Fill in help

- ALWAYS populate every help field. Follow the field rules in `standards://go/cli`.
- `Use` includes the argument signature.
- `Short` is one complete, descriptive sentence.
- `Long` is at least 2–3 sentences.
- `Example` shows at least one realistic invocation with the full binary name.
- `Args` is set when the command takes positional arguments.

## Step 5: Declare flags

- ALWAYS declare flags as typed variables at the top of the file, before the command var.
- ALWAYS bind flags in `init()`, not in `Run`.
- ALWAYS call `MarkFlagRequired` when the command cannot run without the flag.
- Use `PersistentFlags()` ONLY when the flag applies to every subcommand of that parent.
- Use `Flags()` for flags local to this command.

## Step 6: Validate before any side effect

- ALWAYS validate every flag, argument, and combination at the top of `Run` or `RunE`, before any API call, network request, file I/O, or other side effect.
- IF input is illegal, or flags contradict each other, `Fatalf` or return an error immediately. NEVER proceed.
- IF a flag short-circuits the work, such as `--external-only`, do only that branch and return.
- NEVER run work a short-circuit flag says to skip.

## Step 7: Register under the correct parent

- `cobra-cli` writes `<parentCmd>.AddCommand(<newCmd>)` inside `init()`. Keep that when the rest of `cmd/` already registers that way.
- ALWAYS confirm the parent variable name from the parent's source file before writing.
- NEVER put library construction, client dialing, or feature config decoding into that `init()`.
- `Run` or `RunE` asks the container, or the parent `PreRun`, for the service after flags are valid.

## Step 8: Keep cmd thin

- `Run` or `RunE` validates flags, loads the service from the container, calls one method, and prints the result.
- NEVER add request or response structs, clients, or domain methods under `cmd/`.
- ALWAYS put those types in `internal/<feature>/`.
- Libraries NEVER print or exit. The command owns stdout, stderr, and process exit.
- ALWAYS name the file after the command: `cmd/<verbNoun>.go`.

## Checklist

- [ ] Generated with `cobra-cli add`, not hand-written
- [ ] Copyright block removed
- [ ] `Use` includes the argument signature
- [ ] `Short` is a complete, descriptive sentence
- [ ] `Long` has at least 2–3 sentences
- [ ] `Example` shows at least one realistic invocation
- [ ] `Args` validator set if the command takes positional arguments
- [ ] Flags declared as package-level vars and bound in `init()`
- [ ] `Run` validates all flags and args before any side-effecting code
- [ ] Mutually exclusive flags fail before any work begins
- [ ] Short-circuit flags are checked first; skipped work never runs
- [ ] Registered under the correct parent, using the existing registration style
- [ ] No new `init()` beyond the Cobra registration `cobra-cli` generated
- [ ] `cmd/` holds the command only; domain types live in `internal/<feature>/`
- [ ] File follows one-type-per-file naming: `cmd/<verbNoun>.go`
