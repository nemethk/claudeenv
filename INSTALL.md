# 📦 Installation Guide

Step-by-step instructions for installing and configuring `claudeenv`.

---

## ✓ Prerequisites

- **Claude Code** installed and working (`claude --version`)
- **🍎 macOS** or **🐧 Linux** (Windows not supported yet)
- One of: `curl`, `brew`, or `go` (for installation)

---

## Step 1: Install the Binary

Choose one method:

### Option A: curl (recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/nemethk/claudeenv/main/scripts/install.sh | bash
```

The script detects your OS and architecture, downloads the correct binary from the latest release, and places it in `/usr/local/bin`.

### Option B: Homebrew

```bash
brew install nemethk/tap/claudeenv
```

### Option C: go install

```bash
GOBIN=/usr/local/bin go install github.com/nemethk/claudeenv@latest
```

Requires Go 1.26 or later.

---

## Step 2: Verify the Installation

```bash
claudeenv --version
```

You should see the version number. If you get `command not found`, check that `/usr/local/bin` is in your `PATH`.

---

## Step 3: Shell Integration

Add to your shell config file:

**zsh** (`~/.zshrc`):
```bash
eval "$(claudeenv init)"
```

**bash** (`~/.bashrc`):
```bash
eval "$(claudeenv init)"
```

Reload your shell:

```bash
source ~/.zshrc   # or ~/.bashrc
```

This wraps the `claude` command so profiles are loaded automatically before each session.

### Optional: Auto-load on `cd` (zsh only)

If you want profiles to activate the moment you enter a directory (without waiting for `claude`):

```bash
eval "$(claudeenv init --hook)"
```

> **Note:** `--hook` is zsh only. bash users should use `eval "$(claudeenv init)"` without `--hook`.

---

## Step 4: Clone Your Profiles Repo

Clone the profiles repo **wherever you prefer** on your machine:

```bash
git clone git@github.com:you/my-claude-profiles.git ~/my-claude-profiles
```

No forced location — you decide where it lives.

Don't have a profiles repo yet? Scaffold your first profile into the default location and set it up later:

```bash
claudeenv new golang
# creates ~/claudeenv/claude-project/golang/skills/, agents/, rules/
```

---

## Step 5: Set Profile Locations

Tell `claudeenv` where your profiles are by adding env vars to your shell config (`~/.zshrc` or `~/.bashrc`):

```bash
export CLAUDEENV_DIR_GLOBAL=~/my-claude-profiles/claude-global
export CLAUDEENV_DIR_PROJECT=~/my-claude-profiles/claude-project
```

Global and project profiles can live in **different repos** if needed:

```bash
# personal global profiles
export CLAUDEENV_DIR_GLOBAL=~/my-personal-profiles/claude-global

# team project profiles
export CLAUDEENV_DIR_PROJECT=~/work/team-profiles/claude-project
```

Then reload:

```bash
source ~/.zshrc   # or ~/.bashrc
```

Default locations if env vars are not set:

| | Default |
|-|---------|
| Global profiles | `~/claudeenv/claude-global/` |
| Project profiles | `~/claudeenv/claude-project/` |

---

## Step 6: Set Up Your First Project

Go to any project directory and declare its profile:

```bash
cd ~/your-project
claudeenv project use golang
```

This creates a `.claudeenv` file in the current directory:

```ini
[project]
golang
```

Add it to `.gitignore` entries for generated files:

```bash
cat >> .gitignore << 'EOF'
.claudeenv.local
.claude/skills/
.claude/agents/
.claude/rules/
EOF
```

---

## Step 7: Test It

```bash
cd ~/your-project
claude
```

You should see:

```
claudeenv: loaded project profile "golang"
```

Claude Code starts with only the `golang` profile's skills, agents, and rules active.

---

## 🚀 Upgrading

### Option 1: Self-update (recommended)

```bash
sudo claudeenv upgrade
```

This downloads the latest release from GitHub and replaces your current binary. It will:
1. Fetch the latest version
2. Download the binary for your OS and architecture
3. Install it to `/usr/local/bin/claudeenv`
4. Verify the installation

### Option 2: Re-run installation script

```bash
curl -fsSL https://raw.githubusercontent.com/nemethk/claudeenv/main/scripts/install.sh | bash
```

### Option 3: Homebrew

```bash
brew upgrade claudeenv
```

### Option 4: Go

```bash
GOBIN=/usr/local/bin go install github.com/nemethk/claudeenv@latest
```

### ✓ Verify upgrade

```bash
claudeenv --version
```

You should see the latest version number.

---

## 🗑️ Uninstalling

Download and run the uninstall script:

```bash
curl -fsSL https://raw.githubusercontent.com/nemethk/claudeenv/main/scripts/uninstall.sh | bash
```

The script will:
- Remove the binary from `/usr/local/bin`
- Prompt whether to remove profiles (optional)
- Remove shell integration from `~/.zshrc` and `~/.bashrc` (with backups)

Alternatively, remove manually:

```bash
# remove the binary
sudo rm /usr/local/bin/claudeenv

# remove profiles (optional)
rm -rf ~/claudeenv

# remove shell integration
# edit ~/.zshrc and ~/.bashrc, delete the line: eval "$(claudeenv init)"
```

---

## 🔧 Troubleshooting

**`command not found: claudeenv`**
- Check that `/usr/local/bin` is in your `PATH`

**`profile not found`**
- Run `claudeenv project list` to see available profiles
- Check that `CLAUDEENV_DIR_PROJECT` points to the correct directory

**Profile not loading**
- Confirm `eval "$(claudeenv init)"` is in your shell config and the shell was reloaded
- Run `claudeenv load` manually to see any error output

**`.claude/skills/` not populated**
- Run `claudeenv load` in the project directory
- Check that the profile directory contains `skills/`, `agents/`, `rules/` subdirectories
