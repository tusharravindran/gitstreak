package audit

import (
	"fmt"
	"path/filepath"
	"strings"

	gh "github.com/tusharravindran/gitstreak/internal/github"
)

// Verdict is the outcome of running cheat-day rules against a day's commits.
type Verdict struct {
	IsCheatDay bool
	Reasons    []string
}

var lazyMessages = map[string]bool{
	"update": true, "updates": true, "wip": true, "fix": true, "fixes": true,
	"minor changes": true, "minor fix": true, "stuff": true, "changes": true,
	"asdf": true, ".": true, "commit": true, "misc": true, "small fix": true,
}

var nonCodeExt = map[string]bool{
	".md": true, ".txt": true, ".gitignore": true, ".gitattributes": true,
	".yml": true, ".yaml": true, ".json": true, ".lock": true, ".toml": true,
}

var nonCodeNames = map[string]bool{
	"license": true, "readme": true, "changelog": true, "notice": true,
}

const tinyDiffThreshold = 3

// Evaluate runs every cheat-day rule against a single day's auditable commits.
// prevFiles is the single touched file from each of the prior days (most recent last),
// used to detect "same file, day after day" farming. Pass nil/empty if unknown.
func Evaluate(commits []gh.CommitDetail, prevSingleFiles []string) Verdict {
	if len(commits) == 0 {
		return Verdict{}
	}

	var v Verdict

	totalDiff, allFiles := 0, []string{}
	for _, c := range commits {
		totalDiff += c.Additions + c.Deletions
		allFiles = append(allFiles, c.Files...)
	}

	if totalDiff > 0 && totalDiff < tinyDiffThreshold {
		v.IsCheatDay = true
		v.Reasons = append(v.Reasons, fmt.Sprintf("%d line%s changed", totalDiff, plural(totalDiff)))
	}

	if len(allFiles) == 1 {
		v.Reasons = maybeAppendRepeatedFile(v.Reasons, allFiles[0], prevSingleFiles, &v.IsCheatDay)
	}

	for _, c := range commits {
		msg := strings.ToLower(strings.TrimSpace(c.Message))
		msg = strings.SplitN(msg, "\n", 2)[0]
		if lazyMessages[msg] {
			v.IsCheatDay = true
			v.Reasons = append(v.Reasons, fmt.Sprintf("message: %q", msg))
			break
		}
	}

	if len(allFiles) > 0 && allNonCode(allFiles) {
		v.IsCheatDay = true
		v.Reasons = append(v.Reasons, "docs/config only, no source changes")
	}

	return v
}

func maybeAppendRepeatedFile(reasons []string, file string, prevSingleFiles []string, isCheatDay *bool) []string {
	streak := 1
	for i := len(prevSingleFiles) - 1; i >= 0; i-- {
		if prevSingleFiles[i] == file {
			streak++
		} else {
			break
		}
	}
	if streak >= 3 {
		*isCheatDay = true
		reasons = append(reasons, fmt.Sprintf("same file (%s) %d days in a row", file, streak))
	}
	return reasons
}

func allNonCode(files []string) bool {
	for _, f := range files {
		ext := strings.ToLower(filepath.Ext(f))
		name := strings.ToLower(strings.TrimSuffix(filepath.Base(f), ext))
		if nonCodeExt[ext] || nonCodeNames[name] {
			continue
		}
		return false
	}
	return true
}

// SingleFile returns the lone changed file for a day if exactly one file was touched
// across all commits, else "".
func SingleFile(commits []gh.CommitDetail) string {
	var files []string
	for _, c := range commits {
		files = append(files, c.Files...)
	}
	if len(files) == 1 {
		return files[0]
	}
	return ""
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
