package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	gh "github.com/tusharravindran/gitstreak/internal/github"
	"github.com/tusharravindran/gitstreak/internal/streak"
	"github.com/tusharravindran/gitstreak/internal/suggest"
)

var username string

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show your current GitHub streak",
	Run:   runStatus,
}

func init() {
	statusCmd.Flags().StringVarP(&username, "username", "u", "", "GitHub username (or set GITHUB_USERNAME env var)")
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) {
	if username == "" {
		username = os.Getenv("GITHUB_USERNAME")
	}
	if username == "" {
		color.Red("✗ Username required: --username <handle> or set GITHUB_USERNAME")
		os.Exit(1)
	}

	color.New(color.Faint).Print("Fetching contributions for @" + username + "... ")
	days, name, err := gh.FetchContributions(username)
	if err != nil {
		fmt.Println()
		color.Red("✗ " + err.Error())
		os.Exit(1)
	}
	fmt.Println("done")
	fmt.Println()

	result := streak.Calculate(days)

	bold := color.New(color.Bold)
	green := color.New(color.FgGreen, color.Bold)
	yellow := color.New(color.FgYellow, color.Bold)
	red := color.New(color.FgRed, color.Bold)
	faint := color.New(color.Faint)

	// Header
	bold.Printf("  %s\n", name)
	faint.Printf("  github.com/%s\n", username)
	fmt.Println()

	// Today status
	today := time.Now().Format("Mon, 02 Jan 2006")
	if result.CommittedToday {
		green.Printf("  ✅ You've committed today  ")
		faint.Printf("(%d contribution%s)\n", result.TodayCount, plural(result.TodayCount))
	} else {
		red.Printf("  ⚠️  No commits yet today  ")
		faint.Printf("(%s)\n", today)
	}
	fmt.Println()

	// Stats row
	streakColor := green
	if result.CurrentStreak == 0 {
		streakColor = red
	} else if result.CurrentStreak < 3 {
		streakColor = yellow
	}

	fmt.Printf("  %-20s %-20s %-20s\n",
		"Current Streak",
		"Longest Streak",
		"This Year",
	)
	fmt.Printf("  ")
	streakColor.Printf("%-20s", fmt.Sprintf("🔥 %d day%s", result.CurrentStreak, plural(result.CurrentStreak)))
	yellow.Printf("%-20s", fmt.Sprintf("⚡ %d day%s", result.LongestStreak, plural(result.LongestStreak)))
	color.New(color.FgCyan, color.Bold).Printf("%-20s", fmt.Sprintf("📊 %d commits", result.TotalThisYear))
	fmt.Println()
	fmt.Println()

	// Heatmap — last 4 weeks
	fmt.Print("  Last 4 weeks  ")
	recentDays := days
	if len(days) > 28 {
		recentDays = days[len(days)-28:]
	}
	for _, d := range recentDays {
		fmt.Print(heatBlock(d.ContributionCount))
	}
	fmt.Println()
	faint.Println("                " + heatBlock(0) + " 0   " + heatBlock(1) + " 1–2   " + heatBlock(3) + " 3–5   " + heatBlock(6) + " 6+")
	fmt.Println()

	// Suggestions if no commit today
	if !result.CommittedToday {
		yellow.Println("  💡 Quick tasks to stay green:")
		tasks := suggest.Pick(3)
		for _, t := range tasks {
			fmt.Printf("     %s  %s\n", t.Emoji, bold.Sprint(t.Title))
			faint.Printf("        %s  (%s)\n", t.Description, t.TimeEst)
			fmt.Println()
		}
		fmt.Println("  Run " + color.CyanString("gitstreak watch --username "+username) + " to get a 9 PM reminder.")
	}
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

func heatBlock(count int) string {
	blocks := []string{"░", "▒", "▓", "█"}
	switch {
	case count == 0:
		return color.New(color.Faint).Sprint(blocks[0])
	case count <= 2:
		return color.New(color.FgGreen).Sprint(blocks[1])
	case count <= 5:
		return color.New(color.FgGreen, color.Bold).Sprint(blocks[2])
	default:
		return color.New(color.FgHiGreen, color.Bold).Sprint(blocks[3])
	}
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}

func spaces(n int) string {
	return strings.Repeat(" ", n)
}
