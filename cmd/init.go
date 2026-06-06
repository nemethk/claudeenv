package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var hookFlag bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Output shell integration code",
	Long:  "Prints shell code to eval in ~/.zshrc or ~/.bashrc.",
	RunE:  runInit,
}

func init() {
	initCmd.Flags().BoolVar(&hookFlag, "hook", false, "install cd hook for auto-load on directory change (zsh only)")
}

func runInit(cmd *cobra.Command, args []string) error {
	fmt.Print(`claude() {
    claudeenv load && command claude "$@"
}
`)
	if hookFlag {
		fmt.Print(`autoload -U add-zsh-hook
_claudeenv_hook() {
    [ -f .claudeenv ] && claudeenv load --silent
}
add-zsh-hook chpwd _claudeenv_hook
`)
	}
	return nil
}
