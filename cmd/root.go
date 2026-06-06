package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "claudeenv",
	Short:   "Claude Code environment manager",
	Long:    "claudeenv activates domain-specific skills, agents, and rules per project.",
	Version: "dev",
}

func Execute(version string) {
	rootCmd.Version = version
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(loadCmd)
	rootCmd.AddCommand(globalCmd)
	rootCmd.AddCommand(projectCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(upgradeCmd)
}
