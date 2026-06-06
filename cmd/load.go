package cmd

import (
	"github.com/nemethk/claudeenv/internal/config"
	"github.com/nemethk/claudeenv/internal/link"
	"github.com/nemethk/claudeenv/internal/profile"
	"github.com/spf13/cobra"
)

var silentFlag bool

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load profiles into .claude/ based on .claudeenv",
	RunE:  runLoad,
}

func init() {
	loadCmd.Flags().BoolVar(&silentFlag, "silent", false, "suppress output")
}

func runLoad(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// load project profiles into ./.claude/
	projectDir := profile.ProjectDir()
	if err := link.Apply(projectDir, cfg.ProjectProfiles, ".claude", silentFlag); err != nil {
		return err
	}

	// load global profile into ~/.claude/ — .claudeenv wins, falls back to state.json
	globalProfile := cfg.GlobalProfile
	if globalProfile == "" {
		globalProfile = profile.GetGlobal()
	}
	if globalProfile != "" {
		if err := link.ApplyGlobal(cfg.GlobalDir(), globalProfile, silentFlag); err != nil {
			return err
		}
	}

	return nil
}
