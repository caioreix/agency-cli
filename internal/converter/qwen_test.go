package converter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQwen_Convert(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	a := newTestAgent()
	c, _ := Get("qwen")

	files, err := c.Convert(a, dir, ScopeLocal)
	require.NoError(t, err)
	require.Len(t, files, 1)
	assert.Equal(t, filepath.Join(dir, "test-agent.md"), files[0])

	content, err := os.ReadFile(files[0])
	require.NoError(t, err)
	assert.Contains(t, string(content), "name: test-agent")
	assert.Contains(t, string(content), "tools: bash,python")
}

func TestQwen_Convert_NoToolsField(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	a := newTestAgent()
	a.Tools = ""
	c, _ := Get("qwen")

	files, err := c.Convert(a, dir, ScopeLocal)
	require.NoError(t, err)

	content, err := os.ReadFile(files[0])
	require.NoError(t, err)
	assert.NotContains(t, string(content), "tools:")
}

func TestQwen_Convert_GlobalErrors(t *testing.T) {
	t.Parallel()
	c, _ := Get("qwen")
	_, err := c.Convert(newTestAgent(), t.TempDir(), ScopeGlobal)
	assert.Error(t, err)
}
