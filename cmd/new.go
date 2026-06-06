package cmd

import (
	"fmt"

	"github.com/nemethk/claudeenv/internal/profile"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new <profile>",
	Short: "Scaffold a new profile directory",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := profile.Scaffold(args[0]); err != nil {
			return err
		}
		fmt.Printf("scaffolded profile: %s\n", args[0])
		return nil
	},
}
