package converter

import (
	"os"
	"path/filepath"

	"github.com/caioreix/agency-cli/internal/agent"
)

type copilot struct{}

func init() {
	Register("copilot", &copilot{})
}

func (c *copilot) Name() string        { return "Copilot" }
func (c *copilot) Description() string  { return "~/.github/agents/ + ~/.copilot/agents/" }
func (c *copilot) IsProjectScoped() bool { return false }

func (c *copilot) Convert(a *agent.Agent, destDir string) ([]string, error) {
	// Copilot uses original .md files with frontmatter, copied to two locations
	content := "---\n" +
		"name: " + a.Name + "\n" +
		"description: " + a.Description + "\n"
	if a.Color != "" {
		content += "color: " + a.Color + "\n"
	}
	if a.Emoji != "" {
		content += "emoji: " + a.Emoji + "\n"
	}
	if a.Vibe != "" {
		content += "vibe: " + a.Vibe + "\n"
	}
	content += "---\n" + a.Body

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dirs := []string{
		filepath.Join(home, ".github", "agents"),
		filepath.Join(home, ".copilot", "agents"),
	}

	var files []string
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
		outFile := filepath.Join(dir, a.Slug+".md")
		if err := os.WriteFile(outFile, []byte(content), 0o644); err != nil {
			return nil, err
		}
		files = append(files, outFile)
	}

	return files, nil
}
