# agency-cli

CLI tool to browse and install AI agents from [The Agency](https://github.com/msitarzewski/agency-agents) into your preferred agentic tool.

## Preview

![Preview](https://github.com/caioreix/agency-cli/raw/main/assets/preview.gif)

## Interactive TUI

The easiest way to browse and install agents is through the interactive TUI. Just run:

```bash
agency-cli
```

Or with a pre-selected tool:

```bash
agency-cli --tool cursor
```

### Workflow

The TUI guides you through a simple multi-step flow:

1. **Browse agents** - navigate by category or search across all agents
2. **Select a tool** - choose your target agentic tool
3. **Choose scope** - local (project) or global (user), when applicable
4. **Install** - the agent is converted and placed in the right location

### Keyboard shortcuts

| Key | Action |
|-----|--------|
| `↑` / `↓` or `j` / `k` | Navigate list |
| `←` / `→` | Switch categories |
| Type any character | Filter / search |
| `Backspace` | Delete last filter character |
| `Esc` | Clear filter / go back |
| `Enter` | Confirm selection |
| `q` | Quit (when not filtering) |
| `Ctrl+C` | Quit at any time |

> **Tip:** Filtering is smart — results are sorted by relevance: exact name match → prefix → contains → description/vibe/category.

## Installation

### Linux / macOS

```sh
curl -fsSL https://raw.githubusercontent.com/caioreix/agency-cli/main/install.sh | sh
```

The script auto-detects your OS and architecture, downloads the right binary from the [latest release](https://github.com/caioreix/agency-cli/releases/latest), and installs it to `/usr/local/bin` (requires sudo) or `$HOME/.local/bin` as a fallback.

### Windows (PowerShell)

```powershell
$url = "https://github.com/caioreix/agency-cli/releases/latest/download/agency-cli-windows-amd64.exe"
$dest = "$env:LOCALAPPDATA\Programs\agency-cli.exe"
New-Item -ItemType Directory -Force -Path (Split-Path $dest) | Out-Null
Invoke-WebRequest -Uri $url -OutFile $dest
# Add to PATH for current user (run once)
[Environment]::SetEnvironmentVariable("PATH", "$env:PATH;$(Split-Path $dest)", "User")
```

Restart your terminal after running the above so the updated PATH takes effect.

### Go developers

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
| claude-code | `.claude/agents/ + ~/.claude/agents/` | project + user |
| copilot | `.github/agents/ + ~/.copilot/agents/` | project + user |
| cursor | `.cursor/rules/` | project |
| windsurf | `.windsurfrules` | project |
| aider | `CONVENTIONS.md` | project |
| opencode | `.opencode/agents/ + ~/.config/opencode/agents/` | project + user |
| openclaw | `~/.openclaw/agency-agents/` | user |
| antigravity | `~/.gemini/antigravity/skills/` | user |
| gemini-cli | `.gemini/extensions/ + ~/.gemini/extensions/` | project + user |
| kimi-code | `.kimi/agents/ + ~/.kimi/agents/` | project + user |
| qwen | `.qwen/agents/ + ~/.qwen/agents/` | project + user |

## How it works

1. On first use, `agency-cli` clones the [agency-agents](https://github.com/msitarzewski/agency-agents) repository to `~/.cache/agency-cli/agency-agents/`
2. Use `agency-cli sync` to update the local cache
3. When you `get` an agent, the CLI converts the markdown file to the target tool's format and installs it to the correct location
