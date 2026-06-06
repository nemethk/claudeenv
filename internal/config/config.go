package config

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var validProfileName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ValidateProfileName returns an error if the profile name contains unsafe characters.
func ValidateProfileName(name string) error {
	if !validProfileName.MatchString(name) {
		return fmt.Errorf("invalid profile name %q: only letters, numbers, hyphens and underscores allowed", name)
	}
	return nil
}

type Config struct {
	ProjectProfiles []string // from [project] section
	GlobalProfile   string   // from [global] section
	globalDir       string   // from [global.dir] section
}

// DefaultGlobalDir returns the default global profiles directory (~/claudeenv/claude-global).
func DefaultGlobalDir() string {
	return filepath.Join(homeDir(), "claudeenv", "claude-global")
}

// DefaultProjectDir returns the default project profiles directory (~/claudeenv/claude-project).
func DefaultProjectDir() string {
	return filepath.Join(homeDir(), "claudeenv", "claude-project")
}

// homeDir returns the current user's home directory.
// Exits immediately if it cannot be determined — the tool cannot function without it.
func homeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "claudeenv: cannot determine home directory: %v\n", err)
		os.Exit(1)
	}
	return home
}

// expandTilde expands ~ to the user's home directory.
func expandTilde(path string) string {
	if strings.HasPrefix(path, "~") {
		return filepath.Join(homeDir(), path[1:])
	}
	return path
}

// GlobalDir returns the global profiles directory, respecting [global.dir] override.
func (c *Config) GlobalDir() string {
	if c.globalDir != "" {
		return c.globalDir
	}
	if dir := os.Getenv("CLAUDEENV_DIR_GLOBAL"); dir != "" {
		return dir
	}
	return DefaultGlobalDir()
}

// Load reads .claudeenv and merges .claudeenv.local on top.
func Load() (*Config, error) {
	cfg := &Config{}

	if err := parseFile(".claudeenv", cfg, false); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	local := &Config{}
	if err := parseFile(".claudeenv.local", local, true); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	// merge local on top: global replaces, project extends, global.dir overrides
	if local.GlobalProfile != "" {
		cfg.GlobalProfile = local.GlobalProfile
	}
	if local.globalDir != "" {
		cfg.globalDir = local.globalDir
	}
	cfg.ProjectProfiles = append(cfg.ProjectProfiles, local.ProjectProfiles...)

	return cfg, nil
}

// parseFile parses a .claudeenv or .claudeenv.local file.
// When local=true, [project] lines with '+' prefix are additive; lines without '+' are ignored.
func parseFile(path string, cfg *Config, local bool) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	section := ""
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = line[1 : len(line)-1]
			continue
		}
		switch section {
		case "project":
			if local {
				if strings.HasPrefix(line, "+") {
					cfg.ProjectProfiles = append(cfg.ProjectProfiles, strings.TrimPrefix(line, "+"))
				}
				// lines without '+' are silently ignored in .claudeenv.local
			} else {
				cfg.ProjectProfiles = append(cfg.ProjectProfiles, line)
			}
		case "global":
			cfg.GlobalProfile = line
		case "global.dir":
			// only honoured in .claudeenv.local — prevents a committed repo from
			// redirecting teammates' global profiles to an arbitrary path
			if local {
				cfg.globalDir = expandTilde(line)
			}
		}
	}
	return scanner.Err()
}

// WriteProjectProfile replaces [project] in .claudeenv with a single profile,
// preserving any existing [global] (e.g. a committed team-default).
func WriteProjectProfile(p string) error {
	if err := ValidateProfileName(p); err != nil {
		return err
	}
	cfg := &Config{}
	if err := parseFile(".claudeenv", cfg, false); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return writeClaudeEnv([]string{p}, cfg.GlobalProfile)
}

// AddProjectProfile appends a profile to [project] in .claudeenv (no-op if already present).
func AddProjectProfile(p string) error {
	if err := ValidateProfileName(p); err != nil {
		return err
	}
	cfg := &Config{}
	if err := parseFile(".claudeenv", cfg, false); err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	for _, existing := range cfg.ProjectProfiles {
		if existing == p {
			return nil
		}
	}
	cfg.ProjectProfiles = append(cfg.ProjectProfiles, p)
	return writeClaudeEnv(cfg.ProjectProfiles, cfg.GlobalProfile)
}

// ResetProject removes .claudeenv from the current directory.
func ResetProject() error {
	err := os.Remove(".claudeenv")
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	return err
}

func writeClaudeEnv(projectProfiles []string, globalProfile string) error {
	var sb strings.Builder
	if len(projectProfiles) > 0 {
		sb.WriteString("[project]\n")
		for _, p := range projectProfiles {
			sb.WriteString(p)
			sb.WriteByte('\n')
		}
	}
	if globalProfile != "" {
		sb.WriteByte('\n')
		sb.WriteString("[global]\n")
		sb.WriteString(globalProfile)
		sb.WriteByte('\n')
	}
	return os.WriteFile(".claudeenv", []byte(sb.String()), 0644)
}
