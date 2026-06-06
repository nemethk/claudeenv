package profile

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// overrideHome sets HOME to a temp dir for the duration of the test,
// so state.json operations don't touch the real home directory.
func overrideHome(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	return tmp
}

func mustMkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatal(err)
	}
}

func mustWriteFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

// --- GlobalDir ---

func TestGlobalDir_Default(t *testing.T) {
	t.Setenv("CLAUDEENV_DIR_GLOBAL", "")
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, "claudeenv", "claude-global")
	if got := GlobalDir(); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestGlobalDir_EnvVar(t *testing.T) {
	t.Setenv("CLAUDEENV_DIR_GLOBAL", "/custom/global")

	if got := GlobalDir(); got != "/custom/global" {
		t.Errorf("expected /custom/global, got %q", got)
	}
}

// --- ProjectDir ---

func TestProjectDir_Default(t *testing.T) {
	t.Setenv("CLAUDEENV_DIR_PROJECT", "")
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, "claudeenv", "claude-project")
	if got := ProjectDir(); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestProjectDir_EnvVar(t *testing.T) {
	t.Setenv("CLAUDEENV_DIR_PROJECT", "/custom/project")

	if got := ProjectDir(); got != "/custom/project" {
		t.Errorf("expected /custom/project, got %q", got)
	}
}

// --- SetGlobal / GetGlobal ---

func TestSetGlobal_GetGlobal(t *testing.T) {
	overrideHome(t)

	if err := SetGlobal("finance"); err != nil {
		t.Fatal(err)
	}
	if got := GetGlobal(); got != "finance" {
		t.Errorf("expected finance, got %q", got)
	}
}

func TestSetGlobal_Overwrites(t *testing.T) {
	overrideHome(t)

	if err := SetGlobal("finance"); err != nil {
		t.Fatal(err)
	}
	if err := SetGlobal("programming"); err != nil {
		t.Fatal(err)
	}
	if got := GetGlobal(); got != "programming" {
		t.Errorf("expected programming, got %q", got)
	}
}

func TestGetGlobal_NoFile(t *testing.T) {
	overrideHome(t)

	if got := GetGlobal(); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestGetGlobal_InvalidJSON(t *testing.T) {
	home := overrideHome(t)
	path := filepath.Join(home, "claudeenv", "state.json")
	mustWriteFile(t, path, "not valid json")

	if got := GetGlobal(); got != "" {
		t.Errorf("expected empty string on invalid JSON, got %q", got)
	}
}

// --- ResetGlobal ---

func TestResetGlobal(t *testing.T) {
	home := overrideHome(t)

	if err := SetGlobal("finance"); err != nil {
		t.Fatal(err)
	}
	if err := ResetGlobal(); err != nil {
		t.Fatal(err)
	}

	stateFile := filepath.Join(home, "claudeenv", "state.json")
	if _, err := os.Stat(stateFile); !errors.Is(err, fs.ErrNotExist) {
		t.Error("expected state.json to be removed")
	}
	if got := GetGlobal(); got != "" {
		t.Errorf("expected empty after reset, got %q", got)
	}
}

func TestResetGlobal_NoFile(t *testing.T) {
	overrideHome(t)

	if err := ResetGlobal(); err != nil {
		t.Errorf("expected no error when state.json missing, got %v", err)
	}
}

// --- ListGlobal / ListProject ---

func TestListGlobal_NonExistentDir(t *testing.T) {
	os.Setenv("CLAUDEENV_DIR_GLOBAL", "/nonexistent/path/that/does/not/exist")
	defer os.Unsetenv("CLAUDEENV_DIR_GLOBAL")

	profiles, err := ListGlobal()
	if err != nil {
		t.Errorf("expected no error for non-existent dir, got %v", err)
	}
	if len(profiles) != 0 {
		t.Errorf("expected empty list, got %v", profiles)
	}
}

func TestListProject_NonExistentDir(t *testing.T) {
	os.Setenv("CLAUDEENV_DIR_PROJECT", "/nonexistent/path/that/does/not/exist")
	defer os.Unsetenv("CLAUDEENV_DIR_PROJECT")

	profiles, err := ListProject()
	if err != nil {
		t.Errorf("expected no error for non-existent dir, got %v", err)
	}
	if len(profiles) != 0 {
		t.Errorf("expected empty list, got %v", profiles)
	}
}

func TestListGlobal(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("CLAUDEENV_DIR_GLOBAL", tmp)
	defer os.Unsetenv("CLAUDEENV_DIR_GLOBAL")

	mustMkdir(t, filepath.Join(tmp, "finance"))
	mustMkdir(t, filepath.Join(tmp, "programming"))

	profiles, err := ListGlobal()
	if err != nil {
		t.Fatal(err)
	}
	if len(profiles) != 2 {
		t.Fatalf("expected 2 profiles, got %v", profiles)
	}
}

func TestListProject(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("CLAUDEENV_DIR_PROJECT", tmp)
	defer os.Unsetenv("CLAUDEENV_DIR_PROJECT")

	mustMkdir(t, filepath.Join(tmp, "golang"))
	mustMkdir(t, filepath.Join(tmp, "kubernetes"))
	mustMkdir(t, filepath.Join(tmp, "etf"))

	profiles, err := ListProject()
	if err != nil {
		t.Fatal(err)
	}
	if len(profiles) != 3 {
		t.Fatalf("expected 3 profiles, got %v", profiles)
	}
}

func TestListGlobal_SkipsFiles(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("CLAUDEENV_DIR_GLOBAL", tmp)
	defer os.Unsetenv("CLAUDEENV_DIR_GLOBAL")

	mustMkdir(t, filepath.Join(tmp, "finance"))
	mustWriteFile(t, filepath.Join(tmp, "README.md"), "not a profile")

	profiles, err := ListGlobal()
	if err != nil {
		t.Fatal(err)
	}
	if len(profiles) != 1 || profiles[0] != "finance" {
		t.Errorf("expected only finance, got %v", profiles)
	}
}

func TestListGlobal_EmptyDir(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("CLAUDEENV_DIR_GLOBAL", tmp)
	defer os.Unsetenv("CLAUDEENV_DIR_GLOBAL")

	profiles, err := ListGlobal()
	if err != nil {
		t.Fatal(err)
	}
	if len(profiles) != 0 {
		t.Errorf("expected empty list, got %v", profiles)
	}
}

// --- Scaffold ---

func TestScaffold(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("CLAUDEENV_DIR_PROJECT", tmp)
	defer os.Unsetenv("CLAUDEENV_DIR_PROJECT")

	if err := Scaffold("golang"); err != nil {
		t.Fatal(err)
	}

	for _, sub := range []string{"skills", "agents", "rules"} {
		path := filepath.Join(tmp, "golang", sub)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("expected %s to exist: %v", sub, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("expected %s to be a directory", sub)
		}
	}
}

func TestScaffold_RejectsInvalidName(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("CLAUDEENV_DIR_PROJECT", tmp)
	defer os.Unsetenv("CLAUDEENV_DIR_PROJECT")

	for _, name := range []string{"../../etc", "foo/bar", "foo bar", ""} {
		if err := Scaffold(name); err == nil {
			t.Errorf("expected error for profile name %q, got nil", name)
		}
	}
}

func TestScaffold_AlreadyExists(t *testing.T) {
	tmp := t.TempDir()
	os.Setenv("CLAUDEENV_DIR_PROJECT", tmp)
	defer os.Unsetenv("CLAUDEENV_DIR_PROJECT")

	if err := Scaffold("golang"); err != nil {
		t.Fatal(err)
	}
	// scaffold again — should not error
	if err := Scaffold("golang"); err != nil {
		t.Errorf("expected no error on re-scaffold, got %v", err)
	}
}
