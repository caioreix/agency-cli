package converter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClaudeCode_Convert(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	a := newTestAgent()

	c, err := Get("claude-code")
	require.NoError(t, err)

	files, err := c.Convert(a, dir, ScopeLocal)
	require.NoError(t, err)
	require.Len(t, files, 1)

	content, err := os.ReadFile(files[0])
	require.NoError(t, err)

	assert.Equal(t, filepath.Join(dir, "test-agent.md"), files[0])
	assert.Contains(t, string(content), "name: Test Agent")
	assert.Contains(t, string(content), "description: A test agent for unit tests")
	assert.Contains(t, string(content), "color: cyan")
	assert.Contains(t, string(content), "emoji: 🤖")
	assert.Contains(t, string(content), "## Mission")
}

func TestClaudeCode_Convert_IgnoresScope(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	c, _ := Get("claude-code")

	_, err := c.Convert(newTestAgent(), dir, ScopeGlobal)
	assert.NoError(t, err)
}
