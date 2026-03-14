package converter

import (
	"testing"

	"github.com/caioreix/agency-cli/internal/agent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newTestAgent returns a reusable Agent fixture for converter tests.
func newTestAgent() *agent.Agent {
	return &agent.Agent{
		Name:        "Test Agent",
		Description: "A test agent for unit tests",
		Color:       "cyan",
		Emoji:       "🤖",
		Vibe:        "test vibe",
		Tools:       "bash,python",
		Category:    "testing",
		Slug:        "test-agent",
		Body: "## Mission\n\nDo test stuff.\n\n" +
			"## Identity\n\nBe a test.\n\n" +
			"## Communication Style\n\nSpeak plainly.\n",
	}
}

// ── registry ──────────────────────────────────────────────────────────────────

func TestGet_AllSupportedTools(t *testing.T) {
	t.Parallel()
	for _, name := range SupportedTools {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			c, err := Get(name)
			require.NoError(t, err)
			assert.NotNil(t, c)
			assert.NotEmpty(t, c.Name())
			assert.NotEmpty(t, c.Description())
		})
	}
}

func TestGet_UnknownTool(t *testing.T) {
	t.Parallel()
	_, err := Get("does-not-exist")
	assert.Error(t, err)
}

func TestAll_ContainsAllSupportedTools(t *testing.T) {
	t.Parallel()
	all := All()
	for _, name := range SupportedTools {
		assert.Contains(t, all, name, "tool %q missing from registry", name)
	}
}

// ── converter metadata ────────────────────────────────────────────────────────

func TestConverterMetadata(t *testing.T) {
	t.Parallel()
	tests := []struct {
		key        string
		wantScoped bool
	}{
		{"claude-code", true},
		{"copilot", true},
		{"cursor", true},
		{"windsurf", true},
		{"aider", true},
		{"opencode", true},
		{"openclaw", false},
		{"antigravity", false},
		{"gemini-cli", true},
		{"kimi-code", true},
		{"qwen", true},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			t.Parallel()
			c, err := Get(tt.key)
			require.NoError(t, err)
			assert.Equal(t, tt.wantScoped, c.IsProjectScoped(),
				"%s: unexpected IsProjectScoped value", tt.key)
		})
	}
}
