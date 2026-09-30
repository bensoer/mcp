---
uri: workflows://typescript-new-project
name: TypeScript New Project Workflow
description: "Step-by-step workflow for scaffolding a new TypeScript service (Restify, Fastify, or NestJS) or an oclif CLI with pnpm, ESM, ESLint, Prettier, and Jest."
languages:
    - typescript
file_types:
    - "*.ts"
    - "package.json"
priority: required
related_resources:
    - standards://typescript/architecture
    - standards://typescript/syntax
    - standards://typescript/tooling
    - standards://typescript/testing
    - standards://typescript/cli
    - workflows://typescript-oclif-subcommand
---

# TypeScript New Project Workflow
Step-by-step workflow for scaffolding a new TypeScript service or CLI.

Use this for a **new** TypeScript service or CLI. Existing repositories keep their layout and tools. Follow `standards://typescript/architecture`, `standards://typescript/syntax`, `standards://typescript/tooling`, `standards://typescript/testing`, and `standards://typescript/cli`.

- ALWAYS confirm the project kind before writing files.
- IF a slash command invoked this workflow, stay in Plan mode until the user approves.
- NEVER scaffold React apps or shared libraries with this workflow.

| Kind | When |
|------|------|
| Restify service | About 3–5 endpoints, simple persistence, internal tool. Speed of writing matters more than long-term maintainability. |
| Fastify or NestJS service | Larger than that. Ask which one if the user has not chosen. |
| oclif CLI | A command-line tool. |

## Step 1: Apply shared defaults

- ALWAYS use pnpm and ESM (`"type": "module"`).
- ALWAYS use npm scripts only. NEVER add a Makefile.

```json
{
  "scripts": {
    "build": "tsc",
    "typecheck": "tsc --noEmit",
    "lint": "eslint .",
    "format": "prettier --write .",
    "test": "jest"
  }
}
```

- ALWAYS set TypeScript `strict: true`, `"module": "Node16"`, and `"moduleResolution": "Node16"`.
- ALWAYS use `.js` extensions on relative imports.
- ALWAYS add ESLint flat config plus Prettier.
- ALWAYS wire Jest through the `ts-jest` ESM preset. Tests are `ClassName.test.ts` beside the class.
- ALWAYS use zod at the boundary. Map into a class. NEVER use `z.infer` as the application model.
- ALWAYS add a `Settings` class that loads the environment. Domain code receives `Settings`. Domain code NEVER reads `process.env`.
- Axios clients are classes when they hold a base URL, auth, or timeouts.
- ALWAYS use one class per file, a PascalCase filename, named exports, and a folder `index.ts` barrel.
- ALWAYS make methods and functions `async`.
- ALWAYS make instances immutable after construction. ALWAYS use `#` fields for internals.
- ALWAYS return `Result<T>` for expected failures. Throw for unexpected failures.
- Restify and Fastify use a hand-written container. Nest uses Nest providers. NEVER add a second container. NEVER add a factory switch.

## Step 2: Scaffold a Restify or Fastify service

Keep this thin. Add a feature folder when a second class appears. NEVER pre-build action or runtime layers.

```
src/main.ts                      # boot only
src/Application.ts               # composition root: settings, container, listen
src/container/Container.ts
src/container/index.ts
src/settings/Settings.ts
src/settings/index.ts
src/result/Result.ts
src/result/index.ts
src/routes/HealthRoute.ts        # one route class per file
src/routes/index.ts
package.json
tsconfig.json
eslint.config.js
.prettierrc
jest.config.js
```

- `main.ts` constructs `Application` and awaits `start()`.
- Route classes NEVER register themselves. `Application` registers them.
- `Container` stores implementations by key and returns the concrete class, or the `*Interface` the caller asked for.
- A missing or duplicate key throws. That throw is an unexpected wiring failure.

## Step 3: Scaffold NestJS, if that is the kind

- ALWAYS use the Nest CLI (`nest new`) for the skeleton, then adapt it to `standards://typescript/architecture`.
- PascalCase filenames matching the class (`UserService.ts`).
- One class per file. Named exports. Folder `index.ts` barrels.
- Nest modules and providers are the only composition root.
- ALWAYS use zod at the boundary, then a class. NEVER use class-validator.
- NEVER add a second hand-rolled container. NEVER add a factory switch.

## Step 4: Scaffold oclif, if that is the kind

- ALWAYS generate with `oclif generate <name>` and choose ESM. Then adapt.
- Command files stay kebab-case because the path is the command id.
- BEFORE adding commands, follow `workflows://typescript-oclif-subcommand`.

```
src/commands/hello.ts            # kebab-case; oclif command id
src/container/Container.ts
src/settings/Settings.ts
src/result/Result.ts
src/<feature>/SomeService.ts     # PascalCase domain classes
bin/run.js
```

- The composition root builds the container once.
- Commands parse flags, validate, load a service from the container, and print.
- Domain classes NEVER print or exit.

## Step 5: Add Result

`Result<T>` is one class. It is generic only as `Result<T>`.

```typescript
export class Result<T> {
  readonly #value: T | undefined;
  readonly #error: Error | undefined;

  private constructor(value: T | undefined, error: Error | undefined) {
    this.#value = value;
    this.#error = error;
  }

  static async ok<T>(value: T): Promise<Result<T>> {
    return new Result(value, undefined);
  }

  static async fail<T>(error: Error): Promise<Result<T>> {
    return new Result<T>(undefined, error);
  }

  async isOk(): Promise<boolean> {
    return this.#error === undefined;
  }

  async value(): Promise<T> {
    if (this.#error !== undefined || this.#value === undefined) {
      throw new Error("Result has no value");
    }
    return this.#value;
  }

  async error(): Promise<Error> {
    if (this.#error === undefined) {
      throw new Error("Result has no error");
    }
    return this.#error;
  }
}
```

- `value()` and `error()` throw ONLY when the caller ignored `isOk()`. That is a bug.
- Domain misses return `Result.fail`.

## Checklist

- [ ] Kind confirmed: Restify, Fastify, NestJS, or oclif
- [ ] pnpm, ESM, npm scripts, no Makefile
- [ ] `strict: true` on the new tsconfig
- [ ] ESLint, Prettier, and Jest installed and wired to scripts
- [ ] Composition root wires implementations; no factory switch
- [ ] zod maps into classes; `Settings` is a class
- [ ] `Result<T>` used for expected failures
- [ ] One class per file, named exports, folder barrels
- [ ] oclif command paths left kebab-case
