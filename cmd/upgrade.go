package cmd

import (
	"github.com/nemethk/claudeenv/internal/upgrade"
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade claudeenv to the latest version",
	Long: `Download and install the latest version of claudeenv.

Requires sudo to replace the binary at /usr/local/bin/claudeenv.

Usage:
  sudo claudeenv upgrade

This command:
  1. Fetches the latest release from GitHub
  2. Downloads the binary for your OS and architecture
  3. Installs it to /usr/local/bin/claudeenv
  4. Verifies the installation`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return upgrade.Run()
	},
}
