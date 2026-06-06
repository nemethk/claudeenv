package profile

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/nemethk/claudeenv/internal/config"
)

const stateFile = "state.json"

type state struct {
	Global string `json:"global"`
}

// GlobalDir returns the global profiles directory from env or default.
func GlobalDir() string {
	if dir := os.Getenv("CLAUDEENV_DIR_GLOBAL"); dir != "" {
		return dir
	}
	return config.DefaultGlobalDir()
}

// ProjectDir returns the project profiles directory from env or default.
func ProjectDir() string {
	if dir := os.Getenv("CLAUDEENV_DIR_PROJECT"); dir != "" {
		return dir
	}
	return config.DefaultProjectDir()
}

// ListGlobal returns available global profile names.
func ListGlobal() ([]string, error) {
	return listDir(GlobalDir())
}

// ListProject returns available project profile names.
func ListProject() ([]string, error) {
	return listDir(ProjectDir())
}

// SetGlobal persists the active global profile to state.json.
func SetGlobal(p string) error {
	path, err := stateFilePath()
	if err != nil {
		return err
	}
	s := state{Global: p}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// GetGlobal returns the persisted global profile from state.json.
func GetGlobal() string {
	path, err := stateFilePath()
	if err != nil {
		return ""
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var s state
	if err := json.Unmarshal(data, &s); err != nil {
		return ""
	}
	return s.Global
}

// ResetGlobal clears the persisted global profile.
func ResetGlobal() error {
	path, err := stateFilePath()
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	return err
}

// Scaffold creates the directory structure for a new project profile.
func Scaffold(name string) error {
	if err := config.ValidateProfileName(name); err != nil {
		return err
	}
	base := filepath.Join(ProjectDir(), name)
	for _, sub := range []string{"skills", "agents", "rules"} {
		if err := os.MkdirAll(filepath.Join(base, sub), 0755); err != nil {
			return err
		}
	}
	return nil
}

func stateFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home directory: %w", err)
	}
	return filepath.Join(home, "claudeenv", stateFile), nil
}

func listDir(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return []string{}, nil
	}
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}
