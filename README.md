<div align="center">
  <img src="assets/claudeenv-600x400.png" width="300" alt="claudeenv" />

  <h1>claudeenv</h1>

  <p>A pyenv-style environment manager for Claude Code.<br/>
  Activate domain-specific skills, agents, and rules per project.</p>

  <p>
    <a href="https://github.com/nemethk/claudeenv/releases"><img src="https://img.shields.io/github/v/release/nemethk/claudeenv" alt="Latest Release" /></a>
    <a href="https://github.com/nemethk/claudeenv/actions/workflows/ci.yaml"><img src="https://img.shields.io/github/actions/workflow/status/nemethk/claudeenv/ci.yaml?label=ci" alt="CI" /></a>
    <a href="LICENSE"><img src="https://img.shields.io/github/license/nemethk/claudeenv" alt="License" /></a>
    <img src="https://img.shields.io/badge/go-1.26-00ADD8?logo=go" alt="Go 1.26" />
    <img src="https://img.shields.io/badge/platform-linux%20%7C%20macOS-blue" alt="Platform" />
    <a href="https://claude.ai/claude-code"><img src="https://img.shields.io/badge/Claude%20Code-enabled-8A2BE2?logo=anthropic" alt="Claude Code" /></a>
  </p>
</div>

---

## ❓ The Challenge

Claude Code loads all your skills, agents, and rules in every session. When working across multiple domains — Go, Kubernetes, ETF, Crypto — you end up with unrelated tools always visible. There is no native way to say *"in this project, only load Go tools."*

`claudeenv` solves that.

---

## ⚙️ How It Works

Each project declares its domain in a `.claudeenv` file. With the one-time shell setup below in place, `claudeenv` symlinks the matching profile's skills, agents, and rules into `.claude/` — and removes the previous ones — automatically before every `claude` session. (Skipped that step? Run `claudeenv load` directly to do the same thing on demand.)

```ini
# .claudeenv
[project]
golang

[global]
programming
```

```bash
cd ~/DevOps/Go/payment-service
claude    # loads only Go skills, agents, rules

cd ~/Finance/etf-tracker
claude    # loads only ETF skills, rules
```

---

## 💡 Why It Matters

Loading only the necessary profiles has a real impact on Claude Code's behaviour and cost:

| | Loaded at startup | Benefit |
|--|------------------|---------|
| Rules | yes — always | less context consumed, better domain focus |
| Skills | partially | cleaner list, less chance of wrong tool |
| Agents | no | UX only |

**Rules** are loaded into the context window at every session start. Loading irrelevant rules wastes tokens and dilutes Claude's attention. A Go project doesn't need finance rules — and with `claudeenv`, it won't have them.

---

## 👥 Who Uses It

| Audience | Use case |
|----------|----------|
| App developer | domain isolation — load only Go tools in Go projects, only ETF tools in finance projects |
| DevOps / infra | stage isolation — different rules for dev, test, prod environments |
| Skill author | variant comparison — A/B test skills and rules on the same codebase |
| CI/CD pipeline | reproducible context — `.claudeenv` committed to the repo guarantees identical Claude config on every run |

See [GUIDE.md](GUIDE.md) for full walkthroughs of all three scenarios.

---

## 📦 Installation

> For a complete step-by-step guide including profile setup see [INSTALL.md](INSTALL.md).

**macOS / Linux:**

```bash
curl -fsSL https://raw.githubusercontent.com/nemethk/claudeenv/main/scripts/install.sh | bash
```

**Homebrew:**

```bash
brew install nemethk/tap/claudeenv
```

**Go:**

```bash
GOBIN=/usr/local/bin go install github.com/nemethk/claudeenv@latest
```

**One-time shell setup** — add to `~/.zshrc` or `~/.bashrc`:

```bash
eval "$(claudeenv init)"
```

Reload your shell:

```bash
source ~/.zshrc   # or ~/.bashrc
```

---

## 🧪 Try It

See `claudeenv` symlink real profiles into `.claude/` — no setup beyond cloning two repos:

```bash
# profiles repo: golang + kubernetes project profiles, a "programming" global profile
git clone https://github.com/nemethk/claudeenv-demo-profiles.git ~/claudeenv-demo-profiles

# project repo — .claudeenv already declares [project] golang kubernetes / [global] programming
git clone https://github.com/nemethk/claudeenv-demo-project.git ~/claudeenv-demo-project

export CLAUDEENV_DIR_PROJECT=~/claudeenv-demo-profiles/claude-project
export CLAUDEENV_DIR_GLOBAL=~/claudeenv-demo-profiles/claude-global

cd ~/claudeenv-demo-project
claudeenv load
```

```
claudeenv: loaded project profile "golang"
claudeenv: loaded project profile "kubernetes"
claudeenv: loaded global profile "programming"
```

Inspect `.claude/skills/`, `.claude/agents/`, `.claude/rules/` — they're now symlinks into the cloned profiles.

Want this to happen automatically before every `claude` session instead of running `claudeenv load` by hand? Add to `~/.bashrc` or `~/.zshrc`:

```bash
# Claude Code - Claudeenv
eval "$(claudeenv init)"
```

---

## 🛠️ Commands

### Project — scoped to current directory

```bash
claudeenv project use golang          # set project profile, writes .claudeenv
claudeenv project add kubernetes      # add a second profile
claudeenv project list                # list available project profiles
claudeenv project reset               # remove .claudeenv
```

### Global — affects ~/.claude/ everywhere

```bash
claudeenv global use finance          # set global profile
claudeenv global list                 # list available global profiles
claudeenv global reset                # remove global profile
```

### Profile Management

```bash
claudeenv new <name>                  # scaffold a new project profile
claudeenv load                        # manually reload profiles
sudo claudeenv upgrade                # update to latest version
```

---

## `.claudeenv` File Format

**`.claudeenv`** — committed to the repo, shared with the team:

```ini
[project]
golang
kubernetes

[global]
programming
```

**`.claudeenv.local`** — gitignored, personal overrides:

```ini
# override global profile
[global]
finance

# override global profiles directory for this project
[global.dir]
~/my-personal/claudeenv/claude-global

# add to team's project profile (+ prefix required)
[project]
+etf
```

| Rule | Behaviour |
|------|-----------|
| `[global]` in local | replaces team default |
| `[global.dir]` in local | overrides `CLAUDEENV_DIR_GLOBAL` for this project |
| `[project]` in local | additive only — `+` prefix required |

> **Heads up:** `claudeenv project use <profile>` rewrites `.claudeenv` from scratch and drops any existing `[global]` line — if your committed file declares a team-default global (like the example above), use `claudeenv project add <profile>` to add further project profiles without losing it.

---

## ⚡ Configuration

Override where profiles are stored via environment variables:

```bash
# ~/.zshrc or ~/.bashrc
export CLAUDEENV_DIR_GLOBAL=~/my-profiles/claude-global
export CLAUDEENV_DIR_PROJECT=~/work/team-profiles/claude-project
```

| Variable | Default |
|----------|---------|
| `CLAUDEENV_DIR_GLOBAL` | `~/claudeenv/claude-global` |
| `CLAUDEENV_DIR_PROJECT` | `~/claudeenv/claude-project` |

---

## 🎯 Auto-load on `cd` (optional, zsh only)

Load profiles automatically when you enter a directory:

```bash
eval "$(claudeenv init --hook)"
```

> **Note:** `--hook` generates zsh-specific code (`add-zsh-hook chpwd`). bash users should use the default `eval "$(claudeenv init)"` which wraps the `claude` command instead.

---

## 🚫 `.gitignore`

Add to your projects:

```gitignore
.claudeenv.local
.claude/skills/
.claude/agents/
.claude/rules/
```

Commit `.claudeenv` — it declares the domain for the team and CI.

---

## 💫 Inspiration

`claudeenv` follows the conventions of [pyenv](https://github.com/pyenv/pyenv), [goenv](https://github.com/go-nv/goenv), and [direnv](https://github.com/direnv/direnv).

| Tool | File | Purpose |
|------|------|---------|
| pyenv | `.python-version` | Python version per project |
| goenv | `.go-version` | Go version per project |
| claudeenv | `.claudeenv` | Claude Code profile per project |

---

## 🤝 Contributing

Bug fixes, new features, and documentation improvements are welcome — see [CONTRIBUTING.md](CONTRIBUTING.md).

---

## 📜 License

[MIT](LICENSE)
