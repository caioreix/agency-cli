package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/caioreix/agency-cli/internal/agent"
	"github.com/caioreix/agency-cli/internal/color"
	"github.com/caioreix/agency-cli/internal/repo"
	rw "github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
)

var categoryFlag string

const colPadding = 2

var ansiStripper = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// dispWidth returns the terminal display width of a string,
// correctly handling ANSI escape codes and multi-codepoint sequences.
// Uses per-rune RuneWidth so that:
//   - flag emoji (e.g. 🇨🇳 = two regional-indicator runes, each width=1) → 2
//   - VS16 sequences (e.g. 🏗️ = base rune width=1 + U+FE0F width=0) → 1
//   - wide emoji (e.g. 🧬 = single rune width=2) → 2
func dispWidth(s string) int {
	clean := ansiStripper.ReplaceAllString(s, "")
	width := 0
	for _, r := range clean {
		width += rw.RuneWidth(r)
	}
	return width
}

// pad right-pads s to the given display width using spaces.
func pad(s string, width int) string {
	d := width - dispWidth(s)
	if d <= 0 {
		return s
	}
	return s + strings.Repeat(" ", d)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available agents",
	Long:  "List all available agents from the agency-agents repository. Use --category to filter by category.",
	RunE: func(_ *cobra.Command, _ []string) error {
		repoDir, err := repo.EnsureRepo()
		if err != nil {
			return fmt.Errorf("failed to ensure repo: %w", err)
		}

		agents, err := agent.ListAll(repoDir)
		if err != nil {
			return fmt.Errorf("failed to list agents: %w", err)
		}

		type row struct {
			category   string
			name       string
			vibe       string
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

		// Compute max display widths (plain text, no ANSI)
		w0 := dispWidth("CATEGORY")
		w1 := dispWidth("AGENT")
		for _, r := range rows {
			if n := dispWidth(r.name); n > w1 {
				w1 = n
			}
			if n := dispWidth(r.category); n > w0 {
				w0 = n
			}
		}

		const maxVibe = 60

		divider := color.ApplyDim(
			strings.Repeat("─", w0+colPadding) + "┼" +
				strings.Repeat("─", w1+colPadding) + "┼" +
				strings.Repeat("─", maxVibe+colPadding),
		)

		// Header
		fmt.Fprintf(os.Stdout, " %s │ %s │ %s\n",
			pad(color.ApplyBold("CATEGORY"), w0),
			pad(color.ApplyBold("AGENT"), w1),
			color.ApplyBold("VIBE"),
		)
		fmt.Fprintln(os.Stdout, divider)

		prevCat := ""
		for _, r := range rows {
			if r.category != prevCat && prevCat != "" {
				fmt.Fprintln(os.Stdout, divider)
			}
			prevCat = r.category

			vibe := r.vibe
			if dispWidth(vibe) > maxVibe {
				runes := []rune(vibe)
				w := 0
				cut := 0
				for i, ch := range runes {
					w += rw.RuneWidth(ch)
					if w > maxVibe-3 {
						cut = i
						break
					}
				}
				vibe = string(runes[:cut]) + "..."
			}

			nameColored := color.Apply(r.name, r.agentColor)
			catDim := color.ApplyDim(r.category)

			fmt.Fprintf(os.Stdout, " %s │ %s │ %s\n",
				pad(catDim, w0),
				pad(nameColored, w1),
				vibe,
			)
		}

		return nil
	},
}

func init() { //nolint:gochecknoinits // required by cobra/converter
	listCmd.Flags().StringVarP(&categoryFlag, "category", "c", "", "filter by category (e.g., engineering, design)")
	rootCmd.AddCommand(listCmd)
}
