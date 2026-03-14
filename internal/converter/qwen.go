package converter

import (
	"os"
	"path/filepath"

	"github.com/caioreix/agency-cli/internal/agent"
)

type qwen struct{}

func init() { //nolint:gochecknoinits // required by cobra/converter
	Register("qwen", &qwen{})
}

func (c *qwen) Name() string          { return "Qwen Code" }
func (c *qwen) Description() string   { return ".qwen/agents/ + ~/.qwen/agents/" }
func (c *qwen) IsProjectScoped() bool { return true }

func (c *qwen) Convert(a *agent.Agent, _ string, scope string) ([]string, error) {
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
		dir = filepath.Join(home, ".qwen", "agents")
	default:
		dir = filepath.Join(cwd, ".qwen", "agents")
	}

	if mkdirErr := os.MkdirAll(dir, 0o755); mkdirErr != nil { //nolint:gosec // G301: world-traversable
		return nil, mkdirErr
	}

	outFile := filepath.Join(dir, a.Slug+".md")

	content := "---\n" +
		"name: " + a.Slug + "\n" +
		"description: " + a.Description + "\n"
	if a.Tools != "" {
		content += "tools: " + a.Tools + "\n"
	}
	content += "---\n" + a.Body

	if writeErr := os.WriteFile(outFile, []byte(content), 0o644); writeErr != nil { //nolint:gosec // G306: world-readable
		return nil, writeErr
	}

	return []string{outFile}, nil
}
