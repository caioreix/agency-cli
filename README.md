# agency-cli

CLI tool to browse and install AI agents from [The Agency](https://github.com/msitarzewski/agency-agents) into your preferred agentic tool.

## Installation

```bash
go install github.com/caioreix/agency-cli@latest
```

## Usage

### List available agents

```bash
# List all agents
agency-cli list

# Filter by category
agency-cli list --category engineering
```

### Download and install an agent

```bash
# Install an agent for Cursor
agency-cli get code-reviewer --tool cursor

# Install an agent for Copilot
agency-cli get frontend-developer --tool copilot
```

### Sync the agent repository

```bash
# Update the local cache (clone on first run, pull on subsequent runs)
agency-cli sync
```

### List supported tools

```bash
agency-cli tools
```

## Supported Tools

| Tool | Destination | Scope |
|------|-------------|-------|
| claude-code | `~/.claude/agents/` | user |
| copilot | `~/.github/agents/` + `~/.copilot/agents/` | user |
| cursor | `.cursor/rules/` | project |
| windsurf | `.windsurfrules` | project |
| aider | `CONVENTIONS.md` | project |
| opencode | `.opencode/agents/` | project |
| openclaw | `~/.openclaw/agency-agents/` | user |
| antigravity | `~/.gemini/antigravity/skills/` | user |
| gemini-cli | `~/.gemini/extensions/agency-agents/` | user |
| qwen | `.qwen/agents/` | project |

## How it works

1. On first use, `agency-cli` clones the [agency-agents](https://github.com/msitarzewski/agency-agents) repository to `~/.cache/agency-cli/agency-agents/`
2. Use `agency-cli sync` to update the local cache
3. When you `get` an agent, the CLI converts the markdown file to the target tool's format and installs it to the correct location
