package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// chdir changes to a temp directory and registers cleanup to restore the original.
func chdir(t *testing.T) {
	t.Helper()
	orig, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	dir := t.TempDir()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(orig); err != nil {
			t.Errorf("restore working dir: %v", err)
		}
	})
}

func writeFile(t *testing.T, name, content string) {
	t.Helper()
	if err := os.WriteFile(name, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

// --- Load ---

func TestLoad_NoFiles(t *testing.T) {
	chdir(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(cfg.ProjectProfiles) != 0 {
		t.Errorf("expected no project profiles, got %v", cfg.ProjectProfiles)
	}
	if cfg.GlobalProfile != "" {
		t.Errorf("expected no global profile, got %q", cfg.GlobalProfile)
	}
}

func TestLoad_SingleProjectProfile(t *testing.T) {
	chdir(t)
	writeFile(t, ".claudeenv", "[project]\ngolang\n")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.ProjectProfiles) != 1 || cfg.ProjectProfiles[0] != "golang" {
		t.Errorf("expected [golang], got %v", cfg.ProjectProfiles)
	}
}

func TestLoad_MultipleProjectProfiles(t *testing.T) {
	chdir(t)
	writeFile(t, ".claudeenv", "[project]\ngolang\nkubernetes\n")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.ProjectProfiles) != 2 {
		t.Fatalf("expected 2 profiles, got %v", cfg.ProjectProfiles)
	}
	if cfg.ProjectProfiles[0] != "golang" || cfg.ProjectProfiles[1] != "kubernetes" {
		t.Errorf("unexpected profiles: %v", cfg.ProjectProfiles)
	}
}

func TestLoad_ProjectAndGlobal(t *testing.T) {
	chdir(t)
	writeFile(t, ".claudeenv", "[project]\ngolang\n\n[global]\nprogramming\n")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.ProjectProfiles) != 1 || cfg.ProjectProfiles[0] != "golang" {
		t.Errorf("unexpected project profiles: %v", cfg.ProjectProfiles)
	}
	if cfg.GlobalProfile != "programming" {
		t.Errorf("expected global=programming, got %q", cfg.GlobalProfile)
	}
}

func TestLoad_Comments(t *testing.T) {
	chdir(t)
	writeFile(t, ".claudeenv", "# this is a comment\n[project]\n# another comment\ngolang\n")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.ProjectProfiles) != 1 || cfg.ProjectProfiles[0] != "golang" {
		t.Errorf("unexpected profiles: %v", cfg.ProjectProfiles)
	}
}

// --- .claudeenv.local merging ---

func TestLoad_LocalOverridesGlobal(t *testing.T) {
	chdir(t)
	writeFile(t, ".claudeenv", "[project]\ngolang\n\n[global]\nprogramming\n")
	writeFile(t, ".claudeenv.local", "[global]\nfinance\n")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.GlobalProfile != "finance" {
		t.Errorf("expected local to override global, got %q", cfg.GlobalProfile)
	}
}

func TestLoad_LocalExtendsProject(t *testing.T) {
	chdir(t)
	writeFile(t, ".claudeenv", "[project]\ngolang\n")
	writeFile(t, ".claudeenv.local", "[project]\n+kubernetes\n")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.ProjectProfiles) != 2 {
		t.Fatalf("expected 2 profiles, got %v", cfg.ProjectProfiles)
	}
	if cfg.ProjectProfiles[1] != "kubernetes" {
		t.Errorf("expected kubernetes appended, got %v", cfg.ProjectProfiles)
	}
}

func TestLoad_LocalProjectWithoutPlusPrefixIgnored(t *testing.T) {
	chdir(t)
	writeFile(t, ".claudeenv", "[project]\ngolang\n")
	writeFile(t, ".claudeenv.local", "[project]\netf\n") // no '+' — should be ignored

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.ProjectProfiles) != 1 || cfg.ProjectProfiles[0] != "golang" {
		t.Errorf("expected only golang, got %v", cfg.ProjectProfiles)
	}
}

func TestLoad_LocalGlobalDir(t *testing.T) {
	chdir(t)
	writeFile(t, ".claudeenv", "[project]\ngolang\n")
	writeFile(t, ".claudeenv.local", "[global.dir]\n/my/personal/profiles\n\n[global]\nfinance\n")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.GlobalDir() != "/my/personal/profiles" {
		t.Errorf("expected global.dir override, got %q", cfg.GlobalDir())
	}
	if cfg.GlobalProfile != "finance" {
		t.Errorf("expected global=finance, got %q", cfg.GlobalProfile)
	}
}

// --- GlobalDir resolution ---

func TestGlobalDir_Default(t *testing.T) {
	os.Unsetenv("CLAUDEENV_DIR_GLOBAL")
	cfg := &Config{}
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, "claudeenv", "claude-global")
	if cfg.GlobalDir() != expected {
		t.Errorf("expected %q, got %q", expected, cfg.GlobalDir())
	}
}

func TestGlobalDir_EnvVar(t *testing.T) {
	os.Setenv("CLAUDEENV_DIR_GLOBAL", "/custom/global")
	defer os.Unsetenv("CLAUDEENV_DIR_GLOBAL")

	cfg := &Config{}
	if cfg.GlobalDir() != "/custom/global" {
		t.Errorf("expected env var to be used, got %q", cfg.GlobalDir())
	}
}

func TestGlobalDir_LocalOverrideWinsOverEnvVar(t *testing.T) {
	os.Setenv("CLAUDEENV_DIR_GLOBAL", "/env/global")
	defer os.Unsetenv("CLAUDEENV_DIR_GLOBAL")

	cfg := &Config{globalDir: "/local/override"}
	if cfg.GlobalDir() != "/local/override" {
		t.Errorf("expected local override to win, got %q", cfg.GlobalDir())
	}
}

// --- ValidateProfileName ---

func TestValidateProfileName_Valid(t *testing.T) {
	valid := []string{"golang", "go-review", "k8s_validate", "ETF", "my-profile-123"}
	for _, name := range valid {
		if err := ValidateProfileName(name); err != nil {
			t.Errorf("expected %q to be valid, got %v", name, err)
		}
	}
}

func TestValidateProfileName_Invalid(t *testing.T) {
	invalid := []string{"../etc", "../../passwd", "foo/bar", "foo bar", "foo;bar", ""}
	for _, name := range invalid {
		if err := ValidateProfileName(name); err == nil {
			t.Errorf("expected %q to be invalid", name)
		}
	}
}

// --- [global.dir] restricted to .claudeenv.local ---

func TestLoad_GlobalDirIgnoredInCommittedFile(t *testing.T) {
	chdir(t)
	writeFile(t, ".claudeenv", "[global.dir]\n/malicious/path\n\n[global]\nfinance\n")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.globalDir != "" {
		t.Errorf("expected [global.dir] to be ignored in .claudeenv, got %q", cfg.globalDir)
	}
}

func TestLoad_GlobalDirHonouredInLocalFile(t *testing.T) {
	chdir(t)
	writeFile(t, ".claudeenv.local", "[global.dir]\n/my/personal/profiles\n")

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.globalDir != "/my/personal/profiles" {
		t.Errorf("expected [global.dir] to be set from .claudeenv.local, got %q", cfg.globalDir)
	}
}

// --- WriteProjectProfile / AddProjectProfile / ResetProject ---

func TestWriteProjectProfile(t *testing.T) {
	chdir(t)

	if err := WriteProjectProfile("golang"); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.ProjectProfiles) != 1 || cfg.ProjectProfiles[0] != "golang" {
		t.Errorf("unexpected profiles: %v", cfg.ProjectProfiles)
	}
}

func TestAddProjectProfile(t *testing.T) {
	chdir(t)
	writeFile(t, ".claudeenv", "[project]\ngolang\n")

	if err := AddProjectProfile("kubernetes"); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.ProjectProfiles) != 2 {
		t.Fatalf("expected 2 profiles, got %v", cfg.ProjectProfiles)
	}
	if cfg.ProjectProfiles[1] != "kubernetes" {
		t.Errorf("expected kubernetes added, got %v", cfg.ProjectProfiles)
	}
}

func TestAddProjectProfile_ErrorOnUnreadableFile(t *testing.T) {
	chdir(t)
	writeFile(t, ".claudeenv", "[project]\ngolang\n")

	// make file unreadable
	if err := os.Chmod(".claudeenv", 0000); err != nil {
		t.Skip("cannot chmod (may be root)")
	}
	t.Cleanup(func() { os.Chmod(".claudeenv", 0644) })

	err := AddProjectProfile("kubernetes")
	if err == nil {
		t.Error("expected error when .claudeenv is unreadable, got nil")
	}
}

func TestAddProjectProfile_NoDuplicates(t *testing.T) {
	chdir(t)
	writeFile(t, ".claudeenv", "[project]\ngolang\n")

	if err := AddProjectProfile("golang"); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.ProjectProfiles) != 1 {
		t.Errorf("expected 1 profile after adding duplicate, got %v", cfg.ProjectProfiles)
	}
}

func TestResetProject(t *testing.T) {
	chdir(t)
	writeFile(t, ".claudeenv", "[project]\ngolang\n")

	if err := ResetProject(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(".claudeenv"); !errors.Is(err, fs.ErrNotExist) {
		t.Error("expected .claudeenv to be removed")
	}
}

func TestResetProject_NoFile(t *testing.T) {
	chdir(t)

	if err := ResetProject(); err != nil {
		t.Errorf("expected no error when file missing, got %v", err)
	}
}
