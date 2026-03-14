package converter

import (
	"os"
	"path/filepath"

	"github.com/caioreix/agency-cli/internal/agent"
)

type kimiCode struct{}

func init() { //nolint:gochecknoinits // required by cobra/converter
	Register("kimi-code", &kimiCode{})
}

func (c *kimiCode) Name() string          { return "Kimi Code" }
func (c *kimiCode) Description() string   { return ".kimi/agents/ + ~/.kimi/agents/" }
func (c *kimiCode) IsProjectScoped() bool { return true }

func (c *kimiCode) Convert(a *agent.Agent, _ string, scope string) ([]string, error) {
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
		dir = filepath.Join(home, ".kimi", "agents")
	default:
		dir = filepath.Join(cwd, ".kimi", "agents")
	}

	if mkdirErr := os.MkdirAll(dir, 0o755); mkdirErr != nil { //nolint:gosec // G301: world-traversable
		return nil, mkdirErr
	}

	systemFile := filepath.Join(dir, a.Slug+".md")
	if writeErr := os.WriteFile(systemFile, []byte(a.Body), 0o644); writeErr != nil { //nolint:gosec // G306: world-readable
		return nil, writeErr
	}

	yamlContent := "version: 1\n" +
		"agent:\n" +
		"  name: " + a.Name + "\n" +
		"  description: " + a.Description + "\n" +
		"  system_prompt_path: ./" + a.Slug + ".md\n" +
		"  extend: default\n"

	yamlFile := filepath.Join(dir, a.Slug+".yaml")
	if writeErr := os.WriteFile(yamlFile, []byte(yamlContent), 0o644); writeErr != nil { //nolint:gosec // G306: world-readable
		return nil, writeErr
	}

	return []string{yamlFile, systemFile}, nil
}
