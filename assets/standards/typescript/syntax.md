---
uri: standards://typescript/syntax
name: TypeScript Syntax Standards
description: "Standards for how TypeScript service and CLI code should look. Control flow, async, generics, and naming."
languages:
    - typescript
file_types:
    - "*.ts"
priority: required
related_resources:
    - standards://typescript/architecture
    - standards://typescript/tooling
---

# TypeScript Syntax Standards
Standards for how TypeScript service and CLI code should look. Control flow, async, generics, and naming.

- ALWAYS follow the repository when it already standardizes on different syntax. Repository rules win.
- ALWAYS apply these rules to new service and CLI code. NEVER restyle existing files unless asked.
- ALWAYS favor legibility. NEVER write clever one-liners.

## Control flow

- NEVER use a ternary (`cond ? a : b`). ALWAYS use `if` / `else`.
- `?.` and `??` are allowed for null or undefined access and defaults.
- `&&` and `||` are conditions, not value expressions.

```typescript
const name = user?.profile?.name ?? "anonymous";

let label: string;
if (user) {
  label = user.name;
} else {
  label = "anonymous";
}
```

## Generics

- Generics are for containers and results: `Array`, `Promise`, `Map`, `Set`, and `Result<T>`.
- NEVER author generic services, generic repositories, generic helpers, or generic frameworks.
- ALWAYS write the concrete type.

```typescript
async function loadUsers(): Promise<readonly User[]> {
  return [];
}

// ❌ BAD — a generic Repository<T> is the wrong shape
```

## Async

- ALWAYS make methods and functions `async` and return `Promise`, including pure logic, so call sites always `await`.
- NEVER use `.then()` or `.catch()` chains.
- NEVER use callbacks.
- Constructors cannot be async. A builder's `build()` can.

## Names and exports

- ALWAYS use named exports only. NEVER use default exports.
- ALWAYS use PascalCase filenames that match the class, except oclif command paths. See `standards://typescript/cli`.
- ALWAYS use relative `.js` extensions on imports in ESM projects. See `standards://typescript/tooling`.
- ALWAYS keep names globally distinct and include the role.
