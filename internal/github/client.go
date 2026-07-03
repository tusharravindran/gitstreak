package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
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

func token() string {
	t := os.Getenv("GITHUB_TOKEN")
	if t == "" {
		t = os.Getenv("GH_TOKEN")
	}
	return t
}

func FetchContributions(username string) ([]ContributionDay, string, error) {
	tok := token()
	if tok == "" {
		return nil, "", fmt.Errorf("no GitHub token found — set GITHUB_TOKEN or GH_TOKEN")
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
