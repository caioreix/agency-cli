package converter

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenCode_Convert(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	a := newTestAgent()
	c, _ := Get("opencode")

	files, err := c.Convert(a, dir, ScopeLocal)
	require.NoError(t, err)
	require.Len(t, files, 1)
	assert.Equal(t, filepath.Join(dir, "test-agent.md"), files[0])

	content, err := os.ReadFile(files[0])
	require.NoError(t, err)
	assert.Contains(t, string(content), "name: Test Agent")
	assert.Contains(t, string(content), "mode: subagent")
	assert.Contains(t, string(content), "color:")
}

func TestOpenCode_Convert_GlobalErrors(t *testing.T) {
	t.Parallel()
	c, _ := Get("opencode")
	_, err := c.Convert(newTestAgent(), t.TempDir(), ScopeGlobal)
	assert.Error(t, err)
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
