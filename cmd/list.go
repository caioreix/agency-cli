package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/caioreix/agency-cli/internal/agent"
	"github.com/caioreix/agency-cli/internal/repo"
	"github.com/spf13/cobra"
)

var categoryFlag string

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available agents",
	Long:  "List all available agents from the agency-agents repository. Use --category to filter by category.",
	RunE: func(cmd *cobra.Command, args []string) error {
		repoDir, err := repo.EnsureRepo()
		if err != nil {
			return fmt.Errorf("failed to ensure repo: %w", err)
		}

		agents, err := agent.ListAll(repoDir)
		if err != nil {
			return fmt.Errorf("failed to list agents: %w", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "CATEGORY\tAGENT\tDESCRIPTION")

		for _, a := range agents {
			if categoryFlag != "" && a.Category != categoryFlag {
				continue
			}

			desc := a.Description
			if len(desc) > 70 {
				desc = desc[:67] + "..."
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", a.Category, a.Name, desc)
		}

		return w.Flush()
	},
}

func init() {
	listCmd.Flags().StringVarP(&categoryFlag, "category", "c", "", "filter by category (e.g., engineering, design)")
	rootCmd.AddCommand(listCmd)
}
