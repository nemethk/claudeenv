package e2e

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSkillAuthor_SwitchVariants(t *testing.T) {
	te := newTestEnv(t)

	mustMkdir(t, filepath.Join(te.profiles, "go-review-v1", "skills", "go-review"))
	mustMkdir(t, filepath.Join(te.profiles, "go-review-v2", "skills", "go-review"))
	mustMkdir(t, filepath.Join(te.profiles, "golang", "rules"))
	mustWriteFile(t, filepath.Join(te.profiles, "go-review-v1", "skills", "go-review", "SKILL.md"), "# go-review v1 — basic")
	mustWriteFile(t, filepath.Join(te.profiles, "go-review-v2", "skills", "go-review", "SKILL.md"), "# go-review v2 — structured")
	mustWriteFile(t, filepath.Join(te.profiles, "golang", "rules", "golang.md"), "# golang rules")

	// load v1
	te.mustRun(t, "project", "use", "golang")
	te.mustRun(t, "project", "add", "go-review-v1")
	te.mustRun(t, "load", "--silent")
	assertSymlink(t, filepath.Join(te.project, ".claude", "skills", "go-review"))
	assertFileContains(t, filepath.Join(te.project, ".claude", "skills", "go-review", "SKILL.md"), "v1")

	// switch to v2 — symlink must point to v2
	te.mustRun(t, "project", "use", "golang")
	te.mustRun(t, "project", "add", "go-review-v2")
	te.mustRun(t, "load", "--silent")
	assertSymlink(t, filepath.Join(te.project, ".claude", "skills", "go-review"))
	assertFileContains(t, filepath.Join(te.project, ".claude", "skills", "go-review", "SKILL.md"), "v2")
}

func TestNewScaffold_CreatesDirectories(t *testing.T) {
	te := newTestEnv(t)

	te.mustRun(t, "new", "crypto")

	for _, sub := range []string{"skills", "agents", "rules"} {
		path := filepath.Join(te.profiles, "crypto", sub)
		info, err := os.Stat(path)
		if err != nil {
			t.Errorf("expected %s to exist: %v", path, err)
			continue
		}
		if !info.IsDir() {
			t.Errorf("expected %s to be a directory", path)
		}
	}
}

func TestNewScaffold_Idempotent(t *testing.T) {
	te := newTestEnv(t)

	te.mustRun(t, "new", "crypto")
	te.mustRun(t, "new", "crypto") // second call must not error
}
