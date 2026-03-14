package installer

import (
	"os"
	"path/filepath"

	"github.com/caioreix/agency-cli/internal/agent"
	"github.com/caioreix/agency-cli/internal/converter"
)

func Install(a *agent.Agent, toolName string, global bool) ([]string, error) {
	conv, err := converter.Get(toolName)
	if err != nil {
		return nil, err
	}

	destDir, err := DestinationDir(toolName)
	if err != nil {
		return nil, err
	}

	scope := converter.ScopeLocal
	if global {
		scope = converter.ScopeGlobal
	}

	return conv.Convert(a, destDir, scope)
}

func DestinationDir(toolName string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	switch toolName {
	case "claude-code":
		// claude-code handles its own multi-dir logic in the converter
		return filepath.Join(cwd, ".claude", "agents"), nil
	case "copilot":
		// Copilot handles its own multi-dir logic in the converter
		return filepath.Join(cwd, ".github", "agents"), nil
	case "cursor":
		return filepath.Join(cwd, ".cursor", "rules"), nil
	case "windsurf":
		return cwd, nil
	case "aider":
		return cwd, nil
	case "opencode":
		// opencode handles its own multi-dir logic in the converter
		return filepath.Join(cwd, ".opencode", "agents"), nil
	case "openclaw":
		return filepath.Join(home, ".openclaw", "agency-agents"), nil
	case "antigravity":
		return filepath.Join(home, ".gemini", "antigravity", "skills"), nil
	case "gemini-cli":
		// gemini-cli handles its own multi-dir logic in the converter
		return filepath.Join(cwd, ".gemini", "extensions", "agency-agents"), nil
	case "kimi-code":
		// kimi-code handles its own multi-dir logic in the converter
		return filepath.Join(cwd, ".kimi", "agents"), nil
	case "qwen":
		// qwen handles its own multi-dir logic in the converter
		return filepath.Join(cwd, ".qwen", "agents"), nil
	default:
		return "", nil
	}
}
