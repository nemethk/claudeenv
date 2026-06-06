# 📖 Guide

Real-world walkthroughs for three `claudeenv` use cases:

1. **👥 [Team — Multi-Domain Development](#scenario-1-team--multi-domain-development)** — shared profiles repo, per-person global profiles, daily workflow
2. **🎨 [Skill Author — Iterating and Comparing Skills](#scenario-2-skill-author--iterating-and-comparing-skills)** — A/B test skill variants on the same codebase
3. **🚀 [CI/CD Pipelines — Reproducible Claude Context](#scenario-3-cicd-pipelines--reproducible-claude-context)** — consistent Claude config on every pipeline run

---

# Scenario 1: 👥 Team — Multi-Domain Development

## The Team

A small engineering team with two distinct areas of work:

| Person | Role | Daily work |
|--------|------|-----------|
| Alice | Backend engineer | Go microservices on Kubernetes |
| Bob | Backend engineer | Go microservices on Kubernetes |
| Carol | Data / Finance | ETF analysis, Python scripts |

They all use Claude Code. Without `claudeenv`, everyone sees every skill — Go, K8s, ETF, crypto — in every session, regardless of what they're working on.

---

## Step 1: Set Up the Shared Profiles Repo

The team creates a shared git repo for Claude Code profiles:

```
github.com/acme/claude-profiles
├── claude-global/
│   ├── engineering/          # global profile for engineers
│   │   ├── skills/
│   │   │   └── summarize/SKILL.md
│   │   ├── agents/
│   │   │   └── code-reviewer.md
│   │   └── rules/
│   │       └── engineering.md
│   └── finance/              # global profile for finance
│       ├── skills/
│       │   └── etf-analyze/SKILL.md
│       └── rules/
│           └── finance.md
└── claude-project/
    ├── golang/               # Go project profile
    │   ├── skills/
    │   │   ├── go-review/SKILL.md
    │   │   ├── go-bench/SKILL.md
    │   │   └── go-test/SKILL.md
    │   ├── agents/
    │   │   └── go-agent.md
    │   └── rules/
    │       └── golang.md
    ├── kubernetes/           # Kubernetes project profile
    │   ├── skills/
    │   │   ├── k8s-validate/SKILL.md
    │   │   └── k8s-review/SKILL.md
    │   └── rules/
    │       └── kubernetes.md
    └── etf/                  # ETF project profile
        ├── skills/
        │   └── etf-analyze/SKILL.md
        └── rules/
            └── etf.md
```

---

## Step 2: Each Person Clones the Profiles Repo

Everyone clones the shared repo **wherever they prefer** on their machine:

**Alice:**
```bash
git clone git@github.com:acme/claude-profiles.git ~/work/acme/claude-profiles
```

**Bob:**
```bash
git clone git@github.com:acme/claude-profiles.git ~/repos/claude-profiles
```

**Carol:**
```bash
git clone git@github.com:acme/claude-profiles.git ~/claude-profiles
```

No forced location — each person decides where it lives.

---

## Step 3: Set Env Vars (once, per machine)

Each person adds env vars to their shell config (`~/.zshrc` or `~/.bashrc`) pointing to their clone:

**Alice (`~/.zshrc` or `~/.bashrc`):**
```bash
export CLAUDEENV_DIR_GLOBAL=~/work/acme/claude-profiles/claude-global
export CLAUDEENV_DIR_PROJECT=~/work/acme/claude-profiles/claude-project
```

**Bob (`~/.zshrc` or `~/.bashrc`):**
```bash
export CLAUDEENV_DIR_GLOBAL=~/repos/claude-profiles/claude-global
export CLAUDEENV_DIR_PROJECT=~/repos/claude-profiles/claude-project
```

**Carol (`~/.zshrc` or `~/.bashrc`):**
```bash
export CLAUDEENV_DIR_GLOBAL=~/claude-profiles/claude-global
export CLAUDEENV_DIR_PROJECT=~/claude-profiles/claude-project
```

Reload the shell:
```bash
source ~/.zshrc   # or ~/.bashrc
```

Done — `claudeenv` now knows where to find profiles regardless of clone location.

### Separate Global and Project Profile Repos (optional)

Global profiles (personal identity, work style) and project profiles (domain skills) can live in **different repos**:

```bash
# personal global profiles in own private repo
git clone git@github.com:alice/my-claude-global.git ~/my-claude-global
export CLAUDEENV_DIR_GLOBAL=~/my-claude-global

# team project profiles in shared repo
git clone git@github.com:acme/claude-profiles.git ~/work/acme/claude-profiles
export CLAUDEENV_DIR_PROJECT=~/work/acme/claude-profiles/claude-project
```

Useful when global profiles contain personal preferences you don't want to share with the team.

---

## Step 4: Set the Global Profile (personal, per machine)

Each person sets their own global profile once:

**Alice and Bob (engineers):**
```bash
claudeenv global use engineering
```

**Carol (finance):**
```bash
claudeenv global use finance
```

This writes to `~/claudeenv/state.json` — personal, never committed.

---

## Step 5: Add `.claudeenv` to Each Repo

### Payment Service (Go + Kubernetes)

Alice adds `.claudeenv` to the payment service repo:

```bash
cd ~/work/payment-service
claudeenv project use golang
claudeenv project add kubernetes
```

Creates `.claudeenv`:
```ini
[project]
golang
kubernetes
```

She commits it:
```bash
git add .claudeenv
git commit -m "add claudeenv profile"
git push
```

Bob clones the repo — `.claudeenv` is already there. When he runs `claude`, the right profiles load automatically. No setup needed beyond env vars.

### ETF Dashboard (finance)

Carol adds `.claudeenv` to her ETF project:

```bash
cd ~/work/etf-dashboard
claudeenv project use etf
```

```ini
[project]
etf
```

---

## Step 6: Daily Workflow

### Alice — working on the payment service

```bash
cd ~/work/payment-service
claude
# claudeenv: loaded project profile "golang"
# claudeenv: loaded project profile "kubernetes"
# claudeenv: loaded global profile "engineering"
```

She sees only Go and Kubernetes tools. No ETF skills, no unrelated noise.

### Carol — working on ETF analysis

```bash
cd ~/work/etf-dashboard
claude
# claudeenv: loaded project profile "etf"
# claudeenv: loaded global profile "finance"
```

### Switching projects

Alice switches from the payment service to a quick infrastructure task:

```bash
cd ~/work/k8s-infra
claude
# claudeenv: loaded project profile "kubernetes"
# claudeenv: loaded global profile "engineering"
```

No manual switching — just `cd` and run `claude`.

---

## Step 7: Personal Overrides

Bob does part-time finance work. For the ETF dashboard repo, he wants his own global profile.

He creates `.claudeenv.local` in the project (never committed):
```ini
[global]
finance
```

Now when Bob runs `claude` in the ETF dashboard, he gets the `finance` global profile. In every other repo he still gets `engineering`. Alice is unaffected.

---

## Step 8: Adding a New Profile

The team starts a new crypto analytics project. Alice scaffolds a new profile inside the shared profiles repo:

```bash
cd ~/work/acme/claude-profiles
claudeenv new crypto
# scaffolded profile: crypto
```

She adds skills and rules, commits, and pushes:
```bash
git add claude-project/crypto/
git commit -m "add crypto profile"
git push
```

Everyone pulls and the new profile is available immediately:
```bash
cd ~/work/acme/claude-profiles && git pull
claudeenv project list
# golang
# kubernetes
# etf
# crypto     ← new
```

---

## What Each Person Sees

| | Alice (engineer) | Bob (engineer) | Carol (finance) |
|--|-----------------|----------------|-----------------|
| Clone location | `~/work/acme/claude-profiles` | `~/repos/claude-profiles` | `~/claude-profiles` |
| Global profile | `engineering` | `engineering` | `finance` |
| In payment-service | golang + kubernetes | golang + kubernetes | — |
| In etf-dashboard | etf | etf + finance* | etf + finance |
| Skills visible | project + engineering | project + engineering* | project + finance |

*Bob's `.claudeenv.local` in etf-dashboard overrides his global to `finance`.

## Scenario 1 Takeaways

- Clone the profiles repo anywhere — set `CLAUDEENV_DIR_*` once per machine and forget it
- `.claudeenv` committed to the repo means every teammate gets the right tools automatically on clone
- Global profile is personal and machine-specific — set it with `claudeenv global use`, never commit it
- `.claudeenv.local` is the escape hatch for personal overrides without affecting the team
- Profiles are the source of truth — update once, `git pull` everywhere

---

# Scenario 2: 🎨 Skill Author — Iterating and Comparing Skills

A different use case: you are building and maintaining Claude Code skills and rules. You want to test variants of a skill against the same codebase without manually swapping files.

## The Problem

You have written a `go-review` skill. You want to rewrite it and compare the two versions on a real Go project — same code, different skill, side by side.

## Profile Structure

Create one profile per variant in your profiles repo:

```
claude-project/
├── go-review-v1/
│   └── skills/
│       └── go-review/SKILL.md   # current stable version
├── go-review-v2/
│   └── skills/
│       └── go-review/SKILL.md   # experimental rewrite
└── golang/
    └── rules/
        └── golang.md            # shared base rules
```

## The Workflow

```bash
# test the current stable version
cd ~/work/payment-service
claudeenv project use golang
claudeenv project add go-review-v1
claude

# switch to the experimental version — same repo, same code
claudeenv project use golang
claudeenv project add go-review-v2
claude
```

Each `claudeenv project use` call removes the previous symlinks and loads the new ones. No manual file editing.

## Iterating on Rules

The same approach works for rules. If you are tuning rules for a codebase:

```
claude-project/
├── rules-strict/
│   └── rules/golang.md    # strict constraints
└── rules-permissive/
│   └── rules/golang.md    # relaxed constraints
```

```bash
claudeenv project use golang
claudeenv project add rules-strict
claude   # observe Claude's behaviour under strict rules

claudeenv project use golang
claudeenv project add rules-permissive
claude   # compare — same code, different rules
```

## Promoting to Stable

Once you are happy with the new version:

1. Copy the content into the stable profile (`golang/skills/go-review/SKILL.md`)
2. Remove the variant profiles (`go-review-v1`, `go-review-v2`)
3. Commit and push

The development loop: write → profile → test → iterate → compare → promote.

## Scenario 2 Takeaways

- One profile per variant keeps versions explicit and reversible
- Switch variants with two commands — no manual file editing
- Works for both skills and rules — compare any aspect of Claude's behaviour
- `claudeenv project use` replaces the entire profile set; `project add` layers on top

---

# Scenario 3: 🚀 CI/CD Pipelines — Reproducible Claude Context

If Claude Code is part of your automated pipeline — for code review, security scanning, documentation generation, or release preparation — `claudeenv` ensures every pipeline run gets exactly the same Claude context as every developer.

## The Problem

Without `claudeenv`, a CI job running `claude` has no way to know which skills or rules apply to that repo. It either loads everything (noisy and expensive) or nothing (no domain focus). There is no native mechanism to tie a pipeline run to a specific Claude configuration.

## How It Works

Because `.claudeenv` is committed to the repo, CI automatically inherits the project profile on clone. The pipeline just needs access to the profiles repo and a single `claudeenv load` call.

**Setup in CI (once, per environment):**

```yaml
# clone the shared profiles repo (or use a cache)
- name: Set up claudeenv profiles
  run: |
    git clone git@github.com:org/claude-profiles.git ~/profiles
    export CLAUDEENV_DIR_PROJECT=~/profiles/claude-project
    export CLAUDEENV_DIR_GLOBAL=~/profiles/claude-global
    claudeenv load
```

**Per-job profile selection via `.claudeenv`:**

```ini
# committed to the repo — CI picks this up automatically
[project]
golang
code-review
```

## Stage-Specific Profiles

Different pipeline stages can use different profiles. The cleanest approach is to keep a base `.claudeenv` committed to the repo and add stage-specific profiles via `CLAUDEENV_DIR_PROJECT`, or stack additional profiles using `claudeenv project add` as a pipeline step:

| Pipeline stage | Profile | What loads |
|----------------|---------|-----------|
| PR code review | `code-review` | review agents, style rules |
| Security scan | `security` | security-focused rules, OWASP checklist |
| Docs generation | `docs` | documentation skills, formatting rules |
| Release prep | `release` | changelog skill, release checklist |

## Why It Matters for Pipelines

- **Reproducibility** — `.claudeenv` in the repo guarantees every run uses the same Claude configuration, regardless of when or where the pipeline runs
- **Token efficiency** — automated pipelines can run Claude many times per day; loading only relevant rules reduces cost per run
- **No configuration drift** — the profile is versioned alongside the code; when the project evolves, the Claude context evolves with it

## Scenario 3 Takeaways

- No extra pipeline configuration needed — CI clones the repo and gets the profile for free
- Stage-specific behaviour is achieved by stacking profiles with `claudeenv project add` as a pipeline step
- Token cost per Claude invocation drops when only relevant rules are loaded — this compounds across many automated runs

---

## Overall Takeaways

These three scenarios share a single underlying principle: **the `.claudeenv` file is the contract between the codebase and Claude**. Commit it, version it, and everything else follows automatically.

| Scenario | The contract |
|----------|-------------|
| Team | every developer sees the same tools, from day one |
| Skill author | every variant is isolated, testable, and reversible |
| CI/CD | every pipeline run matches the developer environment exactly |

The specific mechanisms — global profiles, `.claudeenv.local`, profile stacking — are all ways to give individuals or stages a consistent starting point without breaking the shared contract.
