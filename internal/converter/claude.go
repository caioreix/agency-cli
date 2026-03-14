package converter

import (
	"os"
	"path/filepath"

	"github.com/caioreix/agency-cli/internal/agent"
)

type claudeCode struct{}

func init() {
	Register("claude-code", &claudeCode{})
}

func (c *claudeCode) Name() string        { return "Claude Code" }
func (c *claudeCode) Description() string  { return "~/.claude/agents/" }
func (c *claudeCode) IsProjectScoped() bool { return false }

func (c *claudeCode) Convert(a *agent.Agent, destDir string, scope string) ([]string, error) {
	// claude-code installs globally regardless of scope
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return nil, err
	}

	// Claude Code uses the original .md files with frontmatter as-is
	outFile := filepath.Join(destDir, a.Slug+".md")
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

	if err := os.WriteFile(outFile, []byte(content), 0o644); err != nil {
		return nil, err
	}

	return []string{outFile}, nil
}
