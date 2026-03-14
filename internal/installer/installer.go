package installer

import (
	"os"
	"path/filepath"

	"github.com/caioreix/agency-cli/internal/agent"
	"github.com/caioreix/agency-cli/internal/converter"
)

func Install(a *agent.Agent, toolName string) ([]string, error) {
	conv, err := converter.Get(toolName)
	if err != nil {
		return nil, err
	}

	destDir, err := DestinationDir(toolName)
	if err != nil {
		return nil, err
	}

	return conv.Convert(a, destDir)
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
		return filepath.Join(home, ".claude", "agents"), nil
	case "copilot":
		// Copilot handles its own multi-dir logic in the converter
		return filepath.Join(home, ".github", "agents"), nil
	case "cursor":
		return filepath.Join(cwd, ".cursor", "rules"), nil
	case "windsurf":
		return cwd, nil
	case "aider":
		return cwd, nil
	case "opencode":
		return filepath.Join(cwd, ".opencode", "agents"), nil
	case "openclaw":
		return filepath.Join(home, ".openclaw", "agency-agents"), nil
	case "antigravity":
		return filepath.Join(home, ".gemini", "antigravity", "skills"), nil
	case "gemini-cli":
		return filepath.Join(home, ".gemini", "extensions", "agency-agents"), nil
	case "qwen":
		return filepath.Join(cwd, ".qwen", "agents"), nil
	default:
		return "", nil
	}
}
