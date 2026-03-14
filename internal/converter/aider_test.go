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

func TestAider_Convert_CreatesConventionsFile(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	a := newTestAgent()
	c, _ := Get("aider")

	files, err := c.Convert(a, dir, ScopeLocal)
	require.NoError(t, err)
	require.Len(t, files, 1)
	assert.Equal(t, filepath.Join(dir, "CONVENTIONS.md"), files[0])

	content, err := os.ReadFile(files[0])
	require.NoError(t, err)
	assert.Contains(t, string(content), "Test Agent")
	assert.Contains(t, string(content), "A test agent for unit tests")
}

func TestAider_Convert_AppendsOnSecondCall(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	c, _ := Get("aider")

	a1 := newTestAgent()
	a2 := &agent.Agent{Name: "Second", Slug: "second", Description: "Desc", Body: "body\n"}

	_, err := c.Convert(a1, dir, ScopeLocal)
	require.NoError(t, err)
	_, err = c.Convert(a2, dir, ScopeLocal)
	require.NoError(t, err)

	content, err := os.ReadFile(filepath.Join(dir, "CONVENTIONS.md"))
	require.NoError(t, err)
	assert.Contains(t, string(content), "Test Agent")
	assert.Contains(t, string(content), "Second")
}

func TestAider_Convert_GlobalErrors(t *testing.T) {
	t.Parallel()
	c, _ := Get("aider")
	_, err := c.Convert(newTestAgent(), t.TempDir(), ScopeGlobal)
	assert.Error(t, err)
}
