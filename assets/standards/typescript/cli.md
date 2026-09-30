---
uri: standards://typescript/cli
name: TypeScript CLI Standards
description: "oclif standards for TypeScript CLIs. Thin commands, help text, flag validation, and where domain types live."
languages:
    - typescript
file_types:
    - "*.ts"
priority: required
related_resources:
    - standards://typescript/architecture
    - standards://typescript/syntax
    - workflows://typescript-oclif-subcommand
    - workflows://typescript-new-project
---

# TypeScript CLI Standards
oclif standards for TypeScript CLIs. Thin commands, help text, flag validation, and where domain types live.

- ALWAYS follow the repository when it already uses another CLI stack. Repository rules win.
- ALWAYS use oclif for a new CLI.
- ALWAYS generate a new command with `oclif generate command`. NEVER hand-write the skeleton. See `workflows://typescript-oclif-subcommand`.

## Command edge

- A command file parses flags, validates them, loads a service from the composition root, calls one method, and prints.
- NEVER put business logic in a command file.
- ALWAYS put request and response classes, clients, and domain methods in `src/<feature>/`, one class per file, PascalCase filename, re-exported from that folder's `index.ts`.
- NEVER put domain classes in `src/commands/`.
- Under `src/commands/`, the path is the command id and stays kebab-case. Example: `postgres/enable-subscription.ts` is `postgres:enable-subscription`.
- The class inside is still PascalCase and includes the role. Example: `EnableSubscriptionCommand`.
- ALWAYS use a named export. IF a generator emitted `export default`, change it to a named export.
- ALWAYS put one command class in the file.
- `run()` is `async`.
- ALWAYS ask the composition root or container for the service.
- NEVER construct backends inside the command.
- NEVER register services from the command file.
- Domain classes return `Result<T>` for expected failures and throw for unexpected ones.
- Domain classes NEVER call `this.error`, `this.exit`, or `process.exit`.
- The command prints.
- Machine-readable JSON uses oclif's `enableJsonFlag` and a returned value.
- Diagnostics stay on the command's warn and error path.

## Help text

- ALWAYS fill every help field.

| Field | Rule |
|-------|------|
| `summary` | One complete sentence. Shown in the parent topic list. |
| `description` | At least 2–3 sentences. What it does, when to use it, and caveats. |
| `examples` | At least one realistic invocation using `<%= config.bin %>`. |
| `args` | Declare positional args. Set `required: true` when the command cannot run without them. |
| `flags` | Declare every flag. Set `required: true` when the command cannot run without it. |

## Validate before I/O

- AFTER `this.parse`, ALWAYS validate flags, args, and combinations BEFORE any network, file, or database call.
- IF flags are illegal or contradictory, call `this.error(...)` and stop.
- IF a flag is a short-circuit, such as `--dry-run` or `--external-only`, do only that branch, then return.
- NEVER run the work a short-circuit flag skips.
