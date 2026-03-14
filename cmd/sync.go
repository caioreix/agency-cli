package cmd

import (
	"fmt"

	"github.com/caioreix/agency-cli/internal/repo"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync the local agent repository cache",
	Long:  "Clone or update the local cache of the agency-agents repository. Runs git clone on first use, git pull on subsequent runs.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("⏳ Syncing agency-agents repo...")

		dir, count, err := repo.Sync()
		if err != nil {
			return fmt.Errorf("sync failed: %w", err)
		}

		if count > 0 {
			fmt.Printf("✓ Updated agency-agents repo (%d new commit(s))\n", count)
		} else {
			fmt.Println("✓ Already up to date")
		}
		fmt.Printf("  Cache: %s\n", dir)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
