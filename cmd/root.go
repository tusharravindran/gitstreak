package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitstreak",
	Short: "Keep your GitHub streak alive",
	Long:  `gitstreak tracks your GitHub contribution streak and reminds you before you break it.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
