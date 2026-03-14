package converter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/caioreix/agency-cli/internal/agent"
)

type opencode struct{}

func init() {
	Register("opencode", &opencode{})
}

func (c *opencode) Name() string          { return "OpenCode" }
func (c *opencode) Description() string   { return ".opencode/agents/ (project-scoped)" }
func (c *opencode) IsProjectScoped() bool { return true }

func (c *opencode) Convert(a *agent.Agent, destDir string, scope string) ([]string, error) {
	if scope == ScopeGlobal {
		return nil, fmt.Errorf("opencode is project-scoped; --scope global is not supported")
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return nil, err
	}

	color := resolveOpenCodeColor(a.Color)

	outFile := filepath.Join(destDir, a.Slug+".md")
	content := "---\n" +
		"name: " + a.Name + "\n" +
		"description: " + a.Description + "\n" +
		"mode: subagent\n" +
		"color: '" + color + "'\n" +
		"---\n" + a.Body

	if err := os.WriteFile(outFile, []byte(content), 0o644); err != nil {
		return nil, err
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
