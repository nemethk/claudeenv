package e2e

import (
	"path/filepath"
	"testing"
)

func TestLoad_LocalOverridesGlobal(t *testing.T) {
	te := newTestEnv(t)
	makeProfile(t, te.profiles, "golang")
	mustMkdir(t, filepath.Join(te.global, "programming", "rules"))
	mustMkdir(t, filepath.Join(te.global, "finance", "rules"))
	mustWriteFile(t, filepath.Join(te.global, "programming", "rules", "engineering.md"), "# engineering")
	mustWriteFile(t, filepath.Join(te.global, "finance", "rules", "finance.md"), "# finance")

	mustWriteFile(t, filepath.Join(te.project, ".claudeenv"), "[project]\ngolang\n\n[global]\nprogramming\n")
	mustWriteFile(t, filepath.Join(te.project, ".claudeenv.local"), "[global]\nfinance\n")

	te.mustRun(t, "load", "--silent")

	assertSymlink(t, filepath.Join(te.home, ".claude", "rules", "finance.md"))
	assertNoSymlink(t, filepath.Join(te.home, ".claude", "rules", "engineering.md"))
}

func TestLoad_LocalExtendsProject(t *testing.T) {
	te := newTestEnv(t)
	makeProfile(t, te.profiles, "golang")
	makeK8sProfile(t, te.profiles)

	mustWriteFile(t, filepath.Join(te.project, ".claudeenv"), "[project]\ngolang\n")
	mustWriteFile(t, filepath.Join(te.project, ".claudeenv.local"), "[project]\n+kubernetes\n")

	te.mustRun(t, "load", "--silent")

	assertSymlink(t, filepath.Join(te.project, ".claude", "skills", "go-review"))
	assertSymlink(t, filepath.Join(te.project, ".claude", "skills", "k8s-validate"))
}

func TestLoad_LocalProjectWithoutPlusIgnored(t *testing.T) {
	te := newTestEnv(t)
	makeProfile(t, te.profiles, "golang")
	makeK8sProfile(t, te.profiles)

	mustWriteFile(t, filepath.Join(te.project, ".claudeenv"), "[project]\ngolang\n")
	// no '+' prefix — should be ignored
	mustWriteFile(t, filepath.Join(te.project, ".claudeenv.local"), "[project]\nkubernetes\n")

	te.mustRun(t, "load", "--silent")

	assertSymlink(t, filepath.Join(te.project, ".claude", "skills", "go-review"))
	assertNoSymlink(t, filepath.Join(te.project, ".claude", "skills", "k8s-validate"))
}
