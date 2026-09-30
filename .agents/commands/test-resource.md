---
description: List registered MCP resources, or read one URI
argument-hint: "[uri]"
---

From the repo root:

```bash
make build
python3 .agents/skills/test-resource-document/scripts/test_resource.py
```

If an argument was given, also read that URI:

```bash
python3 .agents/skills/test-resource-document/scripts/test_resource.py --read ${1}
```
