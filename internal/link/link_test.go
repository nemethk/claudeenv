package link

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

// makeProfile creates a profile directory with skills/agents/rules and sample entries.
func makeProfile(t *testing.T, base, name string) string {
	t.Helper()
	root := filepath.Join(base, name)
	// skills — subdirectories
	mustMkdir(t, filepath.Join(root, "skills", "go-review"))
	mustMkdir(t, filepath.Join(root, "skills", "go-bench"))
	mustWriteFile(t, filepath.Join(root, "skills", "go-review", "SKILL.md"), "# go-review")
	mustWriteFile(t, filepath.Join(root, "skills", "go-bench", "SKILL.md"), "# go-bench")
	// agents — files
	mustWriteFile(t, filepath.Join(root, "agents", "go-agent.md"), "# agent")
	// rules — files
	mustWriteFile(t, filepath.Join(root, "rules", "golang.md"), "# rules")
	return root
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

func isSymlink(t *testing.T, path string) bool {
	t.Helper()
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}

// --- Apply ---

func TestApply_CreatesSymlinks(t *testing.T) {
	tmp := t.TempDir()
	profilesDir := filepath.Join(tmp, "profiles")
	targetDir := filepath.Join(tmp, "target")
	makeProfile(t, profilesDir, "golang")

	if err := Apply(profilesDir, []string{"golang"}, targetDir, true); err != nil {
		t.Fatal(err)
	}

	// skills symlinks
	if !isSymlink(t, filepath.Join(targetDir, "skills", "go-review")) {
		t.Error("expected skills/go-review to be a symlink")
	}
	if !isSymlink(t, filepath.Join(targetDir, "skills", "go-bench")) {
		t.Error("expected skills/go-bench to be a symlink")
	}
	// agents symlink
	if !isSymlink(t, filepath.Join(targetDir, "agents", "go-agent.md")) {
		t.Error("expected agents/go-agent.md to be a symlink")
	}
	// rules symlink
	if !isSymlink(t, filepath.Join(targetDir, "rules", "golang.md")) {
		t.Error("expected rules/golang.md to be a symlink")
	}
}

func TestApply_SymlinkPointsToCorrectSource(t *testing.T) {
	tmp := t.TempDir()
	profilesDir := filepath.Join(tmp, "profiles")
	targetDir := filepath.Join(tmp, "target")
	makeProfile(t, profilesDir, "golang")

	if err := Apply(profilesDir, []string{"golang"}, targetDir, true); err != nil {
		t.Fatal(err)
	}

	link := filepath.Join(targetDir, "skills", "go-review")
	dest, err := os.Readlink(link)
	if err != nil {
		t.Fatal(err)
	}
	expected, err := filepath.Abs(filepath.Join(profilesDir, "golang", "skills", "go-review"))
	if err != nil {
		t.Fatal(err)
	}
	if dest != expected {
		t.Errorf("symlink points to %q, want %q", dest, expected)
	}
}

func TestApply_MultipleProfiles(t *testing.T) {
	tmp := t.TempDir()
	profilesDir := filepath.Join(tmp, "profiles")
	targetDir := filepath.Join(tmp, "target")
	makeProfile(t, profilesDir, "golang")

	// second profile with different skills
	mustMkdir(t, filepath.Join(profilesDir, "kubernetes", "skills", "k8s-validate"))
	mustWriteFile(t, filepath.Join(profilesDir, "kubernetes", "skills", "k8s-validate", "SKILL.md"), "# k8s")

	if err := Apply(profilesDir, []string{"golang", "kubernetes"}, targetDir, true); err != nil {
		t.Fatal(err)
	}

	if !isSymlink(t, filepath.Join(targetDir, "skills", "go-review")) {
		t.Error("expected go-review symlink")
	}
	if !isSymlink(t, filepath.Join(targetDir, "skills", "k8s-validate")) {
		t.Error("expected k8s-validate symlink")
	}
}

func TestApply_ProfileNotFound(t *testing.T) {
	tmp := t.TempDir()
	targetDir := filepath.Join(tmp, "target")

	err := Apply(tmp, []string{"nonexistent"}, targetDir, true)
	if err == nil {
		t.Error("expected error for missing profile")
	}
}

func TestApply_EmptyProfileSubdirs(t *testing.T) {
	tmp := t.TempDir()
	profilesDir := filepath.Join(tmp, "profiles")
	targetDir := filepath.Join(tmp, "target")

	// profile exists but has no skills/agents/rules entries
	mustMkdir(t, filepath.Join(profilesDir, "empty"))

	if err := Apply(profilesDir, []string{"empty"}, targetDir, true); err != nil {
		t.Errorf("expected no error for empty profile, got %v", err)
	}
}

func TestApply_NoProfiles(t *testing.T) {
	tmp := t.TempDir()
	targetDir := filepath.Join(tmp, "target")

	if err := Apply(tmp, []string{}, targetDir, true); err != nil {
		t.Errorf("expected no error for empty profile list, got %v", err)
	}
}

// --- security: path traversal ---

func TestApply_RejectsPathTraversalInProfileName(t *testing.T) {
	tmp := t.TempDir()
	targetDir := filepath.Join(tmp, "target")

	err := Apply(tmp, []string{"../../etc"}, targetDir, true)
	if err == nil {
		t.Error("expected error for path traversal profile name")
	}
}

func TestApply_RejectsSlashInProfileName(t *testing.T) {
	tmp := t.TempDir()
	targetDir := filepath.Join(tmp, "target")

	err := Apply(tmp, []string{"foo/bar"}, targetDir, true)
	if err == nil {
		t.Error("expected error for slash in profile name")
	}
}

// --- removeSymlinks ---

func TestApply_RemovesOldSymlinksOnReload(t *testing.T) {
	tmp := t.TempDir()
	profilesDir := filepath.Join(tmp, "profiles")
	targetDir := filepath.Join(tmp, "target")
	makeProfile(t, profilesDir, "golang")

	// first load — golang
	if err := Apply(profilesDir, []string{"golang"}, targetDir, true); err != nil {
		t.Fatal(err)
	}
	if !isSymlink(t, filepath.Join(targetDir, "skills", "go-review")) {
		t.Fatal("expected go-review after first load")
	}

	// second profile without golang skills
	mustMkdir(t, filepath.Join(profilesDir, "etf", "skills", "etf-analyze"))
	mustWriteFile(t, filepath.Join(profilesDir, "etf", "skills", "etf-analyze", "SKILL.md"), "# etf")

	// reload with etf only
	if err := Apply(profilesDir, []string{"etf"}, targetDir, true); err != nil {
		t.Fatal(err)
	}

	// go-review symlink should be gone
	if isSymlink(t, filepath.Join(targetDir, "skills", "go-review")) {
		t.Error("expected go-review symlink to be removed after profile switch")
	}
	// etf-analyze should be present
	if !isSymlink(t, filepath.Join(targetDir, "skills", "etf-analyze")) {
		t.Error("expected etf-analyze symlink after reload")
	}
}

func TestRemoveSymlinks_LeavesRealFiles(t *testing.T) {
	tmp := t.TempDir()

	// create a real file and a symlink in the same dir
	realFile := filepath.Join(tmp, "real.md")
	mustWriteFile(t, realFile, "real")

	target := filepath.Join(tmp, "link.md")
	src := filepath.Join(tmp, "src.md")
	mustWriteFile(t, src, "src")
	if err := os.Symlink(src, target); err != nil {
		t.Fatal(err)
	}

	if err := removeSymlinks(tmp); err != nil {
		t.Fatal(err)
	}

	// symlink should be gone
	if isSymlink(t, target) {
		t.Error("expected symlink to be removed")
	}
	// real file should remain
	if _, err := os.Stat(realFile); errors.Is(err, fs.ErrNotExist) {
		t.Error("expected real file to remain")
	}
}

func TestRemoveSymlinks_MissingDir(t *testing.T) {
	err := removeSymlinks("/nonexistent/path")
	if err != nil {
		t.Errorf("expected no error for missing dir, got %v", err)
	}
}

func TestRemoveSymlinks_RemovesSymlinks(t *testing.T) {
	tmp := t.TempDir()

	src := filepath.Join(tmp, "src.md")
	link := filepath.Join(tmp, "link.md")
	mustWriteFile(t, src, "content")
	if err := os.Symlink(src, link); err != nil {
		t.Fatal(err)
	}

	if err := removeSymlinks(tmp); err != nil {
		t.Errorf("expected no error removing symlink, got %v", err)
	}
	if _, err := os.Lstat(link); !errors.Is(err, fs.ErrNotExist) {
		t.Error("expected symlink to be removed")
	}
}

func TestApply_PreservesExistingOnValidationFailure(t *testing.T) {
	tmp := t.TempDir()
	profilesDir := filepath.Join(tmp, "profiles")
	targetDir := filepath.Join(tmp, "target")
	makeProfile(t, profilesDir, "golang")

	// load golang successfully
	if err := Apply(profilesDir, []string{"golang"}, targetDir, true); err != nil {
		t.Fatal(err)
	}
	if !isSymlink(t, filepath.Join(targetDir, "skills", "go-review")) {
		t.Fatal("expected go-review symlink after first load")
	}

	// attempt reload with an invalid second profile — should fail without wiping .claude/
	err := Apply(profilesDir, []string{"golang", "../../etc"}, targetDir, true)
	if err == nil {
		t.Fatal("expected error for invalid profile name")
	}

	// existing symlinks must still be intact
	if !isSymlink(t, filepath.Join(targetDir, "skills", "go-review")) {
		t.Error("existing symlinks must be preserved when validation fails")
	}
}

func TestApply_ErrorOnRealFileAtDestination(t *testing.T) {
	tmp := t.TempDir()
	profilesDir := filepath.Join(tmp, "profiles")
	targetDir := filepath.Join(tmp, "target")
	makeProfile(t, profilesDir, "golang")

	// create a real file where a symlink would be placed
	mustMkdir(t, filepath.Join(targetDir, "skills"))
	mustWriteFile(t, filepath.Join(targetDir, "skills", "go-review"), "real file, not a symlink")

	err := Apply(profilesDir, []string{"golang"}, targetDir, true)
	if err == nil {
		t.Error("expected error when destination has a real file, got nil")
	}
}
