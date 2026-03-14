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

func TestGeminiCLI_Convert(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	a := newTestAgent()
	c, _ := Get("gemini-cli")

	files, err := c.Convert(a, dir, ScopeLocal)
	require.NoError(t, err)
	require.Len(t, files, 1)

	skillFile := filepath.Join(dir, "skills", "test-agent", "SKILL.md")
	assert.Equal(t, skillFile, files[0])

	content, err := os.ReadFile(files[0])
	require.NoError(t, err)
	assert.Contains(t, string(content), "name: test-agent")

	manifest, err := os.ReadFile(filepath.Join(dir, "gemini-extension.json"))
	require.NoError(t, err)
	assert.Contains(t, string(manifest), "agency-agents")
}

func TestGeminiCLI_Convert_ManifestNotOverwritten(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	a1 := newTestAgent()
	a2 := &agent.Agent{Name: "Second", Slug: "second", Description: "Desc", Body: "body\n"}
	c, _ := Get("gemini-cli")

	_, err := c.Convert(a1, dir, ScopeLocal)
	require.NoError(t, err)

	customManifest := `{"name":"custom","version":"2.0.0"}`
	require.NoError(t, os.WriteFile(filepath.Join(dir, "gemini-extension.json"), []byte(customManifest), 0o644))

	_, err = c.Convert(a2, dir, ScopeLocal)
	require.NoError(t, err)

	manifest, err := os.ReadFile(filepath.Join(dir, "gemini-extension.json"))
	require.NoError(t, err)
	assert.Equal(t, customManifest, strings.TrimSpace(string(manifest)))
}
