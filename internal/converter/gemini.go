package converter

import (
	"os"
	"path/filepath"

	"github.com/caioreix/agency-cli/internal/agent"
)

type geminiCLI struct{}

func init() { //nolint:gochecknoinits // required by cobra/converter
	Register("gemini-cli", &geminiCLI{})
}

func (c *geminiCLI) Name() string          { return "Gemini CLI" }
func (c *geminiCLI) Description() string   { return ".gemini/extensions/ + ~/.gemini/extensions/" }
func (c *geminiCLI) IsProjectScoped() bool { return true }

func (c *geminiCLI) Convert(a *agent.Agent, _ string, scope string) ([]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var baseDir string
	switch scope {
	case ScopeGlobal:
		baseDir = filepath.Join(home, ".gemini", "extensions", "agency-agents")
	default:
		baseDir = filepath.Join(cwd, ".gemini", "extensions", "agency-agents")
	}

	skillDir := filepath.Join(baseDir, "skills", a.Slug)
	if mkdirErr := os.MkdirAll(skillDir, 0o755); mkdirErr != nil { //nolint:gosec // G301: world-traversable
		return nil, mkdirErr
	}

	outFile := filepath.Join(skillDir, "SKILL.md")
	content := "---\n" +
		"name: " + a.Slug + "\n" +
		"description: " + a.Description + "\n" +
		"---\n" + a.Body

	writeErr := os.WriteFile(outFile, []byte(content), 0o644) //nolint:gosec // G306: world-readable
	if writeErr != nil {
		return nil, writeErr
	}

	// Ensure extension manifest exists
	manifestFile := filepath.Join(baseDir, "gemini-extension.json")
	if _, statErr := os.Stat(manifestFile); os.IsNotExist(statErr) {
		manifest := `{
  "name": "agency-agents",
  "version": "1.0.0"
}
`
		if manifestErr := os.WriteFile( //nolint:gosec // G306: world-readable
			manifestFile,
			[]byte(manifest),
			0o644,
		); manifestErr != nil {
			return nil, manifestErr
		}
	}

	return []string{outFile}, nil
}
