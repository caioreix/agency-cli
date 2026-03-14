package converter

import (
	"os"
	"path/filepath"

	"github.com/caioreix/agency-cli/internal/agent"
)

type copilot struct{}

func init() { //nolint:gochecknoinits // required by cobra/converter
	Register("copilot", &copilot{})
}

func (c *copilot) Name() string          { return "Copilot" }
func (c *copilot) Description() string   { return ".github/agents/ + ~/.copilot/agents/" }
func (c *copilot) IsProjectScoped() bool { return true }

func (c *copilot) Convert(a *agent.Agent, _ string, scope string) ([]string, error) {
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

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var dirs []string
	switch scope {
	case ScopeGlobal:
		dirs = []string{filepath.Join(home, ".copilot", "agents")}
	default:
		dirs = []string{filepath.Join(cwd, ".github", "agents")}
	}

	var files []string
	for _, dir := range dirs {
		if mkdirErr := os.MkdirAll(dir, 0o755); mkdirErr != nil { //nolint:gosec // G301: world-traversable
			return nil, mkdirErr
		}
		outFile := filepath.Join(dir, a.Slug+".md")
		if writeErr := os.WriteFile( //nolint:gosec // G306: world-readable
			outFile,
			[]byte(content),
			0o644,
		); writeErr != nil {
			return nil, writeErr
		}
		files = append(files, outFile)
	}

	return files, nil
}
