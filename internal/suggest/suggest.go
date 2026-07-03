package suggest

import (
	"math/rand"
	"time"
)

type Task struct {
	Title       string
	Description string
	TimeEst     string
	Emoji       string
}

var tasks = []Task{
	{"Write a README", "Add or improve the README for one of your repos. Counts as a commit, improves discoverability.", "10 min", "📝"},
	{"Fix a typo", "Browse your repos for typos in comments, docs, or variable names. Tiny PR, real contribution.", "5 min", "✏️"},
	{"Add a .gitignore", "Find a repo missing a proper .gitignore and add one.", "5 min", "🙈"},
	{"Write a test", "Pick an untested function and write one test. One is better than zero.", "15 min", "🧪"},
	{"Update dependencies", "Run `bundle outdated` or `go get -u` in a project. Commit the lockfile.", "10 min", "📦"},
	{"Add error handling", "Find a function that silently ignores errors. Fix it.", "15 min", "🛡️"},
	{"Write a code comment", "Find a confusing function and add a one-line comment explaining the why.", "5 min", "💬"},
	{"Create a GitHub issue", "Log a bug or feature idea in one of your repos. Planning counts.", "5 min", "🐛"},
	{"Refactor one function", "Pick a function longer than 20 lines. Split it.", "20 min", "🔧"},
	{"Add a Makefile target", "Add a useful `make dev` or `make test` target to a project.", "10 min", "⚙️"},
	{"Document an API endpoint", "Add a comment block to one undocumented endpoint.", "10 min", "📖"},
	{"Write a changelog entry", "Start or update CHANGELOG.md for your latest project.", "10 min", "📋"},
	{"Add badges to README", "Add build status, version, or license badge to a repo.", "5 min", "🏷️"},
	{"Create a GitHub Action", "Add a simple CI workflow — lint, test, or build.", "20 min", "🤖"},
	{"Tag a release", "Pick a stable project and tag a v1.0.0 release with release notes.", "10 min", "🚀"},
}

func Pick(n int) []Task {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := make([]Task, len(tasks))
	copy(shuffled, tasks)
	r.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })
	if n > len(shuffled) {
		n = len(shuffled)
	}
	return shuffled[:n]
}
