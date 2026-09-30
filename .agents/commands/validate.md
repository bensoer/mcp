---
description: Run fmt, typecheck, vet, test, lint, vuln, the resource harness, and pre-commit
---

Validate this repo from the root. Stop at the first non-zero exit.

```bash
make fmt
make fmt-check
make typecheck
make vet
make test
make lint
make vuln
python3 .agents/skills/test-resource-document/scripts/test_resource.py
pre-commit run --all-files
```
