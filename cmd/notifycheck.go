package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/tusharravindran/gitstreak/internal/config"
	gh "github.com/tusharravindran/gitstreak/internal/github"
	"github.com/tusharravindran/gitstreak/internal/notify"
	"github.com/tusharravindran/gitstreak/internal/roast"
	"github.com/tusharravindran/gitstreak/internal/streak"
)

var notifyUsername string
var notifyForce bool
var notifyDryRun bool

var notifyCheckCmd = &cobra.Command{
	Use:    "notify-check",
	Hidden: true,
	Run:    runNotifyCheck,
}

func init() {
	notifyCheckCmd.Flags().StringVarP(&notifyUsername, "username", "u", "", "GitHub username")
	notifyCheckCmd.Flags().BoolVar(&notifyForce, "force", false, "Ignore skip days (for testing)")
	notifyCheckCmd.Flags().BoolVar(&notifyDryRun, "dry-run", false, "Send reminder notification even if already committed today")
	rootCmd.AddCommand(notifyCheckCmd)
}

func runNotifyCheck(cmd *cobra.Command, args []string) {
	cfg, _ := config.Load()

	if notifyUsername == "" {
		notifyUsername = cfg.Username
	}
	if notifyUsername == "" {
		notifyUsername = os.Getenv("GITHUB_USERNAME")
	}
	if notifyUsername == "" {
		return
	}

	// Skip if today is a skip day (unless --force is passed for testing)
	if !notifyForce && cfg.IsTodaySkipped() {
		return
	}

	days, _, err := gh.FetchContributions(notifyUsername)
	if err != nil {
		notify.Send("gitstreak ⚠️", "Could not fetch contributions: "+err.Error())
		return
	}

	result := streak.Calculate(days, cfg.SkipDays)

	if result.CommittedToday && !notifyDryRun {
		notify.SendWithVoice("gitstreak 🎉", roast.PraiseForCommit(result.CurrentStreak), "Samantha")
		return
	}

	// Work out how many days since last commit
	daysMissed := 0
	if result.LastActiveDate != "" {
		last, err := time.Parse("2006-01-02", result.LastActiveDate)
		if err == nil {
			daysMissed = int(time.Since(last).Hours() / 24)
		}
	}

	var title, message string
	if daysMissed > 3 {
		title = "gitstreak 💀"
		message = fmt.Sprintf(roast.ForBrokenStreak(daysMissed), daysMissed)
		notify.SendUrgent(title, message, "Samantha")
	} else if result.CurrentStreak > 0 {
		title = "gitstreak 🔥 streak at risk!"
		message = fmt.Sprintf("You're on a %d-day streak. Don't break it now.", result.CurrentStreak)
		notify.SendUrgent(title, message, "Samantha")
	} else {
		title = "gitstreak — commit something today"
		message = "No commits yet. Even a README update counts."
		notify.Send(title, message)
	}
}
