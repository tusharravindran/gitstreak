package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tusharravindran/gitstreak/internal/config"
	"github.com/tusharravindran/gitstreak/internal/notify"
)

var watchUsername string

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Install a 9 PM daily reminder if you haven't committed",
	Run:   runWatch,
}

var unwatchCmd = &cobra.Command{
	Use:   "unwatch",
	Short: "Remove the daily reminder",
	Run:   runUnwatch,
}

func init() {
	watchCmd.Flags().StringVarP(&watchUsername, "username", "u", "", "GitHub username")
	rootCmd.AddCommand(watchCmd)
	rootCmd.AddCommand(unwatchCmd)
}

func runWatch(cmd *cobra.Command, args []string) {
	if watchUsername == "" {
		watchUsername = os.Getenv("GITHUB_USERNAME")
	}
	if watchUsername == "" {
		color.Red("✗ Username required: --username <handle> or set GITHUB_USERNAME")
		os.Exit(1)
	}

	cfg, _ := config.Load()

	fmt.Println()
	color.New(color.Bold).Println("  Setting up daily streak reminder...")
	fmt.Println()

	if err := notify.Install(watchUsername, cfg.ReminderHour, cfg.ReminderMin); err != nil {
		color.Red("  ✗ " + err.Error())
		os.Exit(1)
	}

	color.New(color.Faint).Printf("  Every day at %s, gitstreak will check if you've committed.\n", cfg.ReminderLabel())
	color.New(color.Faint).Println("  If not, you'll get a macOS notification.")
	fmt.Println()
	color.New(color.Faint).Println("  Change time: gitstreak config --reminder-time 20:00")
	color.New(color.Faint).Println("  Remove:      gitstreak unwatch")
	fmt.Println()
}

func runUnwatch(cmd *cobra.Command, args []string) {
	fmt.Println()
	if err := notify.Uninstall(); err != nil {
		color.Red("  ✗ " + err.Error())
		os.Exit(1)
	}
	fmt.Println()
}
