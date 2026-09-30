---
uri: workflows://typescript-oclif-subcommand
name: TypeScript oclif Subcommand Workflow
description: "Step-by-step workflow for adding an oclif command with oclif generate, kebab-case command paths, full help text, and flag validation before I/O."
languages:
    - typescript
file_types:
    - "*.ts"
priority: required
related_resources:
    - standards://typescript/cli
    - standards://typescript/architecture
    - workflows://typescript-new-project
---

# TypeScript oclif Subcommand Workflow
Step-by-step workflow for adding an oclif command.

Use this when adding a command or subcommand to an existing oclif CLI. Follow `standards://typescript/cli` and `standards://typescript/architecture`.

- ALWAYS keep domain types outside `src/commands/`.
- A command file parses flags, validates them, loads a service from the composition root, calls one method, and prints.
- NEVER put business logic in the command file.

## Step 1: Read the command tree

1. Read `package.json` `oclif` config for the commands directory and bin name.
2. List `src/commands/` and match how topics are nested.
3. IF the parent topic is unclear, ask before generating.

## Step 2: Generate with oclif

- ALWAYS generate the skeleton. NEVER hand-write the file from scratch.

```bash
oclif generate command <topic>:<name>
```

```bash
oclif generate command postgres:enable-subscription
oclif generate command version
```

- ALWAYS keep the path the generator chooses.
- `src/commands/postgres/enable-subscription.ts` is the command id `postgres:enable-subscription`.
- IF `oclif` is not installed, use the project's local binary (`pnpm exec oclif`) or `pnpm add -D oclif` in a CLI that is already an oclif project.

## Step 3: Shape the class

- ALWAYS keep one command class in the file.
- ALWAYS use a named export. IF the generator emitted `export default`, change it to a named export.
- Class name is PascalCase and includes the role. Example: `EnableSubscriptionCommand`.
- The filename stays kebab-case. That is the oclif exception to PascalCase filenames.
- `run()` is `async`.

## Step 4: Fill in help

- ALWAYS fill `summary`, `description`, `examples`, `args`, and `flags`. Follow the field rules in `standards://typescript/cli`.
- `summary` is one complete sentence.
- `description` is at least 2–3 sentences.
- `examples` includes at least one realistic invocation using `<%= config.bin %>`.
- Set `required: true` on args and flags the command cannot run without.

## Step 5: Validate before I/O

- AFTER `this.parse`, ALWAYS validate flags, args, and combinations before any network, file, or database call.
- IF flags are illegal or contradictory, call `this.error(...)` and stop.
- IF a flag short-circuits the work, such as `--dry-run` or `--external-only`, do only that branch, then return.
- NEVER run the work a short-circuit flag skips.

## Step 6: Keep the command thin

- ALWAYS ask the composition root or container for the service.
- NEVER construct backends inside the command.
- NEVER register services from the command file.
- ALWAYS put request and response classes, clients, and domain methods in `src/<feature>/`, one class per file, PascalCase filename, re-exported from that folder's `index.ts`.
- Domain classes return `Result<T>` for expected failures and throw for unexpected ones.
- Domain classes NEVER call `this.error`, `this.exit`, or `process.exit`.
- The command prints.
- Machine-readable JSON uses oclif's `enableJsonFlag` and a returned value.
- Diagnostics stay on the command's warn and error path.

## Checklist

- [ ] Generated with `oclif generate command`, not hand-written
- [ ] Path under `src/commands/` is the kebab-case command id
- [ ] Named export, PascalCase class, one class in the file
- [ ] `summary`, `description` (2–3 sentences), and `examples` are filled in
- [ ] Args and required flags are declared
- [ ] `run` validates before any I/O
- [ ] Short-circuit flags return before the work they skip
- [ ] Service comes from the container; domain types live outside `src/commands/`
