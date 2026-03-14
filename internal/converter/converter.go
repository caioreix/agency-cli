package converter

import (
	"fmt"

	"github.com/caioreix/agency-cli/internal/agent"
)

const (
	ScopeDefault = ""
	ScopeLocal   = "local"
	ScopeGlobal  = "global"
)

type Converter interface {
	Convert(a *agent.Agent, destDir string, scope string) ([]string, error)
	Name() string
	Description() string
	IsProjectScoped() bool
}

var registry = map[string]Converter{} //nolint:gochecknoglobals // package-level converter registry

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

var SupportedTools = []string{ //nolint:gochecknoglobals // exported package-level tool list
	"claude-code",
	"copilot",
	"cursor",
	"windsurf",
	"aider",
	"opencode",
	"openclaw",
	"antigravity",
	"gemini-cli",
	"kimi-code",
	"qwen",
}
