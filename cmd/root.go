package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var toolFlag string

var rootCmd = &cobra.Command{
	Use:   "agency-cli",
	Short: "CLI to browse and install agents from The Agency",
	Long:  "agency-cli lets you list, download, and install AI agents from the agency-agents repository into your preferred tool (Cursor, Copilot, Windsurf, etc.).",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&toolFlag, "tool", "t", "", "target tool (claude-code, copilot, cursor, windsurf, aider, opencode, openclaw, antigravity, gemini-cli, qwen)")
}
