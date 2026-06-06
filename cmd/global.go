package cmd

import (
	"fmt"

	"github.com/nemethk/claudeenv/internal/config"
	"github.com/nemethk/claudeenv/internal/link"
	"github.com/nemethk/claudeenv/internal/profile"
	"github.com/spf13/cobra"
)

var globalCmd = &cobra.Command{
	Use:   "global",
	Short: "Manage global profile (~/.claude/)",
}

var globalUseCmd = &cobra.Command{
	Use:   "use <profile>",
	Short: "Set the active global profile and apply symlinks to ~/.claude/",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.ValidateProfileName(args[0]); err != nil {
			return err
		}
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		// apply symlinks before persisting — state.json is only written on success
		if err := link.ApplyGlobal(cfg.GlobalDir(), args[0], false); err != nil {
			return err
		}
		if err := profile.SetGlobal(args[0]); err != nil {
			return err
		}
		fmt.Printf("global profile set to: %s\n", args[0])
		return nil
	},
}

var globalListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available global profiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		profiles, err := profile.ListGlobal()
		if err != nil {
			return err
		}
		for _, p := range profiles {
			fmt.Println(p)
		}
		return nil
	},
}

var globalResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Remove active global profile and clean up ~/.claude/ symlinks",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := link.ResetGlobal(); err != nil {
			return err
		}
		return profile.ResetGlobal()
	},
}

func init() {
	globalCmd.AddCommand(globalUseCmd)
	globalCmd.AddCommand(globalListCmd)
	globalCmd.AddCommand(globalResetCmd)
}
