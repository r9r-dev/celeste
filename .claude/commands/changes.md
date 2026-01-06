---
description: Write changes from past work
allowed-tools: Bash, Read, Edit, Write
---

Prepare change logs by using `changie` cli tool. 

Change {kind} can be one of the following:
- `Added`
- `Changed`
- `Deprecated`
- `Removed`
- `Fixed`
- `Security`

For each change, run command `changie new -k {kind} -b "{body}"`

## Constraints
- Do NOT read content in .changes/ folder
- Do NOT read or modify CHANGELOG.md (changie handles this separately)
