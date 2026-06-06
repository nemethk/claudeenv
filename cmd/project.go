package cmd

import (
	"fmt"

	"github.com/nemethk/claudeenv/internal/config"
	"github.com/nemethk/claudeenv/internal/profile"
	"github.com/spf13/cobra"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage project profile (./.claude/)",
}

var projectUseCmd = &cobra.Command{
	Use:   "use <profile>",
	Short: "Set the project profile (writes .claudeenv)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return config.WriteProjectProfile(args[0])
	},
}

var projectAddCmd = &cobra.Command{
	Use:   "add <profile>",
	Short: "Add a profile to the project (appends to .claudeenv)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return config.AddProjectProfile(args[0])
	},
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available project profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		profiles, err := profile.ListProject()
		if err != nil {
			return err
		}
		for _, p := range profiles {
			fmt.Println(p)
		}
		return nil
	},
}

var projectResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Remove .claudeenv from current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		return config.ResetProject()
	},
}

func init() {
	projectCmd.AddCommand(projectUseCmd)
	projectCmd.AddCommand(projectAddCmd)
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectResetCmd)
}
