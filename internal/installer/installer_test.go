package installer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/caioreix/agency-cli/internal/installer"
)

func TestDestinationDir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("os.UserHomeDir: %v", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %v", err)
	}

	tests := []struct {
		tool string
		want string
	}{
		{"claude-code", filepath.Join(home, ".claude", "agents")},
		{"copilot", filepath.Join(cwd, ".github", "agents")},
		{"cursor", filepath.Join(cwd, ".cursor", "rules")},
		{"windsurf", cwd},
		{"aider", cwd},
		{"opencode", filepath.Join(cwd, ".opencode", "agents")},
		{"openclaw", filepath.Join(home, ".openclaw", "agency-agents")},
		{"antigravity", filepath.Join(home, ".gemini", "antigravity", "skills")},
		{"gemini-cli", filepath.Join(home, ".gemini", "extensions", "agency-agents")},
		{"qwen", filepath.Join(cwd, ".qwen", "agents")},
	}

	for _, tt := range tests {
		t.Run(tt.tool, func(t *testing.T) {
			got, err := installer.DestinationDir(tt.tool)
			if err != nil {
				t.Fatalf("DestinationDir(%q) returned unexpected error: %v", tt.tool, err)
			}
			if got != tt.want {
				t.Errorf("DestinationDir(%q)\n  got:  %s\n  want: %s", tt.tool, got, tt.want)
			}
		})
	}
}

func TestDestinationDir_Unknown(t *testing.T) {
	got, err := installer.DestinationDir("unknown-tool")
	if err != nil {
		t.Fatalf("DestinationDir(unknown) returned unexpected error: %v", err)
	}
	if got != "" {
		t.Errorf("DestinationDir(unknown) = %q, want empty string", got)
	}
}
