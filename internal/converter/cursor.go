package converter

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/caioreix/agency-cli/internal/agent"
)

type cursor struct{}

func init() {
	Register("cursor", &cursor{})
}

func (c *cursor) Name() string          { return "Cursor" }
func (c *cursor) Description() string   { return ".cursor/rules/ (project-scoped)" }
func (c *cursor) IsProjectScoped() bool { return true }

func (c *cursor) Convert(a *agent.Agent, destDir string, scope string) ([]string, error) {
	if scope == ScopeGlobal {
		return nil, fmt.Errorf("cursor is project-scoped; --scope global is not supported")
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return nil, err
	}

	// Cursor uses .mdc format with specific frontmatter
	outFile := filepath.Join(destDir, a.Slug+".mdc")
	content := "---\n" +
		"description: " + a.Description + "\n" +
		"globs: \"\"\n" +
		"alwaysApply: false\n" +
		"---\n" + a.Body

	if err := os.WriteFile(outFile, []byte(content), 0o644); err != nil {
		return nil, err
	}

	return []string{outFile}, nil
}
