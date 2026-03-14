//nolint:testpackage // tests unexported functions parseHex and nearestAnsi
package color

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ── parseHex ──────────────────────────────────────────────────────────────────

func TestParseHex(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		input   string
		r, g, b uint8
		wantOK  bool
	}{
		{"red with hash", "#FF0000", 255, 0, 0, true},
		{"green lowercase", "#00ff00", 0, 255, 0, true},
		{"blue mixed case", "#3498DB", 52, 152, 219, true},
		{"black", "#000000", 0, 0, 0, true},
		{"white", "#FFFFFF", 255, 255, 255, true},
		{"no hash prefix", "FF0000", 255, 0, 0, true},
		{"too short", "#FFF", 0, 0, 0, false},
		{"too long", "#FFFFFFF", 0, 0, 0, false},
		{"invalid hex chars", "#GGGGGG", 0, 0, 0, false},
		{"empty string", "", 0, 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r, g, b, ok := parseHex(tt.input)
			assert.Equal(t, tt.wantOK, ok)
			if tt.wantOK {
				assert.Equal(t, tt.r, r)
				assert.Equal(t, tt.g, g)
				assert.Equal(t, tt.b, b)
			}
		})
	}
}

// ── nearestAnsi ───────────────────────────────────────────────────────────────

func TestNearestAnsi_ReturnsNonEmpty(t *testing.T) {
	t.Parallel()
	assert.NotEmpty(t, nearestAnsi(0, 0, 0))
	assert.NotEmpty(t, nearestAnsi(255, 255, 255))
	assert.NotEmpty(t, nearestAnsi(255, 0, 0))
}

func TestNearestAnsi_RedLikeColor(t *testing.T) {
	t.Parallel()
	// Bright red (255, 80, 80) is in basicColors → should return the red code.
	assert.Equal(t, "\033[91m", nearestAnsi(255, 80, 80))
}

func TestNearestAnsi_GrayLikeColor(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "\033[90m", nearestAnsi(80, 80, 80))
}

// ── Apply / ApplyBold / ApplyDim (color disabled) ─────────────────────────────

func TestApply_NoColorEnv_ReturnsPlainText(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	assert.Equal(t, "hello", Apply("hello", "cyan"))
	assert.Equal(t, "hello", Apply("hello", "#FF0000"))
	assert.Equal(t, "hello", Apply("hello", "unknown-xyz"))
}

func TestApply_InvalidHex_NoColorEnv(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	assert.Equal(t, "text", Apply("text", "#ZZZZZZ"))
}

func TestApplyBold_NoColorEnv(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	assert.Equal(t, "bold text", ApplyBold("bold text"))
}

func TestApplyDim_NoColorEnv(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	assert.Equal(t, "dim text", ApplyDim("dim text"))
}

func TestApply_DumbTerm_ReturnsPlainText(t *testing.T) {
	t.Setenv("TERM", "dumb")
	assert.Equal(t, "hello", Apply("hello", "blue"))
}

// ── Supported ─────────────────────────────────────────────────────────────────

func TestSupported_ReturnsFalseWhenNoColor(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	assert.False(t, Supported())
}

func TestSupported_ReturnsFalseWhenDumbTerm(t *testing.T) {
	t.Setenv("TERM", "dumb")
	assert.False(t, Supported())
}
