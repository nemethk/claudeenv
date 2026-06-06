package link

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/nemethk/claudeenv/internal/config"
)

// subdirs are the three directories claudeenv manages inside .claude/ and ~/.claude/.
var subdirs = []string{"skills", "agents", "rules"}

// Apply symlinks skills/agents/rules from profilesDir into targetDir.
// All profiles are validated before any filesystem changes are made — an
// invalid profile name never leaves .claude/ in a partially-updated state.
func Apply(profilesDir string, profiles []string, targetDir string, silent bool) error {
	// pass 1: validate every profile and resolve absolute paths — no side effects
	absProfiles, err := filepath.Abs(profilesDir)
	if err != nil {
		return fmt.Errorf("resolve profiles dir: %w", err)
	}

	type entry struct {
		name string
		path string
	}
	entries := make([]entry, 0, len(profiles))
	for _, profileName := range profiles {
		if err := config.ValidateProfileName(profileName); err != nil {
			return err
		}
		profilePath := filepath.Join(profilesDir, profileName)
		absProfile, err := filepath.Abs(profilePath)
		if err != nil {
			return fmt.Errorf("resolve profile path: %w", err)
		}
		if !strings.HasPrefix(absProfile, absProfiles+string(os.PathSeparator)) {
			return fmt.Errorf("profile path escapes profiles directory: %s", profileName)
		}
		if _, err := os.Stat(profilePath); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return fmt.Errorf("profile not found: %s", profilePath)
			}
			return fmt.Errorf("cannot access profile %s: %w", profilePath, err)
		}
		entries = append(entries, entry{profileName, profilePath})
	}

	// pass 2: all profiles valid — safe to remove old symlinks
	for _, sub := range subdirs {
		if err := removeSymlinks(filepath.Join(targetDir, sub)); err != nil {
			return err
		}
	}

	// pass 3: create new symlinks
	for _, e := range entries {
		for _, sub := range subdirs {
			src := filepath.Join(e.path, sub)
			dst := filepath.Join(targetDir, sub)
			if err := os.MkdirAll(dst, 0755); err != nil {
				return err
			}
			if err := symlinkDir(src, dst); err != nil {
				return err
			}
		}
		if !silent {
			fmt.Printf("claudeenv: loaded project profile %q\n", e.name)
		}
	}
	return nil
}

// ApplyGlobal symlinks a global profile into ~/.claude/.
func ApplyGlobal(globalDir, profile string, silent bool) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("home directory: %w", err)
	}
	targetDir := filepath.Join(home, ".claude")
	if err := Apply(globalDir, []string{profile}, targetDir, true); err != nil {
		return err
	}
	if !silent {
		fmt.Printf("claudeenv: loaded global profile %q\n", profile)
	}
	return nil
}

// ResetGlobal removes all claudeenv-managed symlinks from ~/.claude/.
func ResetGlobal() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("home directory: %w", err)
	}
	targetDir := filepath.Join(home, ".claude")
	for _, sub := range subdirs {
		if err := removeSymlinks(filepath.Join(targetDir, sub)); err != nil {
			return err
		}
	}
	return nil
}

// symlinkDir creates symlinks for all entries in src into dst.
// For skills: symlinks directories. For agents/rules: symlinks files.
func symlinkDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}

	for _, e := range entries {
		srcRel := filepath.Join(src, e.Name())
		srcPath, err := filepath.Abs(srcRel)
		if err != nil {
			return fmt.Errorf("resolve path %s: %w", srcRel, err)
		}
		dstPath := filepath.Join(dst, e.Name())

		if info, err := os.Lstat(dstPath); err == nil {
			if info.Mode()&os.ModeSymlink == 0 {
				return fmt.Errorf("cannot create symlink at %s: path exists and is not a symlink", dstPath)
			}
			if err := os.Remove(dstPath); err != nil {
				return fmt.Errorf("remove %s: %w", dstPath, err)
			}
		}

		if err := os.Symlink(srcPath, dstPath); err != nil {
			return err
		}
	}
	return nil
}

// removeSymlinks removes all symlinks in dir (leaves real files/dirs untouched).
func removeSymlinks(dir string) error {
	entries, err := os.ReadDir(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}

	for _, e := range entries {
		path := filepath.Join(dir, e.Name())
		info, err := os.Lstat(path)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return fmt.Errorf("stat %s: %w", path, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("remove symlink %s: %w", path, err)
			}
		}
	}
	return nil
}
