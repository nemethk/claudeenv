package e2e

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestProjectUse_WritesClaudeenv(t *testing.T) {
	te := newTestEnv(t)
	makeProfile(t, te.profiles, "golang")

	te.mustRun(t, "project", "use", "golang")

	assertFileContains(t, filepath.Join(te.project, ".claudeenv"), "golang")
}

func TestProjectUse_LoadCreatesSymlinks(t *testing.T) {
	te := newTestEnv(t)
	makeProfile(t, te.profiles, "golang")

	te.mustRun(t, "project", "use", "golang")
	te.mustRun(t, "load", "--silent")

	assertSymlink(t, filepath.Join(te.project, ".claude", "skills", "go-review"))
	assertSymlink(t, filepath.Join(te.project, ".claude", "skills", "go-bench"))
	assertSymlink(t, filepath.Join(te.project, ".claude", "agents", "go-agent.md"))
	assertSymlink(t, filepath.Join(te.project, ".claude", "rules", "golang.md"))
}

func TestProjectAdd_StacksProfiles(t *testing.T) {
	te := newTestEnv(t)
	makeProfile(t, te.profiles, "golang")
	makeK8sProfile(t, te.profiles)

	te.mustRun(t, "project", "use", "golang")
	te.mustRun(t, "project", "add", "kubernetes")
	te.mustRun(t, "load", "--silent")

	assertSymlink(t, filepath.Join(te.project, ".claude", "skills", "go-review"))
	assertSymlink(t, filepath.Join(te.project, ".claude", "skills", "k8s-validate"))
	assertSymlink(t, filepath.Join(te.project, ".claude", "rules", "golang.md"))
	assertSymlink(t, filepath.Join(te.project, ".claude", "rules", "kubernetes.md"))
}

func TestProfileSwitch_RemovesOldSymlinks(t *testing.T) {
	te := newTestEnv(t)
	makeProfile(t, te.profiles, "golang")
	makeEtfProfile(t, te.profiles)

	te.mustRun(t, "project", "use", "golang")
	te.mustRun(t, "load", "--silent")
	assertSymlink(t, filepath.Join(te.project, ".claude", "skills", "go-review"))

	te.mustRun(t, "project", "use", "etf")
	te.mustRun(t, "load", "--silent")

	assertNoSymlink(t, filepath.Join(te.project, ".claude", "skills", "go-review"))
	assertNoSymlink(t, filepath.Join(te.project, ".claude", "skills", "go-bench"))
	assertSymlink(t, filepath.Join(te.project, ".claude", "skills", "etf-analyze"))
	assertSymlink(t, filepath.Join(te.project, ".claude", "rules", "etf.md"))
}

func TestProjectList_ShowsAvailableProfiles(t *testing.T) {
	te := newTestEnv(t)
	for _, p := range []string{"golang", "kubernetes", "etf"} {
		mustMkdir(t, filepath.Join(te.profiles, p))
	}

	out := te.mustRun(t, "project", "list")

	for _, p := range []string{"golang", "kubernetes", "etf"} {
		assertOutputContains(t, out, p)
	}
}

func TestProjectReset_RemovesClaudeenv(t *testing.T) {
	te := newTestEnv(t)
	mustWriteFile(t, filepath.Join(te.project, ".claudeenv"), "[project]\ngolang\n")

	te.mustRun(t, "project", "reset")

	if _, err := os.Stat(filepath.Join(te.project, ".claudeenv")); !errors.Is(err, fs.ErrNotExist) {
		t.Error("expected .claudeenv to be removed after reset")
	}
}
