package e2e

import (
	"path/filepath"
	"testing"
)

func TestGlobalUse_AppliesSymlinks(t *testing.T) {
	te := newTestEnv(t)
	mustMkdir(t, filepath.Join(te.global, "programming", "skills", "summarize"))
	mustMkdir(t, filepath.Join(te.global, "programming", "agents"))
	mustMkdir(t, filepath.Join(te.global, "programming", "rules"))
	mustWriteFile(t, filepath.Join(te.global, "programming", "skills", "summarize", "SKILL.md"), "# summarize")
	mustWriteFile(t, filepath.Join(te.global, "programming", "agents", "code-reviewer.md"), "# code-reviewer")
	mustWriteFile(t, filepath.Join(te.global, "programming", "rules", "engineering.md"), "# engineering")

	te.mustRun(t, "global", "use", "programming")

	assertSymlink(t, filepath.Join(te.home, ".claude", "skills", "summarize"))
	assertSymlink(t, filepath.Join(te.home, ".claude", "agents", "code-reviewer.md"))
	assertSymlink(t, filepath.Join(te.home, ".claude", "rules", "engineering.md"))
}

func TestGlobalList_ShowsAvailableProfiles(t *testing.T) {
	te := newTestEnv(t)
	for _, p := range []string{"programming", "finance"} {
		mustMkdir(t, filepath.Join(te.global, p))
	}

	out := te.mustRun(t, "global", "list")

	assertOutputContains(t, out, "programming")
	assertOutputContains(t, out, "finance")
}

func TestGlobalReset_RemovesSymlinks(t *testing.T) {
	te := newTestEnv(t)
	mustMkdir(t, filepath.Join(te.global, "programming", "rules"))
	mustWriteFile(t, filepath.Join(te.global, "programming", "rules", "engineering.md"), "# engineering")

	te.mustRun(t, "global", "use", "programming")
	assertSymlink(t, filepath.Join(te.home, ".claude", "rules", "engineering.md"))

	te.mustRun(t, "global", "reset")
	assertNoSymlink(t, filepath.Join(te.home, ".claude", "rules", "engineering.md"))
}

func TestGlobalList_EmptyOnFreshMachine(t *testing.T) {
	te := newTestEnv(t)
	// global dir exists but is empty — should return no profiles, not error
	out := te.mustRun(t, "global", "list")
	_ = out // empty output is fine
}

func TestLoad_FallsBackToStateJSON(t *testing.T) {
	te := newTestEnv(t)
	makeProfile(t, te.profiles, "golang")
	mustMkdir(t, filepath.Join(te.global, "programming", "rules"))
	mustWriteFile(t, filepath.Join(te.global, "programming", "rules", "engineering.md"), "# engineering")

	// set global via state.json (no [global] in .claudeenv)
	te.mustRun(t, "global", "use", "programming")

	// now load from a project that has no [global] in .claudeenv
	mustWriteFile(t, filepath.Join(te.project, ".claudeenv"), "[project]\ngolang\n")
	te.mustRun(t, "load", "--silent")

	// global symlinks should still be applied via state.json fallback
	assertSymlink(t, filepath.Join(te.home, ".claude", "rules", "engineering.md"))
}
