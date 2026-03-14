package converter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestOpenCode_Convert_Local uses t.Chdir which is incompatible with t.Parallel().
func TestOpenCode_Convert_Local(t *testing.T) {
	t.Chdir(t.TempDir())

	cwd, err := os.Getwd()
	require.NoError(t, err)

	a := newTestAgent()
	c, _ := Get("opencode")

	files, err := c.Convert(a, "", ScopeLocal)
	require.NoError(t, err)
	require.Len(t, files, 1)

	expected := filepath.Join(cwd, ".opencode", "agents", "test-agent.md")
	assert.Equal(t, expected, files[0])

	content, err := os.ReadFile(files[0])
	require.NoError(t, err)
	assert.Contains(t, string(content), "name: Test Agent")
	assert.Contains(t, string(content), "mode: subagent")
	assert.Contains(t, string(content), "color:")
}

func TestOpenCode_Convert_Global(t *testing.T) {
	t.Parallel()
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	a := newTestAgent()
	c, _ := Get("opencode")

	files, err := c.Convert(a, "", ScopeGlobal)
	require.NoError(t, err)
	require.Len(t, files, 1)

	expected := filepath.Join(home, ".config", "opencode", "agents", "test-agent.md")
	assert.Equal(t, expected, files[0])

	t.Cleanup(func() { os.Remove(files[0]) })
}

// ── resolveOpenCodeColor ──────────────────────────────────────────────────────

func TestResolveOpenCodeColor_NamedColors(t *testing.T) {
	t.Parallel()
	tests := []struct {
		input string
		want  string
	}{
		{"cyan", "#00FFFF"},
		{"blue", "#3498DB"},
		{"green", "#2ECC71"},
		{"red", "#E74C3C"},
		{"purple", "#9B59B6"},
		{"orange", "#F39C12"},
		{"gray", "#6B7280"},
		{"CYAN", "#00FFFF"},
		{"Blue", "#3498DB"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, resolveOpenCodeColor(tt.input))
		})
	}
}

func TestResolveOpenCodeColor_ValidHex(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "#AABBCC", resolveOpenCodeColor("#aabbcc"))
	assert.Equal(t, "#123456", resolveOpenCodeColor("#123456"))
}

func TestResolveOpenCodeColor_UnknownFallsToGray(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "#6B7280", resolveOpenCodeColor("unknown-color"))
	assert.Equal(t, "#6B7280", resolveOpenCodeColor(""))
}
