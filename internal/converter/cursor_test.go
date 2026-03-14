package converter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCursor_Convert_Local(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	a := newTestAgent()
	c, _ := Get("cursor")

	files, err := c.Convert(a, dir, ScopeLocal)
	require.NoError(t, err)
	require.Len(t, files, 1)

	assert.Equal(t, filepath.Join(dir, "test-agent.mdc"), files[0])

	content, err := os.ReadFile(files[0])
	require.NoError(t, err)
	assert.Contains(t, string(content), "description: A test agent for unit tests")
	assert.Contains(t, string(content), "alwaysApply: false")
}

func TestCursor_Convert_GlobalErrors(t *testing.T) {
	t.Parallel()
	c, _ := Get("cursor")
	_, err := c.Convert(newTestAgent(), t.TempDir(), ScopeGlobal)
	assert.Error(t, err)
}
