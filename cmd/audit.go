package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	auditpkg "github.com/tusharravindran/gitstreak/internal/audit"
	"github.com/tusharravindran/gitstreak/internal/config"
	gh "github.com/tusharravindran/gitstreak/internal/github"
)

var (
	auditUsername string
	auditDays     int
)

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Check whether your recent streak days were genuine or cheat days",
	Run:   runAudit,
}

func init() {
	auditCmd.Flags().StringVarP(&auditUsername, "username", "u", "", "GitHub username (or set GITHUB_USERNAME env var)")
	auditCmd.Flags().IntVar(&auditDays, "days", 7, "Number of past days to audit")
	rootCmd.AddCommand(auditCmd)
}

func runAudit(cmd *cobra.Command, args []string) {
	if auditUsername == "" {
		auditUsername = os.Getenv("GITHUB_USERNAME")
	}
	if auditUsername == "" {
		color.Red("✗ Username required: --username <handle> or set GITHUB_USERNAME")
		os.Exit(1)
	}

	cfg, _ := config.Load()

	bold := color.New(color.Bold)
	faint := color.New(color.Faint)
	green := color.New(color.FgGreen, color.Bold)
	magenta := color.New(color.FgMagenta, color.Bold)

	bold.Printf("  Auditing last %d days for @%s\n\n", auditDays, auditUsername)

	genuine, cheat, unauditable := 0, 0, 0
	prevFiles := []string{}

	for i := auditDays - 1; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i).Format("2006-01-02")

		result, err := gh.FetchCommitDetail(auditUsername, date)
		if err != nil {
			faint.Printf("  %s  ✗ could not audit (%s)\n", date, err.Error())
			continue
		}

		singleFile := auditpkg.SingleFile(result.Commits)

		switch {
		case len(result.Commits) == 0 && !result.HadContributions:
			faint.Printf("  %s  —  no commits\n", date)
		case len(result.Commits) == 0:
			unauditable++
			faint.Printf("  %s  —  unauditable (private/no access)\n", date)
		default:
			verdict := auditpkg.Evaluate(result.Commits, prevFiles)
			if verdict.IsCheatDay {
				cheat++
				magenta.Printf("  %s  😏 cheat day", date)
				faint.Printf("  (%s)\n", joinReasons(verdict.Reasons))
			} else {
				genuine++
				green.Printf("  %s  ✅ genuine\n", date)
			}
		}

		prevFiles = append(prevFiles, singleFile)
		if len(prevFiles) > 7 {
			prevFiles = prevFiles[len(prevFiles)-7:]
		}
		cfg.RecordAuditedFile(date, singleFile)
	}

	_ = config.Save(cfg)

	fmt.Println()
	bold.Printf("  %d genuine, %d cheat day%s, %d unauditable\n", genuine, cheat, plural(cheat), unauditable)
	if cheat > 0 {
		faint.Println("  Your real streak is a bit shorter than your GitHub graph says 👀")
	}
	fmt.Println()
}

func joinReasons(reasons []string) string {
	out := ""
	for i, r := range reasons {
		if i > 0 {
			out += ", "
		}
		out += r
	}
	return out
}
