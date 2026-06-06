package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// testEnv holds the directories and env vars for a single test.
type testEnv struct {
	project  string // working directory for project commands
	profiles string // CLAUDEENV_DIR_PROJECT
	global   string // CLAUDEENV_DIR_GLOBAL
	home     string // HOME (redirected so ~/.claude doesn't pollute the real one)
	env      []string
}

func newTestEnv(t *testing.T) *testEnv {
	t.Helper()
	tmp := t.TempDir()
	te := &testEnv{
		project:  filepath.Join(tmp, "project"),
		profiles: filepath.Join(tmp, "profiles"),
		global:   filepath.Join(tmp, "global"),
		home:     filepath.Join(tmp, "home"),
	}
	for _, d := range []string{
		te.project,
		te.profiles,
		te.global,
		filepath.Join(te.home, ".claude"),
	} {
		mustMkdir(t, d)
	}
	te.env = append(os.Environ(),
		"CLAUDEENV_DIR_PROJECT="+te.profiles,
		"CLAUDEENV_DIR_GLOBAL="+te.global,
		"HOME="+te.home,
	)
	return te
}

// run executes the binary in te.project with te.env.
func (te *testEnv) run(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := exec.Command(binary, args...)
	cmd.Dir = te.project
	cmd.Env = te.env
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// mustRun runs the binary and fails the test if it exits non-zero.
func (te *testEnv) mustRun(t *testing.T, args ...string) string {
	t.Helper()
	out, err := te.run(t, args...)
	if err != nil {
		t.Fatalf("command %v failed: %v\noutput: %s", args, err, out)
	}
	return out
}

// assertSymlink fails if path is not a symlink.
func assertSymlink(t *testing.T, path string) {
	t.Helper()
	info, err := os.Lstat(path)
	if err != nil {
		t.Errorf("expected symlink at %s: %v", path, err)
		return
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Errorf("expected %s to be a symlink, got mode %v", path, info.Mode())
	}
}

// assertNoSymlink fails if path is a symlink.
func assertNoSymlink(t *testing.T, path string) {
	t.Helper()
	info, err := os.Lstat(path)
	if err != nil {
		return // absent — fine
	}
	if info.Mode()&os.ModeSymlink != 0 {
		t.Errorf("expected no symlink at %s", path)
	}
}

// assertFileContains fails if path does not contain substr.
func assertFileContains(t *testing.T, path, substr string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	if !strings.Contains(string(data), substr) {
		t.Errorf("expected %q in %s\ngot:\n%s", substr, path, data)
	}
}

// assertOutputContains fails if out does not contain substr.
func assertOutputContains(t *testing.T, out, substr string) {
	t.Helper()
	if !strings.Contains(out, substr) {
		t.Errorf("expected %q in output\ngot:\n%s", substr, out)
	}
}

// makeProfile creates a full golang-style profile under base/name.
func makeProfile(t *testing.T, base, name string) {
	t.Helper()
	mustMkdir(t, filepath.Join(base, name, "skills", "go-review"))
	mustMkdir(t, filepath.Join(base, name, "skills", "go-bench"))
	mustMkdir(t, filepath.Join(base, name, "agents"))
	mustMkdir(t, filepath.Join(base, name, "rules"))
	mustWriteFile(t, filepath.Join(base, name, "skills", "go-review", "SKILL.md"), "# go-review")
	mustWriteFile(t, filepath.Join(base, name, "skills", "go-bench", "SKILL.md"), "# go-bench")
	mustWriteFile(t, filepath.Join(base, name, "agents", "go-agent.md"), "# go agent")
	mustWriteFile(t, filepath.Join(base, name, "rules", "golang.md"), "# golang rules")
}

func makeK8sProfile(t *testing.T, base string) {
	t.Helper()
	mustMkdir(t, filepath.Join(base, "kubernetes", "skills", "k8s-validate"))
	mustMkdir(t, filepath.Join(base, "kubernetes", "rules"))
	mustWriteFile(t, filepath.Join(base, "kubernetes", "skills", "k8s-validate", "SKILL.md"), "# k8s-validate")
	mustWriteFile(t, filepath.Join(base, "kubernetes", "rules", "kubernetes.md"), "# kubernetes rules")
}

func makeEtfProfile(t *testing.T, base string) {
	t.Helper()
	mustMkdir(t, filepath.Join(base, "etf", "skills", "etf-analyze"))
	mustMkdir(t, filepath.Join(base, "etf", "rules"))
	mustWriteFile(t, filepath.Join(base, "etf", "skills", "etf-analyze", "SKILL.md"), "# etf-analyze")
	mustWriteFile(t, filepath.Join(base, "etf", "rules", "etf.md"), "# etf rules")
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
