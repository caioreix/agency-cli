package converter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestKimiCode_Convert_Local uses t.Chdir which is incompatible with t.Parallel().
func TestKimiCode_Convert_Local(t *testing.T) {
	t.Chdir(t.TempDir())

	cwd, err := os.Getwd()
	require.NoError(t, err)

	a := newTestAgent()
	c, _ := Get("kimi-code")

	files, err := c.Convert(a, "", ScopeLocal)
	require.NoError(t, err)
	require.Len(t, files, 2)

	expectedYAML := filepath.Join(cwd, ".kimi", "agents", "test-agent.yaml")
	expectedMD := filepath.Join(cwd, ".kimi", "agents", "test-agent.md")
	assert.Equal(t, expectedYAML, files[0])
	assert.Equal(t, expectedMD, files[1])

	yamlContent, err := os.ReadFile(files[0])
	require.NoError(t, err)
	assert.Contains(t, string(yamlContent), "name: Test Agent")
	assert.Contains(t, string(yamlContent), "description: A test agent for unit tests")
	assert.Contains(t, string(yamlContent), "system_prompt_path: ./test-agent.md")
	assert.Contains(t, string(yamlContent), "extend: default")

	mdContent, err := os.ReadFile(files[1])
	require.NoError(t, err)
	assert.Contains(t, string(mdContent), "## Mission")
}

// TestKimiCode_Convert_Default uses t.Chdir which is incompatible with t.Parallel().
func TestKimiCode_Convert_Default(t *testing.T) {
	t.Chdir(t.TempDir())

	cwd, err := os.Getwd()
	require.NoError(t, err)

	a := newTestAgent()
	c, _ := Get("kimi-code")

	files, err := c.Convert(a, "", ScopeDefault)
	require.NoError(t, err)
	require.Len(t, files, 2)

	expectedYAML := filepath.Join(cwd, ".kimi", "agents", "test-agent.yaml")
	assert.Equal(t, expectedYAML, files[0])
}

func TestKimiCode_Convert_Global(t *testing.T) {
	t.Parallel()
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	a := newTestAgent()
	c, _ := Get("kimi-code")

	files, err := c.Convert(a, "", ScopeGlobal)
	require.NoError(t, err)
	require.Len(t, files, 2)

	expectedYAML := filepath.Join(home, ".kimi", "agents", "test-agent.yaml")
	expectedMD := filepath.Join(home, ".kimi", "agents", "test-agent.md")
	assert.Equal(t, expectedYAML, files[0])
	assert.Equal(t, expectedMD, files[1])

	t.Cleanup(func() {
		os.Remove(files[0])
		os.Remove(files[1])
	})
}
