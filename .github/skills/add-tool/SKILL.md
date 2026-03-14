---
name: add-tool
description: Adds new tool/IDE converters to agency-cli. Use this skill when asked to add support for a new tool or IDE. Researches official docs, infers config paths for macOS/Linux/Windows, handles local/global scope, and generates production-ready Go code following project conventions.
---

## Mission

You are a senior Go developer specialised in the `agency-cli` codebase. When asked to add support for a new tool, you carry out the full implementation end-to-end: research, code, tests, registration, and verification.

**Steps to follow — in order:**

1. **Research** — Search the tool's official documentation. Find: how it loads AI agents / custom instructions / rules / skills; the expected file format; and where config files live on macOS, Linux, and Windows.
2. **Plan** — Decide scope (project vs global), paths per OS, and file format.
3. **Implement** — Create `internal/converter/<toolname>.go`.
4. **Test** — Create `internal/converter/<toolname>_test.go`.
5. **Register** — Update `internal/installer/installer.go`, `internal/converter/converter.go`, and `cmd/root.go`.
6. **Verify** — Run `go test ./...` and confirm everything passes.

---

## Project Layout

```
internal/
  converter/
    converter.go          ← registry, Converter interface, SupportedTools
    <toolname>.go         ← one file per tool
    <toolname>_test.go
  installer/
    installer.go          ← DestinationDir switch
cmd/
  root.go                 ← --tool flag description
```

---

## Converter Interface

```go
type Converter interface {
    Convert(a *agent.Agent, destDir string, scope string) ([]string, error)
    Name() string          // display name, e.g. "My Tool"
    Description() string   // install path shown in `agency-cli tools`
    IsProjectScoped() bool
}
```

Register via `init()`:

```go
func init() { //nolint:gochecknoinits // required by cobra/converter
    Register("toolname", &myTool{})
}
```

---

## Agent Fields

```go
type Agent struct {
    Name        string // e.g. "DevOps Engineer"
    Description string // one-line summary
    Color       string // e.g. "cyan", "blue"
    Emoji       string // e.g. "🤖"
    Vibe        string // short personality note
    Tools       string // comma-separated tool list
    Category    string // directory category
    Slug        string // kebab-case, e.g. "devops-engineer"
    Body        string // full markdown body after frontmatter
    FilePath    string // source path
}
```

---

## Scope Handling

| Scenario | `IsProjectScoped()` | Behaviour |
|---|---|---|
| Project-only | `true` | Return `errors.New("... is project-scoped; --scope global is not supported")` when `scope == ScopeGlobal` |
| Global-only | `false` | Ignore scope, always install to user home; use `_` for scope param |
| Dual-scope | `true` | Branch on `scope`: project dir for `ScopeLocal`/`ScopeDefault`, user home for `ScopeGlobal` |

Available constants: `ScopeLocal = "local"`, `ScopeGlobal = "global"`, `ScopeDefault = ""`.

---

## Cross-Platform Path Rules

- **Always use `os.UserHomeDir()`** for the user home directory — returns `~` on macOS/Linux and `%USERPROFILE%` on Windows.
- **Always use `filepath.Join()`** for path construction — handles OS path separators automatically.
- **`%APPDATA%` on Windows** — some tools store config in `%APPDATA%` instead of `~`. Use an env-var lookup with fallback:

```go
func appDataDir() (string, error) {
    if appdata := os.Getenv("APPDATA"); appdata != "" {
        return appdata, nil
    }
    return os.UserHomeDir()
}
```

- Prefer `os.UserHomeDir()` + tool-specific subfolder unless the tool's documentation explicitly calls out a Windows-specific path (e.g. `%APPDATA%\Tool`).

---

## File Permissions (lint-compliant)

```go
os.MkdirAll(dir, 0o755)                          //nolint:gosec // G301: world-traversable
os.WriteFile(file, []byte(content), 0o644)        //nolint:gosec // G306: world-readable
```

---

## Existing Converter Reference

### Simple global — `claude.go`

Writes a `.md` file with frontmatter. Ignores scope (always global).

```go
func (c *claudeCode) Convert(a *agent.Agent, destDir string, _ string) ([]string, error) {
    if err := os.MkdirAll(destDir, 0o755); err != nil { //nolint:gosec // G301: world-traversable
        return nil, err
    }
    outFile := filepath.Join(destDir, a.Slug+".md")
    content := "---\nname: " + a.Name + "\ndescription: " + a.Description + "\n---\n" + a.Body
    if err := os.WriteFile(outFile, []byte(content), 0o644); err != nil { //nolint:gosec // G306: world-readable
        return nil, err
    }
    return []string{outFile}, nil
}
```

### Dual-scope — `copilot.go`

Different directories for local vs global.

```go
func (c *copilot) Convert(a *agent.Agent, _ string, scope string) ([]string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }
    cwd, err := os.Getwd()
    if err != nil {
        return nil, err
    }
    var dirs []string
    switch scope {
    case ScopeGlobal:
        dirs = []string{filepath.Join(home, ".copilot", "agents")}
    default:
        dirs = []string{filepath.Join(cwd, ".github", "agents")}
    }
    var files []string
    for _, dir := range dirs {
        if mkdirErr := os.MkdirAll(dir, 0o755); mkdirErr != nil { //nolint:gosec // G301: world-traversable
            return nil, mkdirErr
        }
        outFile := filepath.Join(dir, a.Slug+".md")
        if writeErr := os.WriteFile(outFile, []byte(content), 0o644); writeErr != nil { //nolint:gosec // G306: world-readable
            return nil, writeErr
        }
        files = append(files, outFile)
    }
    return files, nil
}
```

### Project-scoped with global error — `cursor.go`

```go
func (c *cursor) Convert(a *agent.Agent, destDir string, scope string) ([]string, error) {
    if scope == ScopeGlobal {
        return nil, errors.New("cursor is project-scoped; --scope global is not supported")
    }
    // ...
}
```

### Append-to-single-file — `windsurf.go` / `aider.go`

Read existing file and append; write header only on first run.

```go
var content string
if existing, err := os.ReadFile(outFile); err == nil {
    content = string(existing) + entry
} else {
    content = header + entry
}
```

### Subdirectory + manifest — `gemini.go`

Create a per-agent subdirectory and only write a manifest if it doesn't exist yet.

```go
skillDir := filepath.Join(destDir, "skills", a.Slug)
if err := os.MkdirAll(skillDir, 0o755); err != nil { //nolint:gosec // G301: world-traversable
    return nil, err
}
manifestFile := filepath.Join(destDir, "gemini-extension.json")
if _, err := os.Stat(manifestFile); os.IsNotExist(err) {
    // write manifest once
}
```

### Multi-file split — `openclaw.go`

Split agent body into domain-specific files (SOUL.md, AGENTS.md, IDENTITY.md).

---

## Files to Update

### 1. `internal/converter/converter.go` — append to `SupportedTools`

```go
var SupportedTools = []string{
    // ...existing...
    "new-tool",
}
```

### 2. `internal/installer/installer.go` — add case to `DestinationDir`

```go
case "new-tool":
    return filepath.Join(home, ".newtool", "agents"), nil
// or for project-scoped:
case "new-tool":
    return filepath.Join(cwd, ".newtool", "agents"), nil
```

### 3. `cmd/root.go` — extend `--tool` flag description

```go
const toolDesc = "target tool (claude-code, copilot, cursor, ..., new-tool)"
```

---

## Test Pattern

```go
package converter

import (
    "os"
    "path/filepath"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestMyTool_Convert(t *testing.T) {
    t.Parallel()
    dir := t.TempDir()
    a := newTestAgent()
    c, _ := Get("my-tool")

    files, err := c.Convert(a, dir, ScopeLocal)
    require.NoError(t, err)
    require.Len(t, files, 1)
    assert.Equal(t, filepath.Join(dir, "test-agent.md"), files[0])

    content, err := os.ReadFile(files[0])
    require.NoError(t, err)
    assert.Contains(t, string(content), "name: Test Agent")
    assert.Contains(t, string(content), "## Mission")
}

func TestMyTool_Convert_GlobalErrors(t *testing.T) {
    t.Parallel()
    c, _ := Get("my-tool")
    _, err := c.Convert(newTestAgent(), t.TempDir(), ScopeGlobal)
    assert.Error(t, err)
}
```

Minimum test cases per converter:

| Case | Required |
|---|---|
| Local install — verify file path and key content | ✅ |
| Global install — success or expected error | ✅ |
| Scope ignored (global-only tools) | ✅ |
| Append behaviour (single-file converters) | ✅ |
| Optional fields omitted (e.g. no `Tools`, no `Emoji`) | when applicable |

---

## Checklist

Before finishing, verify all of the following:

- [ ] `internal/converter/<tool>.go` — `init()`, struct, all 4 interface methods implemented
- [ ] `internal/converter/converter.go` — tool name added to `SupportedTools`
- [ ] `internal/installer/installer.go` — case added in `DestinationDir`
- [ ] `cmd/root.go` — `--tool` flag description updated
- [ ] `internal/converter/<tool>_test.go` — all required test cases present
- [ ] `go test ./...` — all tests pass
- [ ] Paths verified for macOS, Linux, and Windows
