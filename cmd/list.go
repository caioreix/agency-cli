package cmd

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/caioreix/agency-cli/internal/agent"
	"github.com/caioreix/agency-cli/internal/color"
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

		type row struct {
			category string
			name     string
			vibe     string
			agentColor string
		}

		var rows []row
		for _, a := range agents {
			if categoryFlag != "" && a.Category != categoryFlag {
				continue
			}

			name := a.Name
			if a.Emoji != "" {
				name = a.Emoji + " " + name
			}

			vibe := a.Vibe
			if vibe == "" {
				vibe = a.Description
			}

			rows = append(rows, row{
				category:   a.Category,
				name:       name,
				vibe:       vibe,
				agentColor: a.Color,
			})
		}

		// Compute visible column widths (no ANSI codes)
		w0, w1 := utf8.RuneCountInString("CATEGORY"), utf8.RuneCountInString("AGENT")
		for _, r := range rows {
			if n := utf8.RuneCountInString(r.category); n > w0 {
				w0 = n
			}
			if n := utf8.RuneCountInString(r.name); n > w1 {
				w1 = n
			}
		}

		// Cap vibe width to keep lines reasonable
		const maxVibe = 60

		sep := color.ApplyDim(strings.Repeat("─", w0+2) + "┼" + strings.Repeat("─", w1+2) + "┼" + strings.Repeat("─", maxVibe+2))

		// Header
		header := fmt.Sprintf(" %-*s │ %-*s │ %s",
			w0, color.ApplyBold("CATEGORY"),
			w1, color.ApplyBold("AGENT"),
			color.ApplyBold("VIBE"),
		)
		fmt.Println(header)
		fmt.Println(sep)

		// Rows
		prevCat := ""
		for _, r := range rows {
			// Add separator between categories
			if r.category != prevCat && prevCat != "" {
				fmt.Println(sep)
			}
			prevCat = r.category

			vibe := r.vibe
			if utf8.RuneCountInString(vibe) > maxVibe {
				runes := []rune(vibe)
				vibe = string(runes[:maxVibe-3]) + "..."
			}

			nameColored := color.Apply(r.name, r.agentColor)
			catDim := color.ApplyDim(r.category)

			fmt.Printf(" %-*s │ %-*s │ %s\n",
				w0+visiblePadding(r.category, catDim),
				catDim,
				w1+visiblePadding(r.name, nameColored),
				nameColored,
				vibe,
			)
		}

		return nil
	},
}

// visiblePadding returns extra padding needed to compensate for invisible ANSI codes.
func visiblePadding(plain, withCodes string) int {
	return len(withCodes) - len(plain)
}

func init() {
	listCmd.Flags().StringVarP(&categoryFlag, "category", "c", "", "filter by category (e.g., engineering, design)")
	rootCmd.AddCommand(listCmd)
}
