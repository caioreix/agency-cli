package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/caioreix/agency-cli/internal/converter"
	"github.com/spf13/cobra"
)

const tabwriterPadding = 2

var toolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "List supported tools",
	Long:  "List all supported target tools and their installation destinations.",
	Run: func(_ *cobra.Command, _ []string) {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, tabwriterPadding, ' ', 0)
		fmt.Fprintln(w, "TOOL\tDESTINATION\tSCOPE")

		for _, name := range converter.SupportedTools {
			c, err := converter.Get(name)
			if err != nil {
				continue
			}

			scope := "user"
			if c.IsProjectScoped() {
				scope = "project"
			}

			fmt.Fprintf(w, "%s\t%s\t%s\n", name, c.Description(), scope)
		}

		if err := w.Flush(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	},
}

func init() { //nolint:gochecknoinits // required by cobra/converter
	rootCmd.AddCommand(toolsCmd)
}
