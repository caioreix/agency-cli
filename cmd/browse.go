package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/caioreix/agency-cli/internal/tui"
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "Interactively browse and install agents",
	Long:  "Launch an interactive TUI to browse, select, and install agents.",
	RunE:  runTUI,
}

func runTUI(cmd *cobra.Command, args []string) error {
	p := tea.NewProgram(tui.New(), tea.WithAltScreen())
	_, err := p.Run()
	return err
}

func init() {
	rootCmd.AddCommand(browseCmd)
}
