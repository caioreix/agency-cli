package agent_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/caioreix/agency-cli/internal/agent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── Slugify ──────────────────────────────────────────────────────────────────

func TestSlugify(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"already lowercase", "hello", "hello"},
		{"uppercase", "HELLO", "hello"},
		{"spaces become dashes", "Hello World", "hello-world"},
		{"multiple spaces collapse", "hello   world", "hello-world"},
		{"leading and trailing spaces", "  hello  ", "hello"},
		{"special characters", "hello!@#world", "hello-world"},
		{"numbers preserved", "agent 42", "agent-42"},
		{"already slugified", "my-agent", "my-agent"},
		{"leading trailing dashes stripped", "-hello-", "hello"},
		{"mixed complex", "My Super Agent Name!", "my-super-agent-name"},
		{"only special chars", "!@#", ""},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, agent.Slugify(tt.input))
		})
	}
}

// ── Parse helpers ─────────────────────────────────────────────────────────────

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
	return path
}

// ── Parse ─────────────────────────────────────────────────────────────────────

func TestParse_ValidFullAgent(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := writeFile(t, dir, "agent.md", "---\n"+
		"name: My Agent\n"+
		"description: Does stuff\n"+
		"color: cyan\n"+
		"emoji: 🤖\n"+
		"vibe: chill\n"+
		"tools: bash\n"+
		"---\n"+
		"## Section\n\nBody text here.\n")

	a, err := agent.Parse(path, "engineering")
	require.NoError(t, err)

	assert.Equal(t, "My Agent", a.Name)
	assert.Equal(t, "Does stuff", a.Description)
	assert.Equal(t, "cyan", a.Color)
	assert.Equal(t, "🤖", a.Emoji)
	assert.Equal(t, "chill", a.Vibe)
	assert.Equal(t, "bash", a.Tools)
	assert.Equal(t, "engineering", a.Category)
	assert.Equal(t, "my-agent", a.Slug)
	assert.Equal(t, path, a.FilePath)
	assert.Contains(t, a.Body, "## Section")
	assert.Contains(t, a.Body, "Body text here.")
}

func TestParse_MinimalAgent(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := writeFile(t, dir, "min.md", "---\nname: Minimal\n---\nbody\n")

	a, err := agent.Parse(path, "design")
	require.NoError(t, err)

	assert.Equal(t, "Minimal", a.Name)
	assert.Equal(t, "design", a.Category)
	assert.Equal(t, "minimal", a.Slug)
	assert.Empty(t, a.Description)
	assert.Empty(t, a.Color)
}

func TestParse_DoubleQuotedValues(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := writeFile(t, dir, "dq.md", "---\nname: \"Quoted Agent\"\ndescription: \"quoted desc\"\n---\nbody\n")

	a, err := agent.Parse(path, "testing")
	require.NoError(t, err)

	assert.Equal(t, "Quoted Agent", a.Name)
	assert.Equal(t, "quoted desc", a.Description)
}

func TestParse_SingleQuotedValues(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := writeFile(t, dir, "sq.md", "---\nname: 'Single Agent'\ndescription: 'single desc'\n---\nbody\n")

	a, err := agent.Parse(path, "testing")
	require.NoError(t, err)

	assert.Equal(t, "Single Agent", a.Name)
	assert.Equal(t, "single desc", a.Description)
}

func TestParse_MissingFrontmatter(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := writeFile(t, dir, "nofm.md", "just body text\n")

	_, err := agent.Parse(path, "x")
	assert.Error(t, err)
}

func TestParse_NoNameField(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := writeFile(t, dir, "noname.md", "---\ndescription: no name here\n---\nbody\n")

	_, err := agent.Parse(path, "x")
	assert.Error(t, err)
}

func TestParse_FileNotFound(t *testing.T) {
	t.Parallel()
	_, err := agent.Parse("/nonexistent/path/agent.md", "x")
	assert.Error(t, err)
}

func TestParse_EmptyBody(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	path := writeFile(t, dir, "empty.md", "---\nname: Empty Body\n---\n")

	a, err := agent.Parse(path, "x")
	require.NoError(t, err)
	assert.Empty(t, a.Body)
}
