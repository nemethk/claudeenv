# Release Process

Releases are fully automated via goreleaser and GitHub Actions. Cutting a release is a single command.

> Read the full release setup walkthrough: [Medium article](https://medium.com/@nemethk/TODO)

---

## Versioning

Follow [Semantic Versioning](https://semver.org):

| Change | Version bump | Example |
|--------|-------------|---------|
| Bug fix | patch | `v0.1.0` → `v0.1.1` |
| New feature (backward compatible) | minor | `v0.1.0` → `v0.2.0` |
| Breaking change | major | `v0.1.0` → `v1.0.0` |

---

## Commit Messages

Use [Conventional Commits](https://www.conventionalcommits.org) — goreleaser generates release notes from them automatically:

```bash
feat: add auto-load on cd hook
fix: remove old symlinks before reload
docs: update GUIDE.md with team scenario
chore: update dependencies
test: add path traversal security tests
```

---

## Cutting a Release

**1. Make sure all changes are committed and CI is green**

```bash
git status           # clean working tree
git log --oneline    # review commits since last release
```

**2. Tag the release**

```bash
git tag v0.1.0
git push origin v0.1.0
```

**3. GitHub Actions takes over**

goreleaser automatically:
- builds 4 binaries (linux/darwin × amd64/arm64)
- creates a GitHub Release with binaries and checksums
- updates the Homebrew tap
- generates release notes from commit messages

Monitor the release at:
```
https://github.com/nemethk/claudeenv/actions
```

**4. Verify the release**

```
https://github.com/nemethk/claudeenv/releases
```

Check that binaries are attached and release notes look correct.

---

## First Release Checklist

- [ ] Repository is public on GitHub
- [ ] `GITHUB_TOKEN` has write permissions (default for public repos)
- [ ] Homebrew tap repo exists: `github.com/nemethk/homebrew-tap`
- [ ] goreleaser config validated: `goreleaser check`
- [ ] All tests pass: `make test-all`
- [ ] Version builds cleanly: `make build`
