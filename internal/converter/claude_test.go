//nolint:testpackage // shares newTestAgent helper and tests unexported functions
package converter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClaudeCode_Convert_Local uses t.Chdir which is incompatible with t.Parallel().
func TestClaudeCode_Convert_Local(t *testing.T) {
	t.Chdir(t.TempDir())

	cwd, err := os.Getwd()
	require.NoError(t, err)

	a := newTestAgent()
	c, err := Get("claude-code")
	require.NoError(t, err)

	files, err := c.Convert(a, "", ScopeLocal)
	require.NoError(t, err)
	require.Len(t, files, 1)

	expected := filepath.Join(cwd, ".claude", "agents", "test-agent.md")
	assert.Equal(t, expected, files[0])

	content, err := os.ReadFile(files[0])
	require.NoError(t, err)
	assert.Contains(t, string(content), "name: Test Agent")
	assert.Contains(t, string(content), "description: A test agent for unit tests")
	assert.Contains(t, string(content), "color: cyan")
	assert.Contains(t, string(content), "emoji: 🤖")
	assert.Contains(t, string(content), "## Mission")
}

// TestClaudeCode_Convert_Default uses t.Chdir which is incompatible with t.Parallel().
func TestClaudeCode_Convert_Default(t *testing.T) {
	t.Chdir(t.TempDir())

	cwd, err := os.Getwd()
	require.NoError(t, err)

	a := newTestAgent()
	c, _ := Get("claude-code")

	files, err := c.Convert(a, "", ScopeDefault)
	require.NoError(t, err)
	require.Len(t, files, 1)

	expected := filepath.Join(cwd, ".claude", "agents", "test-agent.md")
	assert.Equal(t, expected, files[0])
}

func TestClaudeCode_Convert_Global(t *testing.T) {
	t.Parallel()
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	a := newTestAgent()
	c, _ := Get("claude-code")

	files, err := c.Convert(a, "", ScopeGlobal)
	require.NoError(t, err)
	require.Len(t, files, 1)

	expected := filepath.Join(home, ".claude", "agents", "test-agent.md")
	assert.Equal(t, expected, files[0])

	t.Cleanup(func() { os.Remove(files[0]) })
}
