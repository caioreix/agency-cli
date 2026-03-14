package converter

import (
	"os"
	"path/filepath"

	"github.com/caioreix/agency-cli/internal/agent"
)

type geminiCLI struct{}

func init() {
	Register("gemini-cli", &geminiCLI{})
}

func (c *geminiCLI) Name() string        { return "Gemini CLI" }
func (c *geminiCLI) Description() string  { return "~/.gemini/extensions/agency-agents/" }
func (c *geminiCLI) IsProjectScoped() bool { return false }

func (c *geminiCLI) Convert(a *agent.Agent, destDir string, scope string) ([]string, error) {
	// gemini-cli installs globally regardless of scope
	skillDir := filepath.Join(destDir, "skills", a.Slug)
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		return nil, err
	}

	outFile := filepath.Join(skillDir, "SKILL.md")
	content := "---\n" +
		"name: " + a.Slug + "\n" +
		"description: " + a.Description + "\n" +
		"---\n" + a.Body

	if err := os.WriteFile(outFile, []byte(content), 0o644); err != nil {
		return nil, err
	}

	// Ensure extension manifest exists
	manifestFile := filepath.Join(destDir, "gemini-extension.json")
	if _, err := os.Stat(manifestFile); os.IsNotExist(err) {
		manifest := `{
  "name": "agency-agents",
  "version": "1.0.0"
}
`
		if err := os.WriteFile(manifestFile, []byte(manifest), 0o644); err != nil {
			return nil, err
		}
	}

	return []string{outFile}, nil
}
