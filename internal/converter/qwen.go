package converter

import (
	"os"
	"path/filepath"

	"github.com/caioreix/agency-cli/internal/agent"
)

type qwen struct{}

func init() {
	Register("qwen", &qwen{})
}

func (c *qwen) Name() string        { return "Qwen Code" }
func (c *qwen) Description() string  { return ".qwen/agents/ (project-scoped)" }
func (c *qwen) IsProjectScoped() bool { return true }

func (c *qwen) Convert(a *agent.Agent, destDir string) ([]string, error) {
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return nil, err
	}

	outFile := filepath.Join(destDir, a.Slug+".md")

	content := "---\n" +
		"name: " + a.Slug + "\n" +
		"description: " + a.Description + "\n"
	if a.Tools != "" {
		content += "tools: " + a.Tools + "\n"
	}
	content += "---\n" + a.Body

	if err := os.WriteFile(outFile, []byte(content), 0o644); err != nil {
		return nil, err
	}

	return []string{outFile}, nil
}
