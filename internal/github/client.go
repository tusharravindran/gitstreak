package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const graphqlURL = "https://api.github.com/graphql"

type ContributionDay struct {
	Date              string `json:"date"`
	ContributionCount int    `json:"contributionCount"`
}

type ContributionWeek struct {
	ContributionDays []ContributionDay `json:"contributionDays"`
}

type ContributionCalendar struct {
	TotalContributions int                `json:"totalContributions"`
	Weeks              []ContributionWeek `json:"weeks"`
}

type ContributionsCollection struct {
	ContributionCalendar ContributionCalendar `json:"contributionCalendar"`
}

type User struct {
	Login                   string                  `json:"login"`
	Name                    string                  `json:"name"`
	ContributionsCollection ContributionsCollection `json:"contributionsCollection"`
}

type GraphQLResponse struct {
	Data struct {
		User User `json:"user"`
	} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors"`
}

// CommitDetail describes a single commit's shape for cheat-day auditing.
type CommitDetail struct {
	SHA       string
	Message   string
	Files     []string
	Additions int
	Deletions int
}

// AuditResult holds everything discovered while auditing a given day's contributions.
type AuditResult struct {
	Commits          []CommitDetail
	UnauditableRepos int
	HadContributions bool
}

type userEvent struct {
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	Repo      struct {
		Name string `json:"name"` // "owner/repo"
	} `json:"repo"`
}

type restCommit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Message string `json:"message"`
	} `json:"commit"`
	Stats *struct {
		Additions int `json:"additions"`
		Deletions int `json:"deletions"`
	} `json:"stats"`
	Files []struct {
		Filename string `json:"filename"`
	} `json:"files"`
}

func token() string {
	t := os.Getenv("GITSTREAK_GH_TOKEN")
	if t == "" {
		t = os.Getenv("GITHUB_TOKEN")
	}
	if t == "" {
		t = os.Getenv("GH_TOKEN")
	}
	return t
}

func FetchContributions(username string) ([]ContributionDay, string, error) {
	tok := token()
	if tok == "" {
		return nil, "", fmt.Errorf("no GitHub token found — set GITSTREAK_GH_TOKEN (or GITHUB_TOKEN/GH_TOKEN)")
	}

	now := time.Now()
	from := now.AddDate(-1, 0, 0).Format(time.RFC3339)
	to := now.Format(time.RFC3339)

	query := fmt.Sprintf(`{
		"query": "query { user(login: \"%s\") { login name contributionsCollection(from: \"%s\", to: \"%s\") { contributionCalendar { totalContributions weeks { contributionDays { date contributionCount } } } } } }"
	}`, username, from, to)

	req, err := http.NewRequest("POST", graphqlURL, bytes.NewBufferString(query))
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var result GraphQLResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, "", fmt.Errorf("failed to parse response: %w", err)
	}
	if len(result.Errors) > 0 {
		return nil, "", fmt.Errorf("GitHub API error: %s", result.Errors[0].Message)
	}

	var days []ContributionDay
	for _, week := range result.Data.User.ContributionsCollection.ContributionCalendar.Weeks {
		days = append(days, week.ContributionDays...)
	}

	name := result.Data.User.Name
	if name == "" {
		name = result.Data.User.Login
	}

	return days, name, nil
}

// FetchCommitDetail returns commit-level detail (files, diff size, message) for the
// given username on the given date (YYYY-MM-DD), across every repo the token can read.
// Repos the token can't read (private, no access, or anonymized) are counted as
// unauditable rather than erroring the whole call.
//
// Repo discovery goes through the events API rather than the GraphQL
// contributionsCollection aggregate: the aggregate lags well behind real time for
// same-day activity (observed empty for hours after a push), while events show up
// within seconds — necessary for auditing "today" while the streak is still being decided.
func FetchCommitDetail(username, date string) (AuditResult, error) {
	tok := token()
	if tok == "" {
		return AuditResult{}, fmt.Errorf("no GitHub token found — set GITSTREAK_GH_TOKEN (or GITHUB_TOKEN/GH_TOKEN)")
	}

	day, err := time.Parse("2006-01-02", date)
	if err != nil {
		return AuditResult{}, err
	}

	repoSet := map[string]bool{}
	for page := 1; page <= 3; page++ {
		events, more, err := fetchUserEventsPage(tok, username, page)
		if err != nil {
			return AuditResult{}, err
		}
		stop := false
		for _, e := range events {
			evDate, err := time.Parse(time.RFC3339, e.CreatedAt)
			if err != nil {
				continue
			}
			if evDate.Before(day) {
				stop = true // events are newest-first; once we're before the target day, no more matches
				continue
			}
			if e.Type == "PushEvent" && evDate.Format("2006-01-02") == date {
				repoSet[e.Repo.Name] = true
			}
		}
		if stop || !more {
			break
		}
	}

	audit := AuditResult{HadContributions: len(repoSet) > 0}
	for full := range repoSet {
		owner, repo, ok := strings.Cut(full, "/")
		if !ok {
			continue
		}
		commits, ok := fetchRepoCommitsForDay(tok, owner, repo, username, day)
		if !ok {
			audit.UnauditableRepos++
			continue
		}
		audit.Commits = append(audit.Commits, commits...)
	}

	return audit, nil
}

func fetchUserEventsPage(tok, username string, page int) ([]userEvent, bool, error) {
	url := fmt.Sprintf("https://api.github.com/users/%s/events?per_page=100&page=%d", username, page)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, false, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("GitHub events API returned %s", resp.Status)
	}

	body, _ := io.ReadAll(resp.Body)
	var events []userEvent
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, false, fmt.Errorf("failed to parse events response: %w", err)
	}
	return events, len(events) == 100, nil
}

// fetchRepoCommitsForDay lists a user's commits in a repo for a single day and fetches
// per-commit stats. ok=false means the repo was not readable with the current token
// (private/no-access) and should be treated as unauditable, not an error.
func fetchRepoCommitsForDay(tok, owner, repo, username string, day time.Time) ([]CommitDetail, bool) {
	since := day.Format(time.RFC3339)
	until := day.AddDate(0, 0, 1).Format(time.RFC3339)
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits?author=%s&since=%s&until=%s",
		owner, repo, username, since, until)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, false
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusNotFound {
		return nil, false
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false
	}

	body, _ := io.ReadAll(resp.Body)
	var list []restCommit
	if err := json.Unmarshal(body, &list); err != nil {
		return nil, false
	}

	var details []CommitDetail
	for _, c := range list {
		d, ok := fetchSingleCommit(tok, owner, repo, c.SHA)
		if !ok {
			continue
		}
		details = append(details, d)
	}
	return details, true
}

func fetchSingleCommit(tok, owner, repo, sha string) (CommitDetail, bool) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", owner, repo, sha)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return CommitDetail{}, false
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return CommitDetail{}, false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return CommitDetail{}, false
	}

	body, _ := io.ReadAll(resp.Body)
	var c restCommit
	if err := json.Unmarshal(body, &c); err != nil {
		return CommitDetail{}, false
	}

	detail := CommitDetail{SHA: c.SHA, Message: c.Commit.Message}
	if c.Stats != nil {
		detail.Additions = c.Stats.Additions
		detail.Deletions = c.Stats.Deletions
	}
	for _, f := range c.Files {
		detail.Files = append(detail.Files, f.Filename)
	}
	return detail, true
}
