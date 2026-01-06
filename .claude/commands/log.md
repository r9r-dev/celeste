---
description: Write changes from past work
allowed-tools: Bash, Read, Edit, Write
---

# Objective

Prepare change logs by using `changie` cli tool. To do this, launch a haiku agent to follow steps precisely and respect constraints.

## Steps

1. For each change, run command `changie new -k {kind} -b "{body}"`

Change {kind} can be one of the following:
- `Added` for new features
- `Changed` for changes in existing functionality
- `Deprecated` for soon-to-be removed features
- `Removed` for now removed features
- `Fixed` for any bug fixes
- `Security` in case of vulnerabilities

2. When finished, commit any uncommited files.

## Constraints
- Do NOT read content in .changes/ folder
- Do NOT read or modify CHANGELOG.md (changie handles this separately)
