---
uri: standards://typescript/tooling
name: TypeScript Tooling Standards
description: "Use when setting up pnpm, ESM, ESLint, Prettier, Jest, or tsconfig for a TypeScript service or CLI."
languages:
    - typescript
file_types:
    - "*.ts"
    - package.json
    - tsconfig.json
priority: required
related_resources:
    - uri: standards://typescript/architecture
      when: "the package layout is also being chosen"
    - uri: standards://typescript/testing
      when: "the Jest setup is part of the change"
    - uri: workflows://typescript-new-project
      when: "scaffolding a new TypeScript project"
---

# TypeScript Tooling Standards
Use when setting up pnpm, ESM, ESLint, Prettier, Jest, or tsconfig for a TypeScript service or CLI.

Read this file. Read a related resource only when its `when` condition matches this step.

- ALWAYS follow the repository when it already has a toolchain. Repository rules win.
- ALWAYS apply these rules when the project has no toolchain yet.
- NEVER restyle an existing repository's toolchain unless asked.

## Package manager and scripts

- ALWAYS use pnpm.
- ALWAYS use npm scripts only.
- NEVER add a Makefile.

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

## Language and format

- ALWAYS set `"type": "module"`.
- ALWAYS set TypeScript `"module": "Node16"` and `"moduleResolution": "Node16"`.
- ALWAYS use `.js` extensions on relative imports.
- ALWAYS set `strict: true` on a new `tsconfig`.
- On an existing repository, ALWAYS match that repository's `tsconfig`. NEVER tighten it unless asked.
- ALWAYS use ESLint flat config plus Prettier.
- ALWAYS typecheck with `tsc --noEmit`.
- ALWAYS use Jest, via the `ts-jest` ESM preset, as the test runner. See `standards://typescript/testing`.
