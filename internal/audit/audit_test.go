package audit

import (
	"testing"

	gh "github.com/tusharravindran/gitstreak/internal/github"
)

func TestEvaluate_NoCommits(t *testing.T) {
	v := Evaluate(nil, nil)
	if v.IsCheatDay {
		t.Fatalf("expected no verdict for no commits, got %+v", v)
	}
}

func TestEvaluate_TinyDiff(t *testing.T) {
	commits := []gh.CommitDetail{{Files: []string{"main.go"}, Additions: 1, Deletions: 0, Message: "typo"}}
	v := Evaluate(commits, nil)
	if !v.IsCheatDay {
		t.Fatalf("expected cheat day for tiny diff, got %+v", v)
	}
}

func TestEvaluate_GenuineCommit(t *testing.T) {
	commits := []gh.CommitDetail{{
		Files:     []string{"internal/foo.go", "internal/foo_test.go"},
		Additions: 80,
		Deletions: 12,
		Message:   "add retry logic to the fetch client",
	}}
	v := Evaluate(commits, nil)
	if v.IsCheatDay {
		t.Fatalf("expected genuine commit to pass, got %+v", v)
	}
}

func TestEvaluate_LazyMessage(t *testing.T) {
	commits := []gh.CommitDetail{{
		Files:     []string{"a.go", "b.go"},
		Additions: 40,
		Deletions: 10,
		Message:   "wip",
	}}
	v := Evaluate(commits, nil)
	if !v.IsCheatDay {
		t.Fatalf("expected cheat day for lazy message, got %+v", v)
	}
}

func TestEvaluate_DocsOnly(t *testing.T) {
	commits := []gh.CommitDetail{{
		Files:     []string{"README.md", "CHANGELOG.md"},
		Additions: 50,
		Deletions: 10,
		Message:   "document the new API",
	}}
	v := Evaluate(commits, nil)
	if !v.IsCheatDay {
		t.Fatalf("expected cheat day for docs-only commit, got %+v", v)
	}
}

func TestEvaluate_RepeatedSingleFile(t *testing.T) {
	commits := []gh.CommitDetail{{
		Files:     []string{"internal/foo.go"},
		Additions: 20,
		Deletions: 5,
		Message:   "expand retry handling",
	}}
	v := Evaluate(commits, []string{"internal/foo.go", "internal/foo.go"})
	if !v.IsCheatDay {
		t.Fatalf("expected cheat day for repeated single file, got %+v", v)
	}

	v2 := Evaluate(commits, []string{"internal/foo.go", "main.go"})
	if v2.IsCheatDay {
		t.Fatalf("expected no cheat day when streak is broken, got %+v", v2)
	}
}

func TestSingleFile(t *testing.T) {
	commits := []gh.CommitDetail{{Files: []string{"README.md"}}}
	if got := SingleFile(commits); got != "README.md" {
		t.Fatalf("expected README.md, got %q", got)
	}

	multi := []gh.CommitDetail{{Files: []string{"a.go", "b.go"}}}
	if got := SingleFile(multi); got != "" {
		t.Fatalf("expected empty string for multi-file commit, got %q", got)
	}
}
