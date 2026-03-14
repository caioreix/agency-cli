package converter

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/caioreix/agency-cli/internal/agent"
)

type opencode struct{}

func init() { //nolint:gochecknoinits // required by cobra/converter
	Register("opencode", &opencode{})
}

func (c *opencode) Name() string          { return "OpenCode" }
func (c *opencode) Description() string   { return ".opencode/agents/ + ~/.config/opencode/agents/" }
func (c *opencode) IsProjectScoped() bool { return true }

func (c *opencode) Convert(a *agent.Agent, _ string, scope string) ([]string, error) {
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
		dir = filepath.Join(home, ".config", "opencode", "agents")
	default:
		dir = filepath.Join(cwd, ".opencode", "agents")
	}

	if mkdirErr := os.MkdirAll(dir, 0o755); mkdirErr != nil { //nolint:gosec // G301: world-traversable
		return nil, mkdirErr
	}

	color := resolveOpenCodeColor(a.Color)

	outFile := filepath.Join(dir, a.Slug+".md")
	content := "---\n" +
		"name: " + a.Name + "\n" +
		"description: " + a.Description + "\n" +
		"mode: subagent\n" +
		"color: '" + color + "'\n" +
		"---\n" + a.Body

	writeErr := os.WriteFile(outFile, []byte(content), 0o644) //nolint:gosec // G306: world-readable
	if writeErr != nil {
		return nil, writeErr
	}

	return []string{outFile}, nil
}

func resolveOpenCodeColor(color string) string {
	c := strings.TrimSpace(strings.ToLower(color))
	colors := map[string]string{
		"cyan":          "#00FFFF",
		"blue":          "#3498DB",
		"green":         "#2ECC71",
		"red":           "#E74C3C",
		"purple":        "#9B59B6",
		"orange":        "#F39C12",
		"teal":          "#008080",
		"indigo":        "#6366F1",
		"pink":          "#E84393",
		"gold":          "#EAB308",
		"amber":         "#F59E0B",
		"neon-green":    "#10B981",
		"neon-cyan":     "#06B6D4",
		"metallic-blue": "#3B82F6",
		"yellow":        "#EAB308",
		"violet":        "#8B5CF6",
		"rose":          "#F43F5E",
		"lime":          "#84CC16",
		"gray":          "#6B7280",
		"fuchsia":       "#D946EF",
	}

	if hex, ok := colors[c]; ok {
		return hex
	}

	if len(c) == 7 && c[0] == '#' {
		return strings.ToUpper(c)
	}

	return "#6B7280"
}
