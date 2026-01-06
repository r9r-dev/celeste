---
description: Prepare a new release by updating version files
allowed-tools: Bash, Read, Edit, Write
---

# Objective

Prepare a new release by updating all version references in the codebase.

To do this, launch a haiku agent to follow steps precisely and respect constraints.

## Steps


1. Run `changie batch auto` to see if there are unreleased changes. Could fail with Exit code 1 but it doesn't matter (no changes).
2. Run `changie merge`
3. Get the last version number by running `changie latest`
4. Update `backend/internal/version/version.go` with the new version
5. Update `frontend/package.json` with the new version
6. Commit all changes with message "chore: bump version to {version}"
7. Push to remote
8. Create and push git tag `v{version}`
9. Return Exit Code 1 with related issues if anything failed but `changie batch auto`

## Constraints

- Do NOT read or modify CHANGELOG.md (changie handles this separately)
- The version format is semver (e.g., 0.1.0, 1.2.3)
- Strip any leading 'v' from changie output before using
