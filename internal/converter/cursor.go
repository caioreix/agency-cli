package converter

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/caioreix/agency-cli/internal/agent"
)

type cursor struct{}

func init() { //nolint:gochecknoinits // required by cobra/converter
	Register("cursor", &cursor{})
}

func (c *cursor) Name() string          { return "Cursor" }
func (c *cursor) Description() string   { return ".cursor/rules/ (project-scoped)" }
func (c *cursor) IsProjectScoped() bool { return true }

func (c *cursor) Convert(a *agent.Agent, destDir string, scope string) ([]string, error) {
	if scope == ScopeGlobal {
		return nil, errors.New("cursor is project-scoped; --scope global is not supported")
	}
	if err := os.MkdirAll(destDir, 0o755); err != nil { //nolint:gosec // G301: world-traversable
		return nil, err
	}

	// Cursor uses .mdc format with specific frontmatter
	outFile := filepath.Join(destDir, a.Slug+".mdc")
	content := "---\n" +
		"description: " + a.Description + "\n" +
		"globs: \"\"\n" +
		"alwaysApply: false\n" +
		"---\n" + a.Body

	if err := os.WriteFile(outFile, []byte(content), 0o644); err != nil { //nolint:gosec // G306: world-readable
		return nil, err
	}

	return []string{outFile}, nil
}
