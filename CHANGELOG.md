# 📝 Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### ✨ Added
- Initial public release with core claudeenv features
- Profile management: global and project-scoped
- Symlink-based profile activation
- Shell integration (zsh and bash)
- Automated install/uninstall scripts
- Self-update command: `sudo claudeenv upgrade` — downloads latest release
- Comprehensive test suite (unit + end-to-end + upgrade tests)
- goreleaser configuration for multi-platform releases

### 🔒 Security
- Path traversal validation for profile names and paths
- [global.dir] setting honored in .claudeenv.local only, ignored in committed files
- Modern Go error handling (errors.Is instead of deprecated os.IsNotExist)

### 📚 Documentation
- README.md with use cases and command reference
- INSTALL.md with step-by-step setup and uninstall
- GUIDE.md with three real-world scenarios
- CONTRIBUTING.md for developers

### 🔮 Future (Planned)
- Windows support (currently Linux/macOS only)

---

For development history, see [git log](https://github.com/nemethk/claudeenv/commits/main).
