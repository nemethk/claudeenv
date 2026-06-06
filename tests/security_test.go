package e2e

import (
	"path/filepath"
	"testing"
)

func TestSecurity_PathTraversalInProjectUseRejected(t *testing.T) {
	te := newTestEnv(t)

	cases := []string{"../../etc", "../passwd", "foo/bar", "foo bar", ""}
	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := te.run(t, "project", "use", name)
			if err == nil {
				t.Errorf("expected error for profile name %q, got success", name)
			}
		})
	}
}

func TestSecurity_PathTraversalInClaudeenvFileRejected(t *testing.T) {
	te := newTestEnv(t)
	mustWriteFile(t, filepath.Join(te.project, ".claudeenv"), "[project]\n../../etc\n")

	_, err := te.run(t, "load", "--silent")
	if err == nil {
		t.Error("expected error when .claudeenv contains path traversal, got success")
	}
}

func TestSecurity_GlobalDirIgnoredInCommittedFile(t *testing.T) {
	te := newTestEnv(t)
	mustWriteFile(t, filepath.Join(te.project, ".claudeenv"), "[global.dir]\n/malicious/path\n")

	// must not error — [global.dir] in committed file is silently ignored
	out, err := te.run(t, "load", "--silent")
	if err != nil {
		t.Errorf("expected no error, got: %v\noutput: %s", err, out)
	}
}

func TestSecurity_GlobalDirHonouredInLocalFile(t *testing.T) {
	te := newTestEnv(t)
	altGlobal := filepath.Join(t.TempDir(), "alt-global")
	mustMkdir(t, filepath.Join(altGlobal, "custom-profile", "rules"))
	mustWriteFile(t, filepath.Join(altGlobal, "custom-profile", "rules", "custom.md"), "# custom")

	mustWriteFile(t, filepath.Join(te.project, ".claudeenv.local"),
		"[global.dir]\n"+altGlobal+"\n\n[global]\ncustom-profile\n")

	te.mustRun(t, "load", "--silent")

	assertSymlink(t, filepath.Join(te.home, ".claude", "rules", "custom.md"))
}

func TestSecurity_GlobalUseRejectsPathTraversal(t *testing.T) {
	te := newTestEnv(t)

	cases := []string{"../../etc", "../passwd", "foo/bar", "foo bar", ""}
	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := te.run(t, "global", "use", name)
			if err == nil {
				t.Errorf("expected error for profile name %q, got success", name)
			}
		})
	}
}

func TestSecurity_NewScaffoldRejectsPathTraversal(t *testing.T) {
	te := newTestEnv(t)

	cases := []string{"../../etc", "../passwd", "foo/bar", "foo bar", ""}
	for _, name := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := te.run(t, "new", name)
			if err == nil {
				t.Errorf("expected error for profile name %q, got success", name)
			}
		})
	}
}
