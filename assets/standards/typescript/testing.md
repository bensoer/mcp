---
uri: standards://typescript/testing
name: TypeScript Testing Standards
description: "Use when writing or reviewing Jest tests for a TypeScript service or CLI. Covers colocated tests, public methods, and hand-written fakes."
languages:
    - typescript
file_types:
    - "*.test.ts"
priority: required
related_resources:
    - uri: standards://typescript/architecture
      when: "a test is pressuring production structure"
    - uri: standards://typescript/tooling
      when: "Jest config or the test script is changing"
---

# TypeScript Testing Standards
Use when writing or reviewing Jest tests for a TypeScript service or CLI. Covers colocated tests, public methods, and hand-written fakes.

Read this file. Read a related resource only when its `when` condition matches this step.

- ALWAYS follow the repository when it already uses another test runner. Repository rules win.
- ALWAYS use Jest as the runner on a new project.
- ALWAYS colocate `ClassName.test.ts` next to `ClassName.ts`.
- ALWAYS test through the public method.
- ALWAYS pass a hand-written fake into the constructor or options object.
- NEVER extract, export, or reshape production code so a test can call it.
- NEVER widen `#` field visibility so a test can reach internal state.
- Module mocks are NOT a reason to change production shape.
