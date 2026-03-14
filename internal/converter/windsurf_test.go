//nolint:testpackage // shares newTestAgent helper and tests unexported functions
package converter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/caioreix/agency-cli/internal/agent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWindsurf_Convert_CreatesRulesFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	a := newTestAgent()
	c, _ := Get("windsurf")

	files, err := c.Convert(a, dir, ScopeLocal)
	require.NoError(t, err)
	require.Len(t, files, 1)
	assert.Equal(t, filepath.Join(dir, ".windsurfrules"), files[0])

	content, err := os.ReadFile(files[0])
	require.NoError(t, err)
	assert.Contains(t, string(content), "Test Agent")
	assert.Contains(t, string(content), "A test agent for unit tests")
}

func TestWindsurf_Convert_AppendsOnSecondCall(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	c, _ := Get("windsurf")

	a1 := newTestAgent()
	a2 := &agent.Agent{Name: "Second Agent", Slug: "second-agent", Description: "Second", Body: "second body\n"}

	_, err := c.Convert(a1, dir, ScopeLocal)
	require.NoError(t, err)
	_, err = c.Convert(a2, dir, ScopeLocal)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(dir, ".windsurfrules"))
	require.NoError(t, err)
	assert.Contains(t, string(content), "Test Agent")
	assert.Contains(t, string(content), "Second Agent")
}

func TestWindsurf_Convert_GlobalErrors(t *testing.T) {
	t.Parallel()
	c, _ := Get("windsurf")
	_, err := c.Convert(newTestAgent(), t.TempDir(), ScopeGlobal)
	assert.Error(t, err)
}
