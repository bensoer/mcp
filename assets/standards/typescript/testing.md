---
uri: standards://typescript/testing
name: TypeScript Testing Standards
description: "Jest standards for TypeScript services and CLIs. Colocated tests, public methods, and hand-written fakes."
languages:
    - typescript
file_types:
    - "*.test.ts"
priority: required
related_resources:
    - standards://typescript/architecture
    - standards://typescript/tooling
---

# TypeScript Testing Standards
Jest standards for TypeScript services and CLIs. Colocated tests, public methods, and hand-written fakes.

- ALWAYS follow the repository when it already uses another test runner. Repository rules win.
- ALWAYS use Jest as the runner on a new project.
- ALWAYS colocate `ClassName.test.ts` next to `ClassName.ts`.
- ALWAYS test through the public method.
- ALWAYS pass a hand-written fake into the constructor or options object.
- NEVER extract, export, or reshape production code so a test can call it.
- NEVER widen `#` field visibility so a test can reach internal state.
- Module mocks are NOT a reason to change production shape.
