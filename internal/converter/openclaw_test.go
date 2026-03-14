package converter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenClaw_Convert_CreatesThreeFiles(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	a := newTestAgent()
	c, _ := Get("openclaw")

	files, err := c.Convert(a, dir, ScopeLocal)
	require.NoError(t, err)
	require.Len(t, files, 3)

	agentDir := filepath.Join(dir, "test-agent")
	assert.Equal(t, filepath.Join(agentDir, "SOUL.md"), files[0])
	assert.Equal(t, filepath.Join(agentDir, "AGENTS.md"), files[1])
	assert.Equal(t, filepath.Join(agentDir, "IDENTITY.md"), files[2])

	soul, err := os.ReadFile(files[0])
	require.NoError(t, err)
	assert.Contains(t, string(soul), "## Identity")

	identity, err := os.ReadFile(files[2])
	require.NoError(t, err)
	assert.Contains(t, string(identity), "🤖")
	assert.Contains(t, string(identity), "Test Agent")
}

// ── splitOpenClawSections ─────────────────────────────────────────────────────

func TestSplitOpenClawSections_IdentityGoesToSoul(t *testing.T) {
	t.Parallel()
	body := "## Identity\n\nPersona stuff.\n\n## Mission\n\nDo things.\n"
	soul, agents := splitOpenClawSections(body)

	assert.Contains(t, soul, "## Identity")
	assert.Contains(t, agents, "## Mission")
	assert.NotContains(t, soul, "## Mission")
	assert.NotContains(t, agents, "## Identity")
}

func TestSplitOpenClawSections_CommunicationGoesToSoul(t *testing.T) {
	t.Parallel()
	body := "## Communication Style\n\nBe concise.\n\n## Deliverables\n\nProduce code.\n"
	soul, agents := splitOpenClawSections(body)

	assert.Contains(t, soul, "## Communication Style")
	assert.Contains(t, agents, "## Deliverables")
}

func TestSplitOpenClawSections_NoSoulSections(t *testing.T) {
	t.Parallel()
	body := "## Mission\n\nDo stuff.\n\n## Workflow\n\nStep by step.\n"
	soul, agents := splitOpenClawSections(body)

	assert.Empty(t, soul)
	assert.Contains(t, agents, "## Mission")
	assert.Contains(t, agents, "## Workflow")
}

func TestSplitOpenClawSections_EmptyBody(t *testing.T) {
	t.Parallel()
	soul, agents := splitOpenClawSections("")
	assert.Empty(t, strings.TrimSpace(soul))
	assert.Empty(t, strings.TrimSpace(agents))
}
