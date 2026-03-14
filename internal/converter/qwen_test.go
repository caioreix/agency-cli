//nolint:testpackage // shares newTestAgent helper and tests unexported functions
package converter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestQwen_Convert_Local uses t.Chdir which is incompatible with t.Parallel().
func TestQwen_Convert_Local(t *testing.T) {
	t.Chdir(t.TempDir())

	cwd, err := os.Getwd()
	require.NoError(t, err)

	a := newTestAgent()
	c, _ := Get("qwen")

	files, err := c.Convert(a, "", ScopeLocal)
	require.NoError(t, err)
	require.Len(t, files, 1)

	expected := filepath.Join(cwd, ".qwen", "agents", "test-agent.md")
	assert.Equal(t, expected, files[0])

	content, err := os.ReadFile(files[0])
	require.NoError(t, err)
	assert.Contains(t, string(content), "name: test-agent")
	assert.Contains(t, string(content), "tools: bash,python")
}

// TestQwen_Convert_NoToolsField uses t.Chdir which is incompatible with t.Parallel().
func TestQwen_Convert_NoToolsField(t *testing.T) {
	t.Chdir(t.TempDir())

	a := newTestAgent()
	a.Tools = ""
	c, _ := Get("qwen")

	files, err := c.Convert(a, "", ScopeLocal)
	require.NoError(t, err)

	content, err := os.ReadFile(files[0])
	require.NoError(t, err)
	assert.NotContains(t, string(content), "tools:")
}

func TestQwen_Convert_Global(t *testing.T) {
	t.Parallel()
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	a := newTestAgent()
	c, _ := Get("qwen")

	files, err := c.Convert(a, "", ScopeGlobal)
	require.NoError(t, err)
	require.Len(t, files, 1)

	expected := filepath.Join(home, ".qwen", "agents", "test-agent.md")
	assert.Equal(t, expected, files[0])

	t.Cleanup(func() { os.Remove(files[0]) })
}
