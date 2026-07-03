package cmd

import (
	"os"

	"github.com/spf13/cobra"
	gh "github.com/tusharravindran/gitstreak/internal/github"
	"github.com/tusharravindran/gitstreak/internal/notify"
	"github.com/tusharravindran/gitstreak/internal/streak"
)

var notifyUsername string

// notify-check is called by launchd at 9pm — not meant for direct use
var notifyCheckCmd = &cobra.Command{
	Use:    "notify-check",
	Hidden: true,
	Run:    runNotifyCheck,
}

func init() {
	notifyCheckCmd.Flags().StringVarP(&notifyUsername, "username", "u", "", "GitHub username")
	rootCmd.AddCommand(notifyCheckCmd)
}

func runNotifyCheck(cmd *cobra.Command, args []string) {
	if notifyUsername == "" {
		notifyUsername = os.Getenv("GITHUB_USERNAME")
	}
	if notifyUsername == "" {
		return
	}

	days, _, err := gh.FetchContributions(notifyUsername)
	if err != nil {
		notify.Send("gitstreak", "⚠️ Could not fetch contributions: "+err.Error())
		return
	}

	result := streak.Calculate(days)

	if !result.CommittedToday {
		msg := "No commits yet today"
		if result.CurrentStreak > 0 {
			msg = "🔥 Your streak is at risk! " + msg
		}
		notify.Send("gitstreak — keep your streak alive", msg)
	}
}
