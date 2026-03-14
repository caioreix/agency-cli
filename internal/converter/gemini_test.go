package converter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/caioreix/agency-cli/internal/agent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGeminiCLI_Convert_Local uses t.Chdir which is incompatible with t.Parallel().
func TestGeminiCLI_Convert_Local(t *testing.T) {
	t.Chdir(t.TempDir())

	cwd, err := os.Getwd()
	require.NoError(t, err)

	a := newTestAgent()
	c, _ := Get("gemini-cli")

	files, err := c.Convert(a, "", ScopeLocal)
	require.NoError(t, err)
	require.Len(t, files, 1)

	skillFile := filepath.Join(cwd, ".gemini", "extensions", "agency-agents", "skills", "test-agent", "SKILL.md")
	assert.Equal(t, skillFile, files[0])

	content, err := os.ReadFile(files[0])
	require.NoError(t, err)
	assert.Contains(t, string(content), "name: test-agent")

	manifest, err := os.ReadFile(filepath.Join(cwd, ".gemini", "extensions", "agency-agents", "gemini-extension.json"))
	require.NoError(t, err)
	assert.Contains(t, string(manifest), "agency-agents")
}

func TestGeminiCLI_Convert_Global(t *testing.T) {
	t.Parallel()
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	a := newTestAgent()
	c, _ := Get("gemini-cli")

	files, err := c.Convert(a, "", ScopeGlobal)
	require.NoError(t, err)
	require.Len(t, files, 1)

	skillFile := filepath.Join(home, ".gemini", "extensions", "agency-agents", "skills", "test-agent", "SKILL.md")
	assert.Equal(t, skillFile, files[0])

	t.Cleanup(func() { os.Remove(files[0]) })
}

// TestGeminiCLI_Convert_ManifestNotOverwritten uses t.Chdir which is incompatible with t.Parallel().
func TestGeminiCLI_Convert_ManifestNotOverwritten(t *testing.T) {
	t.Chdir(t.TempDir())

	cwd, err := os.Getwd()
	require.NoError(t, err)

	a1 := newTestAgent()
	a2 := &agent.Agent{Name: "Second", Slug: "second", Description: "Desc", Body: "body\n"}
	c, _ := Get("gemini-cli")

	_, err = c.Convert(a1, "", ScopeLocal)
	require.NoError(t, err)

	manifestPath := filepath.Join(cwd, ".gemini", "extensions", "agency-agents", "gemini-extension.json")
	customManifest := `{"name":"custom","version":"2.0.0"}`
	require.NoError(t, os.WriteFile(manifestPath, []byte(customManifest), 0o644))

	_, err = c.Convert(a2, "", ScopeLocal)
	require.NoError(t, err)

	manifest, err := os.ReadFile(manifestPath)
	require.NoError(t, err)
	assert.Equal(t, customManifest, strings.TrimSpace(string(manifest)))
}
