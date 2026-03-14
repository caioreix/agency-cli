package installer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/caioreix/agency-cli/internal/installer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDestinationDir(t *testing.T) {
	t.Parallel()
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	cwd, err := os.Getwd()
	require.NoError(t, err)

	tests := []struct {
		tool string
		want string
	}{
		{"claude-code", filepath.Join(cwd, ".claude", "agents")},
		{"copilot", filepath.Join(cwd, ".github", "agents")},
		{"cursor", filepath.Join(cwd, ".cursor", "rules")},
		{"windsurf", cwd},
		{"aider", cwd},
		{"opencode", filepath.Join(cwd, ".opencode", "agents")},
		{"openclaw", filepath.Join(home, ".openclaw", "agency-agents")},
		{"antigravity", filepath.Join(home, ".gemini", "antigravity", "skills")},
		{"gemini-cli", filepath.Join(home, ".gemini", "extensions", "agency-agents")},
		{"kimi-code", filepath.Join(cwd, ".kimi", "agents")},
		{"qwen", filepath.Join(cwd, ".qwen", "agents")},
	}

	for _, tt := range tests {
		t.Run(tt.tool, func(t *testing.T) {
			t.Parallel()
			got, err := installer.DestinationDir(tt.tool)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDestinationDir_Unknown(t *testing.T) {
	t.Parallel()
	got, err := installer.DestinationDir("unknown-tool")
	require.NoError(t, err)
	assert.Empty(t, got)
}
