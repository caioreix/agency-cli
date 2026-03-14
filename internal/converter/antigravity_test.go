//nolint:testpackage // shares newTestAgent helper and tests unexported functions
package converter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAntigravity_Convert(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	a := newTestAgent()
	c, _ := Get("antigravity")

	files, err := c.Convert(a, dir, ScopeLocal)
	require.NoError(t, err)
	require.Len(t, files, 1)

	skillDir := filepath.Join(dir, "agency-test-agent")
	assert.Equal(t, filepath.Join(skillDir, "SKILL.md"), files[0])

	content, err := os.ReadFile(files[0])
	require.NoError(t, err)
	assert.Contains(t, string(content), "name: agency-test-agent")
	assert.Contains(t, string(content), "risk: low")
}
