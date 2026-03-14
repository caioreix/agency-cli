package cmd

import (
	"fmt"

	"github.com/caioreix/agency-cli/internal/agent"
	"github.com/caioreix/agency-cli/internal/installer"
	"github.com/caioreix/agency-cli/internal/repo"
	"github.com/spf13/cobra"
)

var globalFlag bool

var getCmd = &cobra.Command{
	Use:   "get <agent-slug>",
	Short: "Download, convert and install an agent",
	Long:  "Download a specific agent by its slug, convert it to the target tool format, and install it to the correct destination.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		slug := args[0]

		if toolFlag == "" {
			return fmt.Errorf("--tool flag is required. Use 'agency-cli tools' to see available tools")
		}

		fmt.Println("⏳ Ensuring agency-agents repo is available...")
		repoDir, err := repo.EnsureRepo()
		if err != nil {
			return fmt.Errorf("failed to ensure repo: %w", err)
		}
		fmt.Println("✓ Repo ready")

		a, err := agent.FindBySlug(repoDir, slug)
		if err != nil {
			return fmt.Errorf("agent not found: %w", err)
		}

		fmt.Printf("⏳ Converting \"%s\" for %s...\n", a.Name, toolFlag)
		files, err := installer.Install(a, toolFlag, globalFlag)
		if err != nil {
			return fmt.Errorf("failed to install agent: %w", err)
		}

		fmt.Printf("✓ Converted \"%s\" for %s\n", a.Name, toolFlag)
		for _, f := range files {
			fmt.Printf("  → %s\n", f)
		}

		return nil
	},
}

func init() {
	getCmd.Flags().BoolVarP(&globalFlag, "global", "g", false, "install to global location instead of current project")
	rootCmd.AddCommand(getCmd)
}
