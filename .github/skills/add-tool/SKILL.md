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
6. **Document** — Update `README.md` Supported Tools table.
7. **Verify** — Run `go test ./...` and `make lint` and confirm everything passes.

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

Always use distinct variable names and split the `WriteFile` call onto its own line — the `//nolint` comment pushes inline `if` forms over the 120-char golines limit:

```go
if mkdirErr := os.MkdirAll(dir, 0o755); mkdirErr != nil { //nolint:gosec // G301: world-traversable
    return nil, mkdirErr
}

writeErr := os.WriteFile(file, []byte(content), 0o644) //nolint:gosec // G306: world-readable
if writeErr != nil {
    return nil, writeErr
}
```

If a function writes multiple files, use distinct variable names (`writeErr`, `manifestErr`, `writeConfigErr`, …) to avoid `govet` shadow warnings.

---

## Existing Converter Reference

### Dual-scope — `claude.go` / `copilot.go` / `kimi.go` / `opencode.go` / `gemini.go` / `qwen.go`

`Description()` shows both paths separated by ` + `:

```go
func (c *myTool) Description() string   { return ".mytool/agents/ + ~/.mytool/agents/" }
func (c *myTool) IsProjectScoped() bool { return true }
```

`Convert()` resolves directories itself — ignores `destDir`:

```go
func (c *myTool) Convert(a *agent.Agent, _ string, scope string) ([]string, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return nil, err
    }
    cwd, err := os.Getwd()
    if err != nil {
        return nil, err
    }

    var dir string
    switch scope {
    case ScopeGlobal:
        dir = filepath.Join(home, ".mytool", "agents")
    default:
        dir = filepath.Join(cwd, ".mytool", "agents")
    }

    if mkdirErr := os.MkdirAll(dir, 0o755); mkdirErr != nil { //nolint:gosec // G301: world-traversable
        return nil, mkdirErr
    }

    outFile := filepath.Join(dir, a.Slug+".md")
    content := "---\nname: " + a.Name + "\n---\n" + a.Body

    writeErr := os.WriteFile(outFile, []byte(content), 0o644) //nolint:gosec // G306: world-readable
    if writeErr != nil {
        return nil, writeErr
    }

    return []string{outFile}, nil
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

### Global-only — `openclaw.go` / `antigravity.go`

`IsProjectScoped()` returns `false`. Scope param is `_`. `Description()` shows only the global path.

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

### Multi-file per agent — `kimi.go`

Some tools produce multiple files per agent (e.g. a YAML config + a Markdown system prompt). Return all file paths from `Convert()`:

```go
systemFile := filepath.Join(dir, a.Slug+".md")
writeErr := os.WriteFile(systemFile, []byte(a.Body), 0o644) //nolint:gosec // G306: world-readable
if writeErr != nil {
    return nil, writeErr
}

configFile := filepath.Join(dir, a.Slug+".yaml")
writeConfigErr := os.WriteFile(configFile, []byte(yamlContent), 0o644) //nolint:gosec // G306: world-readable
if writeConfigErr != nil {
    return nil, writeConfigErr
}

return []string{configFile, systemFile}, nil
```

### Subdirectory + manifest — `gemini.go`

Create a per-agent subdirectory and only write a manifest if it doesn't exist yet.

```go
skillDir := filepath.Join(baseDir, "skills", a.Slug)
if mkdirErr := os.MkdirAll(skillDir, 0o755); mkdirErr != nil { //nolint:gosec // G301: world-traversable
    return nil, mkdirErr
}
manifestFile := filepath.Join(baseDir, "gemini-extension.json")
if _, statErr := os.Stat(manifestFile); os.IsNotExist(statErr) {
    // write manifest once — use a distinct variable name to avoid shadow
    manifestErr := os.WriteFile(manifestFile, []byte(manifest), 0o644) //nolint:gosec // G306: world-readable
    if manifestErr != nil {
        return nil, manifestErr
    }
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

**Global-only** (converter uses `destDir`):
```go
case "new-tool":
    return filepath.Join(home, ".newtool", "agents"), nil
```

**Project-only** (converter uses `destDir`):
```go
case "new-tool":
    return filepath.Join(cwd, ".newtool", "agents"), nil
```

**Dual-scope** (converter resolves paths itself — return the local/cwd path as default):
```go
case "new-tool":
    // new-tool handles its own multi-dir logic in the converter
    return filepath.Join(cwd, ".newtool", "agents"), nil
```

### 3. `cmd/root.go` — extend `--tool` flag description

```go
const toolDesc = "target tool (claude-code, copilot, cursor, ..., new-tool)"
```

### 4. `README.md` — add row to Supported Tools table

```markdown
| new-tool | `.newtool/agents/ + ~/.newtool/agents/` | project + user |
```

Use `project` for project-only, `user` for global-only, `project + user` for dual-scope.

---

## Test Pattern

All test files must start with the `//nolint:testpackage` directive (tests share the `newTestAgent()` helper defined in `converter_test.go`):

```go
//nolint:testpackage // shares newTestAgent helper and tests unexported functions
package converter
```

### Local scope test — uses `t.Chdir` (incompatible with `t.Parallel`)

```go
// TestMyTool_Convert_Local uses t.Chdir which is incompatible with t.Parallel().
func TestMyTool_Convert_Local(t *testing.T) {
    t.Chdir(t.TempDir())

    cwd, err := os.Getwd()
    require.NoError(t, err)

    a := newTestAgent()
    c, _ := Get("my-tool")

    files, err := c.Convert(a, "", ScopeLocal)
    require.NoError(t, err)
    require.Len(t, files, 1)

    assert.Equal(t, filepath.Join(cwd, ".mytool", "agents", "test-agent.md"), files[0])

    content, err := os.ReadFile(files[0])
    require.NoError(t, err)
    assert.Contains(t, string(content), "name: Test Agent")
    assert.Contains(t, string(content), "## Mission")
}
```

### Default scope test — same as local (uses `t.Chdir`)

```go
// TestMyTool_Convert_Default uses t.Chdir which is incompatible with t.Parallel().
func TestMyTool_Convert_Default(t *testing.T) {
    t.Chdir(t.TempDir())

    cwd, err := os.Getwd()
    require.NoError(t, err)

    c, _ := Get("my-tool")
    files, err := c.Convert(newTestAgent(), "", ScopeDefault)
    require.NoError(t, err)

    assert.Equal(t, filepath.Join(cwd, ".mytool", "agents", "test-agent.md"), files[0])
}
```

### Global scope test — uses `t.Parallel` + cleanup

```go
func TestMyTool_Convert_Global(t *testing.T) {
    t.Parallel()
    home, err := os.UserHomeDir()
    require.NoError(t, err)

    c, _ := Get("my-tool")
    files, err := c.Convert(newTestAgent(), "", ScopeGlobal)
    require.NoError(t, err)
    require.Len(t, files, 1)

    assert.Equal(t, filepath.Join(home, ".mytool", "agents", "test-agent.md"), files[0])

    t.Cleanup(func() { os.Remove(files[0]) })
}
```

### Project-only global error test

```go
func TestMyTool_Convert_GlobalErrors(t *testing.T) {
    t.Parallel()
    c, _ := Get("my-tool")
    _, err := c.Convert(newTestAgent(), t.TempDir(), ScopeGlobal)
    assert.Error(t, err)
}
```

### Variable naming in assertions — avoid `encoded-compare` lint

Do **not** use variable names containing `yaml`, `json`, or `xml` in `assert.Equal` calls — `testifylint` will flag them. Use neutral names:

```go
// ✅ good
wantConfigFile := filepath.Join(cwd, ".mytool", "agents", "test-agent.yaml")
assert.Equal(t, wantConfigFile, files[0])

// ❌ bad — triggers testifylint encoded-compare
expectedYAMLPath := filepath.Join(cwd, ".mytool", "agents", "test-agent.yaml")
assert.Equal(t, expectedYAMLPath, files[0])
```

Minimum test cases per converter:

| Case | Required |
|---|---|
| Local install — verify file path and key content | ✅ |
| Default scope — verify same path as local | ✅ (dual-scope only) |
| Global install — success or expected error | ✅ |
| Scope ignored (global-only tools) | ✅ |
| Append behaviour (single-file converters) | ✅ |
| Optional fields omitted (e.g. no `Tools`, no `Emoji`) | when applicable |

---

## Checklist

Before finishing, verify all of the following:

- [ ] `internal/converter/<tool>.go` — `init()`, struct, all 4 interface methods implemented
- [ ] `internal/converter/converter.go` — tool name added to `SupportedTools`
- [ ] `internal/installer/installer.go` — case added in `DestinationDir` (with `// handles own multi-dir logic` comment for dual-scope)
- [ ] `cmd/root.go` — `--tool` flag description updated
- [ ] `internal/converter/<tool>_test.go` — starts with `//nolint:testpackage`, all required test cases present
- [ ] `README.md` — row added to Supported Tools table
- [ ] `go test ./...` — all tests pass
- [ ] `make lint` — 0 issues
- [ ] Paths verified for macOS, Linux, and Windows
