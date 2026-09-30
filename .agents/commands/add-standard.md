---
description: Add one assets/standards markdown file and prove the server serves it
argument-hint: "<category> <topic>"
---

1. Read `.agents/skills/create-resource-document` and follow it.
2. Path: `assets/standards/${1}/${2}.md`.
3. URI: `standards://${1}/${2}`. Required frontmatter: `uri` and `name`. Related links must be `related_resources`.
4. Prove it:

```bash
make build
python3 .agents/skills/test-resource-document/scripts/test_resource.py --read standards://${1}/${2}
pre-commit run --all-files
```
