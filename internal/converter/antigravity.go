package converter

import (
	"os"
	"path/filepath"

	"github.com/caioreix/agency-cli/internal/agent"
)

type antigravity struct{}

func init() {
	Register("antigravity", &antigravity{})
}

func (c *antigravity) Name() string        { return "Antigravity" }
func (c *antigravity) Description() string  { return "~/.gemini/antigravity/skills/" }
func (c *antigravity) IsProjectScoped() bool { return false }

func (c *antigravity) Convert(a *agent.Agent, destDir string, scope string) ([]string, error) {
	// antigravity installs globally regardless of scope
	slug := "agency-" + a.Slug
	skillDir := filepath.Join(destDir, slug)
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		return nil, err
	}

	outFile := filepath.Join(skillDir, "SKILL.md")
	content := "---\n" +
		"name: " + slug + "\n" +
		"description: " + a.Description + "\n" +
		"risk: low\n" +
		"source: community\n" +
		"---\n" + a.Body

	if err := os.WriteFile(outFile, []byte(content), 0o644); err != nil {
		return nil, err
	}

	return []string{outFile}, nil
}
