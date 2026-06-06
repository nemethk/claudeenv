# Contributing

Contributions are welcome — bug fixes, new features, documentation improvements.

---

## Development Setup

**Prerequisites:**
- Go 1.26 or later
- Git

**Clone and build:**

```bash
git clone git@github.com:nemethk/claudeenv.git
cd claudeenv
go mod download
make build
```

**Verify:**

```bash
./dist/claudeenv --version
```

**Install locally:**

```bash
GOBIN=/usr/local/bin make install
claudeenv --version
```

---

## Running Tests

### Unit Tests

Go unit tests cover the three internal packages — `config`, `link`, and `profile`.

```bash
# all unit tests
make test

# with verbose output (shows each test name)
make test-verbose

# specific package
go test ./internal/config/...
go test ./internal/link/...
go test ./internal/profile/...
```

### End-to-End Tests

The `tests/` package contains Go integration tests that build and exercise the compiled binary. Each test creates an isolated environment, runs `claudeenv` commands via `os/exec`, and asserts on the filesystem and command output.

```bash
# run e2e tests (builds binary automatically)
make test-e2e

# run everything — unit + e2e
make test-all
```

| File | What it tests |
|------|---------------|
| `project_test.go` | `project use`, `add`, `list`, `reset`, profile switching |
| `global_test.go` | `global use` applies symlinks, `global list`, state.json fallback |
| `load_test.go` | `.claudeenv.local` overrides global, extends project, ignores non-`+` lines |
| `security_test.go` | path traversal rejected, `[global.dir]` ignored in committed file, honoured in local |
| `scenarios_test.go` | skill author variant switching, `new` scaffold, idempotent re-scaffold |

Each test uses `newTestEnv(t)` which creates isolated `project/`, `profiles/`, `global/`, and `home/` directories under `t.TempDir()`. `HOME` is redirected so tests never touch `~/.claude/` on the real machine.

All tests must pass before submitting a PR.

---

## Project Structure

```
claudeenv/
├── cmd/              # CLI commands (cobra)
│   ├── root.go       # root command, Execute()
│   ├── init.go       # claudeenv init
│   ├── load.go       # claudeenv load
│   ├── global.go     # claudeenv global use/list/reset
│   ├── project.go    # claudeenv project use/add/list/reset
│   ├── new.go        # claudeenv new
│   └── upgrade.go    # claudeenv upgrade
├── internal/
│   ├── config/       # parse .claudeenv and .claudeenv.local
│   ├── link/         # symlink operations
│   └── profile/      # profile directory resolution, state.json
├── tests/            # Go integration tests against the compiled binary
│   ├── main_test.go  # TestMain — builds binary before running tests
│   ├── helpers_test.go # testEnv, assert helpers, profile builders
│   ├── project_test.go
│   ├── global_test.go
│   ├── load_test.go
│   ├── security_test.go
│   └── scenarios_test.go
├── assets/           # logo and static assets
├── main.go
└── go.mod
```

---

## Making Changes

**1. Fork and create a branch**

```bash
git checkout -b feat/my-feature
```

**2. Write code and tests**

- add tests for new functionality
- run `make test-all` before committing

**3. Follow commit message conventions**

```
feat: short description of new feature
fix: short description of bug fix
docs: documentation changes
test: adding or updating tests
chore: maintenance, dependencies
```

These feed into automatic release notes — keep them clear and descriptive.

**4. Open a PR**

- describe what the change does and why
- link any related issues
- make sure CI is green

---

## Code Style

- run `go fmt ./...` before committing
- follow standard Go conventions
- keep functions small and focused
- add comments only when the **why** is non-obvious

---

## Planned Features

These are features planned for future releases:

- **Windows support** — currently Linux and macOS only
- **Profile validation** — lint profiles for common issues
- **Profile dependencies** — declare and manage profile dependencies

See [CHANGELOG.md](CHANGELOG.md) for the full roadmap.

---

## Reporting Bugs

Open an issue with:
- `claudeenv --version`
- OS and shell
- steps to reproduce
- expected vs actual behaviour
