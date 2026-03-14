package converter

import (
	"os"
	"path/filepath"

	"github.com/caioreix/agency-cli/internal/agent"
)

type claudeCode struct{}

func init() { //nolint:gochecknoinits // required by cobra/converter
	Register("claude-code", &claudeCode{})
}

func (c *claudeCode) Name() string          { return "Claude Code" }
func (c *claudeCode) Description() string   { return ".claude/agents/ + ~/.claude/agents/" }
func (c *claudeCode) IsProjectScoped() bool { return true }

func (c *claudeCode) Convert(a *agent.Agent, _ string, scope string) ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var dir string
	switch scope {
	case ScopeGlobal:
		dir = filepath.Join(home, ".claude", "agents")
	default:
		dir = filepath.Join(cwd, ".claude", "agents")
	}

	if mkdirErr := os.MkdirAll(dir, 0o755); mkdirErr != nil { //nolint:gosec // G301: world-traversable
		return nil, mkdirErr
	}

	// Claude Code uses .md files with YAML frontmatter
	outFile := filepath.Join(dir, a.Slug+".md")
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

	writeErr := os.WriteFile(outFile, []byte(content), 0o644) //nolint:gosec // G306: world-readable
	if writeErr != nil {
		return nil, writeErr
	}

	return []string{outFile}, nil
}
