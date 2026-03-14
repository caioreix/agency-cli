package converter

import (
	"fmt"

	"github.com/caioreix/agency-cli/internal/agent"
)

type Converter interface {
	Convert(a *agent.Agent, destDir string) ([]string, error)
	Name() string
	Description() string
	IsProjectScoped() bool
}

var registry = map[string]Converter{}

func Register(name string, c Converter) {
	registry[name] = c
}

func Get(name string) (Converter, error) {
	c, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
	return c, nil
}

func All() map[string]Converter {
	return registry
}

var SupportedTools = []string{
	"claude-code",
	"copilot",
	"cursor",
	"windsurf",
	"aider",
	"opencode",
	"openclaw",
	"antigravity",
	"gemini-cli",
	"qwen",
}
