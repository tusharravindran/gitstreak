package cmd

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/tusharravindran/gitstreak/internal/suggest"
)

var suggestCmd = &cobra.Command{
	Use:   "suggest",
	Short: "Get quick task ideas to keep your streak green",
	Run:   runSuggest,
}

func init() {
	rootCmd.AddCommand(suggestCmd)
}

func runSuggest(cmd *cobra.Command, args []string) {
	bold := color.New(color.Bold)
	faint := color.New(color.Faint)
	yellow := color.New(color.FgYellow, color.Bold)

	fmt.Println()
	yellow.Println("  💡 Things you can commit today:")
	fmt.Println()

	tasks := suggest.Pick(5)
	for i, t := range tasks {
		fmt.Printf("  %d.  %s  %s\n", i+1, t.Emoji, bold.Sprint(t.Title))
		faint.Printf("       %s\n", t.Description)
		faint.Printf("       ⏱ %s\n", t.TimeEst)
		fmt.Println()
	}
}
