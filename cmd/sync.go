package cmd

import (
	"fmt"
	"os"

	"github.com/caioreix/agency-cli/internal/repo"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync the local agent repository cache",
	Long:  "Clone or update the local cache of the agency-agents repository. Runs git clone on first use, git pull on subsequent runs.",
	RunE: func(_ *cobra.Command, _ []string) error {
		fmt.Fprintln(os.Stdout, "⏳ Syncing agency-agents repo...")

		dir, count, err := repo.Sync()
		if err != nil {
			return fmt.Errorf("sync failed: %w", err)
		}

		if count > 0 {
			fmt.Fprintf(os.Stdout, "✓ Updated agency-agents repo (%d new commit(s))\n", count)
		} else {
			fmt.Fprintln(os.Stdout, "✓ Already up to date")
		}
		fmt.Fprintf(os.Stdout, "  Cache: %s\n", dir)

		return nil
	},
}

func init() { //nolint:gochecknoinits // required by cobra/converter
	rootCmd.AddCommand(syncCmd)
}
