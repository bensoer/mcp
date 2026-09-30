---
uri: standards://typescript/architecture
name: TypeScript Coding Architecture Practices
description: "Best practices for TypeScript service and CLI organisation, composition, contracts, and wiring."
languages:
    - typescript
file_types:
    - "*.ts"
priority: required
related_resources:
    - standards://typescript/syntax
    - standards://typescript/cli
    - standards://typescript/tooling
    - standards://typescript/testing
    - workflows://typescript-new-project
---

# TypeScript Coding Architecture Practices
Best practices for TypeScript service and CLI organisation, composition, contracts, and wiring.

- ALWAYS follow the repository when it already standardizes on different tools. Repository rules win.
- ALWAYS apply these rules when writing new TypeScript for a service or CLI.
- NEVER restyle existing files to match these rules unless asked.
- NEVER apply these rules to React or other UI, or to shared libraries, even in a `.ts` file. Ask before applying them there.
- Layer depth is a per-repository decision. NEVER import a foreign multi-layer stack. NEVER invent one.
- Logging is intentionally unspecified here. Follow the repository.

## Structure

- ALWAYS use classes for services, backends, clients, commands, and anything that owns dependencies or a lifecycle.
- Plain functions are for pure transforms with no owner.
- ALWAYS put one class per file.
- ALWAYS name the file PascalCase, matching the class. Example: `UserRepository.ts`.
- ALWAYS keep the role in the name. Example: `UserRepository`, `NotificationFacade`, `InvoiceBuilder`.
- ALWAYS keep names globally distinct within the project.
- ALWAYS use named exports only.
- ALWAYS give each folder an `index.ts` that re-exports its public symbols.
- ALWAYS import the folder, not the leaf file.

```typescript
import { UserRepository } from "./repositories/index.js";
```

- oclif command files are the filename exception. Under `src/commands/`, the path is the command id (`postgres/enable-subscription.ts`). The class inside is still PascalCase (`EnableSubscriptionCommand`) and is a named export.
- NEVER put domain classes in `src/commands/`. See `standards://typescript/cli`.
- Files with only functions, constants, a zod schema, or re-exports are allowed.

## Inheritance and contracts

- PREFER composition over inheritance.
- A pure contract is an `abstract class` with the `Interface` suffix and only abstract methods. Example: `StorageInterface`.
- An abstract class with shared implementation uses the `Abstract` prefix. Example: `AbstractPinger`.
- Use inheritance ONLY when every concrete subclass meaningfully implements every abstract method, the relationship is a genuine "is-a", and no subclass would implement a hook as an empty body.
- An empty abstract implementation means composition. Extract a collaborator instead of a template-method base.
- NEVER exceed two levels of inheritance depth.
- ALWAYS export the contract.
- Callers depend on `StorageInterface` where a second implementation exists.
- Constructors return the concrete class. NEVER return the interface from a constructor.

## Constructors and builders

- Options objects and builder classes are allowed.
- A small constructor MAY take an options object. Example: `constructor(options: UserRepositoryOptions)`.
- A builder is its own class, in its own file. Example: `UserRepositoryBuilder`.
- The builder MAY change its draft. `build()` returns the immutable instance.
- NEVER use functional-option callbacks. No `new Foo(withTimeout(1))`. Use an options object or a builder.

## Wiring

- NEVER write a factory switch that grows a branch per implementation.
- ALWAYS register implementations and look them up by key.
- Registration lives in the composition root. Objects NEVER register themselves.
- Restify and Fastify ALWAYS use a hand-written container class.
- NEVER use a DI framework or injection decorators with Restify or Fastify.
- NestJS modules and providers ARE the composition root. NEVER add a second container.
- Nest decorators are part of NestJS. NEVER add other decorator libraries.
- Unused constructors and unregistered implementations are unfinished. Wire them or remove them.
- A missing or duplicate container key throws. That throw is an unexpected wiring failure.

## Which HTTP stack

| Situation | Stack |
|-----------|-------|
| About 3–5 endpoints, simple persistence, internal tool, speed of writing matters more than long-term maintainability | Restify |
| Larger service | Fastify or NestJS |
| Command-line tool | oclif |
| HTTP client that holds a base URL, auth, or timeouts | axios, behind a class |

- ALWAYS ask which larger stack when the user has not chosen between Fastify and NestJS.
- ALWAYS use axios for HTTP clients. Put the client behind a class when it holds a base URL, auth, or timeouts.
- See `standards://typescript/cli` for oclif command rules.

## Errors

- ALWAYS return `Result<T>` for expected domain outcomes, such as not found or validation rejected.
- NEVER throw for an expected domain outcome.
- ALWAYS throw for unexpected failures: bugs, downed I/O, broken invariants.
- The command or HTTP edge turns a `Result` into an exit status or an HTTP response.
- Domain classes NEVER call `process.exit`.
- `Result<T>` is one class in its own file. It is an allowed generic. See `workflows://typescript-new-project` for the class shape.
- `value()` and `error()` throw ONLY when the caller ignored `isOk()`. That is a bug.
- Domain misses return `Result.fail`.

## Values

- ALWAYS make instances immutable after construction.
- ALWAYS mark public fields `readonly`.
- Methods return a new instance instead of assigning `this`.
- The owning builder MAY mutate its draft until `build()`.
- ALWAYS use a string union when there are 3 or fewer variants and the set will stay that size.
- ALWAYS use a string enum when there are 4 or more variants, or when a 3-variant set is likely to grow.
- `Partial`, `Pick`, and `Omit` are for a boundary DTO, such as a patch body.
- Application code names its own shape. NEVER use conditional types. NEVER use mapped types.
- ALWAYS use `#` private fields for internal state.
- Tests call public methods. NEVER widen visibility so a test can reach `#` state.
- `null` and `undefined` follow the boundary. JSON `null` stays `null`.

## Validation and settings

- ALWAYS use zod to parse input and serialize output at the boundary: HTTP body, CLI payload, env strings.
- The schema maps into a class. Application code uses that class.
- NEVER use `z.infer` as the application model.
- The zod schema MAY sit in the same file as the class it constructs. It is not a second class.
- NEVER use class-validator.
- Settings are their own class, loaded from the environment by a loader.
- Domain code receives `Settings`. Domain code NEVER reads `process.env`.
- NEVER pass a `z.infer` settings type through the program.

## Method extraction

- ALWAYS write logic inline at the single call site first.
- Extract a method when the same logic is used in two or more places, OR when leaving it inline makes that call site hard to read.
- PREFER a method on the owning class over a free function.
- NEVER import a helper from a sibling file to reuse it. Move it to a real home first.
- Interface methods and exported API methods MAY exist when called once. They are contracts.
- NEVER extract, export, or reshape production code so a test can call it.
