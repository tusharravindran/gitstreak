package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tusharravindran/gitstreak/internal/config"
	"github.com/tusharravindran/gitstreak/internal/roast"
)

var (
	reminderTime string
	configUser   string
	skipDays     string
	clearSkip    bool
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or update gitstreak settings",
	Run:   runConfig,
}

func init() {
	configCmd.Flags().StringVar(&reminderTime, "reminder-time", "", "Time for daily reminder, e.g. 20:00 or 8:30pm")
	configCmd.Flags().StringVar(&configUser, "username", "", "Default GitHub username")
	configCmd.Flags().StringVar(&skipDays, "skip-days", "", "Days to skip reminders, e.g. Sat,Sun or Mon,Tue,Wed")
	configCmd.Flags().BoolVar(&clearSkip, "clear-skip-days", false, "Remove all skip days")
	rootCmd.AddCommand(configCmd)
}

func runConfig(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		color.Red("✗ Could not load config: " + err.Error())
		os.Exit(1)
	}

	changed := false
	yellow := color.New(color.FgYellow, color.Bold)
	green := color.New(color.FgGreen, color.Bold)

	if configUser != "" {
		cfg.Username = configUser
		changed = true
	}

	if reminderTime != "" {
		h, m, err := parseTime(reminderTime)
		if err != nil {
			color.Red("✗ Invalid time format. Use HH:MM (e.g. 20:00) or H:MMam/pm (e.g. 8:00pm)")
			os.Exit(1)
		}
		cfg.ReminderHour = h
		cfg.ReminderMin = m
		changed = true
	}

	if clearSkip {
		cfg.SkipDays = []int{}
		changed = true
		green.Println("\n  ✓ Skip days cleared — reminders active every day")
	}

	if skipDays != "" {
		parsed, err := config.ParseDays(skipDays)
		if err != nil {
			color.Red("✗ " + err.Error())
			os.Exit(1)
		}

		consecutive := config.ConsecutiveSkipCount(parsed)
		if consecutive > 3 {
			fmt.Println()
			yellow.Printf("  %s\n", fmt.Sprintf(roast.ForSkipDays(consecutive), consecutive))
			fmt.Println()
			color.New(color.Faint).Println("  (Saved anyway — your call 🤷)")
			fmt.Println()
		}

		cfg.SkipDays = parsed
		changed = true
	}

	if changed {
		if err := config.Save(cfg); err != nil {
			color.Red("✗ Could not save config: " + err.Error())
			os.Exit(1)
		}
		if skipDays == "" && !clearSkip {
			green.Println("\n  ✓ Config saved")
		}
		fmt.Println()
	}

	printConfig(cfg)
}

func printConfig(cfg config.Config) {
	bold := color.New(color.Bold)
	faint := color.New(color.Faint)

	bold.Println("  gitstreak config")
	fmt.Println()

	uname := cfg.Username
	if uname == "" {
		uname = os.Getenv("GITHUB_USERNAME")
	}
	if uname == "" {
		uname = faint.Sprint("(not set — use --username or GITHUB_USERNAME env)")
	}

	skipLabel := faint.Sprint("none (reminders every day)")
	if len(cfg.SkipDays) > 0 {
		skipLabel = strings.Join(cfg.SkipDayNames(), ", ")
	}

	fmt.Printf("  %-22s %s\n", "username", uname)
	fmt.Printf("  %-22s %s\n", "reminder time", cfg.ReminderLabel())
	fmt.Printf("  %-22s %s\n", "skip days", skipLabel)
	fmt.Println()
	faint.Println("  Config file: ~/.config/gitstreak/config.json")
	fmt.Println()
	faint.Println("  Examples:")
	faint.Println("    gitstreak config --reminder-time 20:00")
	faint.Println("    gitstreak config --reminder-time 8:30pm")
	faint.Println("    gitstreak config --username tusharravindran")
	faint.Println("    gitstreak config --skip-days Sat,Sun")
	faint.Println("    gitstreak config --skip-days Mon,Tue,Wed,Thu,Fri   # only remind on weekends")
	faint.Println("    gitstreak config --clear-skip-days")
	fmt.Println()
	faint.Println("  After changing time or skip days, re-run: gitstreak watch")
	fmt.Println()
}

func parseTime(s string) (hour, minute int, err error) {
	s = strings.TrimSpace(strings.ToLower(s))
	isPM := strings.HasSuffix(s, "pm")
	isAM := strings.HasSuffix(s, "am")
	s = strings.TrimSuffix(s, "pm")
	s = strings.TrimSuffix(s, "am")
	s = strings.TrimSpace(s)

	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid format")
	}

	h, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return 0, 0, err
	}
	m, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return 0, 0, err
	}

	if isPM && h != 12 {
		h += 12
	}
	if isAM && h == 12 {
		h = 0
	}

	if h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, 0, fmt.Errorf("out of range")
	}

	return h, m, nil
}
