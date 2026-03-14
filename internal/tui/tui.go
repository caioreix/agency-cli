package tui

import (
"fmt"
"sort"
"strings"
"unicode"

"github.com/charmbracelet/bubbles/spinner"
tea "github.com/charmbracelet/bubbletea"
"github.com/charmbracelet/lipgloss"
rw "github.com/mattn/go-runewidth"

"github.com/caioreix/agency-cli/internal/agent"
"github.com/caioreix/agency-cli/internal/converter"
"github.com/caioreix/agency-cli/internal/installer"
"github.com/caioreix/agency-cli/internal/repo"
)

// ─── Steps ────────────────────────────────────────────────────────────────────

type step int

const (
stepLoading step = iota
stepAgent
stepTool
stepScope
stepInstalling
stepDone
stepErr
)

// ─── Messages ─────────────────────────────────────────────────────────────────

type repoReadyMsg string
type agentsLoadedMsg []*agent.Agent
type installDoneMsg []string
type errMsg struct{ err error }

// ─── Row types ────────────────────────────────────────────────────────────────

type rowKind int

const (
rowAgent    rowKind = iota
rowCategory         // non-selectable section header
)

type row struct {
kind  rowKind
agent *agent.Agent
cat   string
score int // lower = better match (for sorting)
}

// ─── Styles ───────────────────────────────────────────────────────────────────

var (
accentColor = lipgloss.Color("205")
dimColor    = lipgloss.Color("240")

selStyle     = lipgloss.NewStyle().Foreground(accentColor).Bold(true)
dimStyle     = lipgloss.NewStyle().Foreground(dimColor)
successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
errStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
fileStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
arrowStyle   = lipgloss.NewStyle().Foreground(accentColor)
filterStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
catStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true)
helpStyle    = lipgloss.NewStyle().Foreground(dimColor)
breadStyle   = lipgloss.NewStyle().Foreground(dimColor)
)

// ─── Scope option ─────────────────────────────────────────────────────────────

type scopeOpt struct {
label  string
desc   string
global bool
}

// ─── Category order ───────────────────────────────────────────────────────────

var categoryOrder = []string{
"design", "engineering", "game-development", "marketing",
"paid-media", "sales", "product", "project-management",
"testing", "support", "spatial-computing", "specialized",
}

// ─── Model ────────────────────────────────────────────────────────────────────

type Model struct {
step    step
width   int
height  int
spinner spinner.Model

// Agent step
allAgents  []*agent.Agent
categories []string
catIdx     int
filter     string
cursor     int // index in visibleRows()
offset     int // scroll offset in visibleRows()

// Tool / scope cursors
toolCursor  int
	toolFilter  string
scopeCursor int

// Selections
selectedAgent  *agent.Agent
selectedTool   string
selectedGlobal bool

// Result
result []string
err    error
}

func New() Model {
s := spinner.New()
s.Spinner = spinner.Dot
s.Style = lipgloss.NewStyle().Foreground(accentColor)
return Model{step: stepLoading, spinner: s, width: 80, height: 24}
}

// ─── Tea interface ────────────────────────────────────────────────────────────

func (m Model) Init() tea.Cmd {
return tea.Batch(m.spinner.Tick, ensureRepo)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
switch msg := msg.(type) {

case tea.WindowSizeMsg:
m.width, m.height = msg.Width, msg.Height
return m, nil

case tea.KeyMsg:
return m.handleKey(msg)

case repoReadyMsg:
return m, loadAgents(string(msg))

case agentsLoadedMsg:
m.allAgents = []*agent.Agent(msg)
m.categories = extractCategories(m.allAgents)
m.step = stepAgent
return m, nil

case installDoneMsg:
m.result = []string(msg)
m.step = stepDone
return m, nil

case errMsg:
m.err = msg.err
m.step = stepErr
return m, nil

case spinner.TickMsg:
if m.step == stepLoading || m.step == stepInstalling {
var cmd tea.Cmd
m.spinner, cmd = m.spinner.Update(msg)
return m, cmd
}
}

return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
if msg.String() == "ctrl+c" {
return m, tea.Quit
}
switch m.step {
case stepAgent:
return m.handleAgentKey(msg)
case stepTool:
return m.handleToolKey(msg)
case stepScope:
return m.handleScopeKey(msg)
case stepDone, stepErr:
return m, tea.Quit
}
return m, nil
}

func (m Model) handleAgentKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
rows := m.visibleRows()

switch msg.String() {
case "esc":
if m.filter != "" {
m.filter = ""
m.cursor, m.offset = m.firstAgentRow(m.visibleRows()), 0
}
return m, nil

case "backspace":
if len(m.filter) > 0 {
runes := []rune(m.filter)
m.filter = string(runes[:len(runes)-1])
m.cursor, m.offset = m.firstAgentRow(m.visibleRows()), 0
}
return m, nil

case "left":
if m.filter == "" && len(m.categories) > 0 {
m.catIdx = (m.catIdx - 1 + len(m.categories)) % len(m.categories)
m.cursor, m.offset = 0, 0
}
return m, nil

case "right":
if m.filter == "" && len(m.categories) > 0 {
m.catIdx = (m.catIdx + 1) % len(m.categories)
m.cursor, m.offset = 0, 0
}
return m, nil

case "up":
m.moveCursor(rows, -1)
return m, nil

case "down":
m.moveCursor(rows, 1)
return m, nil

case "enter":
if m.cursor < len(rows) && rows[m.cursor].kind == rowAgent {
m.selectedAgent = rows[m.cursor].agent
m.step = stepTool
m.toolCursor = 0
}
return m, nil
}

// j/k/q only when filter is empty
if m.filter == "" {
switch msg.String() {
case "k":
m.moveCursor(rows, -1)
return m, nil
case "j":
m.moveCursor(rows, 1)
return m, nil
case "q":
return m, tea.Quit
}
}

// Printable chars → filter (searches all categories)
if msg.Type == tea.KeyRunes {
for _, r := range msg.Runes {
if unicode.IsPrint(r) {
m.filter += string(r)
newRows := m.visibleRows()
m.cursor = m.firstAgentRow(newRows)
m.offset = 0
}
}
}

return m, nil
}

func (m *Model) moveCursor(rows []row, dir int) {
pos := m.cursor + dir
for pos >= 0 && pos < len(rows) {
if rows[pos].kind == rowAgent {
m.cursor = pos
vh := m.listHeight()
if m.cursor < m.offset {
m.offset = m.cursor
} else if m.cursor >= m.offset+vh {
m.offset = m.cursor - vh + 1
}
return
}
pos += dir
}
}

func (m Model) handleToolKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	tools := m.filteredTools()

	switch msg.String() {
	case "esc":
		if m.toolFilter != "" {
			m.toolFilter = ""
			m.toolCursor = 0
		} else {
			m.step = stepAgent
		}
		return m, nil

	case "q":
		if m.toolFilter != "" {
			m.toolFilter = ""
			m.toolCursor = 0
		} else {
			m.step = stepAgent
		}
		return m, nil

	case "backspace":
		if len(m.toolFilter) > 0 {
			runes := []rune(m.toolFilter)
			m.toolFilter = string(runes[:len(runes)-1])
			m.toolCursor = 0
		}
		return m, nil

	case "up", "k":
		if m.toolCursor > 0 {
			m.toolCursor--
		}
		return m, nil

	case "down", "j":
		if m.toolCursor < len(tools)-1 {
			m.toolCursor++
		}
		return m, nil

	case "enter":
		if m.toolCursor < len(tools) {
			m.selectedTool = tools[m.toolCursor].key
			m.step = stepScope
			m.scopeCursor = 0
			m.toolFilter = ""
		}
		return m, nil
	}

	if msg.Type == tea.KeyRunes {
		for _, r := range msg.Runes {
			if unicode.IsPrint(r) {
				m.toolFilter += string(r)
				m.toolCursor = 0
			}
		}
	}

	return m, nil
}

func (m Model) handleScopeKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
scopes := m.scopeOptions()
switch msg.String() {
case "q", "esc":
m.step = stepTool
case "up", "k":
if m.scopeCursor > 0 {
m.scopeCursor--
}
case "down", "j":
if m.scopeCursor < len(scopes)-1 {
m.scopeCursor++
}
case "enter":
if len(scopes) > 0 {
m.selectedGlobal = scopes[m.scopeCursor].global
m.step = stepInstalling
return m, tea.Batch(m.spinner.Tick, doInstall(m.selectedAgent, m.selectedTool, m.selectedGlobal))
}
}
return m, nil
}

// ─── View ─────────────────────────────────────────────────────────────────────

func (m Model) View() string {
switch m.step {
case stepLoading:
return "\n  " + m.spinner.View() + " Loading agents…\n"
case stepAgent:
return m.viewAgents()
case stepTool:
return m.viewTools()
case stepScope:
return m.viewScope()
case stepInstalling:
return "\n  " + m.spinner.View() + " Installing…\n"
case stepDone:
return m.viewDone()
case stepErr:
return "\n  " + errStyle.Render("✗ "+m.err.Error()) + "\n\n" +
"  " + dimStyle.Render("press any key to exit") + "\n"
}
return ""
}

func (m Model) viewAgents() string {
w := maxInt(m.width, 40)
div := dimStyle.Render(strings.Repeat("─", w))
rows := m.visibleRows()

// ── line 1: category nav ──
var catLine string
if m.filter != "" {
// count matching agents
nAgents := 0
for _, r := range rows {
if r.kind == rowAgent {
nAgents++
}
}
catLine = " " + catStyle.Render(fmt.Sprintf("%d agents", nAgents)) +
dimStyle.Render(" across all categories")
} else {
cat := ""
if len(m.categories) > 0 {
cat = m.categories[m.catIdx]
}
catLine = " " + dimStyle.Render("◀ ") +
catStyle.Render(cat) +
dimStyle.Render(fmt.Sprintf("  %d/%d  ▶", m.catIdx+1, len(m.categories)))
}

// ── line 2: filter input ──
var filterLine string
if m.filter != "" {
filterLine = " " + dimStyle.Render("/ ") + filterStyle.Render(m.filter) + filterStyle.Render("█")
} else {
filterLine = " " + dimStyle.Render("/ ") + dimStyle.Render("type to search…")
}

// ── list ──
nameW := 28
if w > 110 {
nameW = 38
} else if w > 90 {
nameW = 32
}
vibeW := w - nameW - 5
if vibeW < 10 {
vibeW = 10
}

vh := m.listHeight()
end := m.offset + vh
if end > len(rows) {
end = len(rows)
}

var sb strings.Builder
sb.WriteString(catLine + "\n")
sb.WriteString(filterLine + "\n")
sb.WriteString(div + "\n")

for i := m.offset; i < end; i++ {
r := rows[i]
if r.kind == rowCategory {
// Category separator header
label := "  " + r.cat + " "
line := dimStyle.Render("  ── ") + catStyle.Render(r.cat) + dimStyle.Render(" "+strings.Repeat("─", maxInt(0, w-len(r.cat)-8)))
_ = label
sb.WriteString(line + "\n")
continue
}
a := r.agent
name := a.Name
if a.Emoji != "" {
name = a.Emoji + " " + a.Name
}
vibe := a.Vibe
if vibe == "" {
vibe = a.Description
}
namePad := padW(truncW(name, nameW), nameW)
vibeStr := truncW(vibe, vibeW)

if i == m.cursor {
sb.WriteString(arrowStyle.Render("❯ ") + selStyle.Render(namePad) + "  " + dimStyle.Render(vibeStr) + "\n")
} else {
sb.WriteString("  " + namePad + "  " + dimStyle.Render(vibeStr) + "\n")
}
}

// Fill blank rows so divider stays at consistent position
for i := end - m.offset; i < vh; i++ {
sb.WriteString("\n")
}

// ── footer ──
scroll := ""
if len(rows) > vh {
scroll = dimStyle.Render(fmt.Sprintf("%d-%d/%d  ", m.offset+1, end, len(rows)))
}
sb.WriteString(div + "\n")
if m.filter != "" {
sb.WriteString(" " + scroll + helpStyle.Render("↑↓ navigate  enter select  esc clear filter  ctrl+c quit") + "\n")
} else {
sb.WriteString(" " + scroll + helpStyle.Render("↑↓ j k  ←→ category  type to search  enter select  ctrl+c quit") + "\n")
}
return sb.String()
}

func (m Model) viewTools() string {
	w := maxInt(m.width, 40)
	div := dimStyle.Render(strings.Repeat("─", w))

	agentName := m.selectedAgent.Name
	if m.selectedAgent.Emoji != "" {
		agentName = m.selectedAgent.Emoji + " " + m.selectedAgent.Name
	}

	tools := m.filteredTools()

	var filterLine string
	if m.toolFilter != "" {
		filterLine = " " + dimStyle.Render("/ ") + filterStyle.Render(m.toolFilter) + filterStyle.Render("█")
	} else {
		filterLine = " " + dimStyle.Render("/ ") + dimStyle.Render("type to search…")
	}

	var sb strings.Builder
	sb.WriteString(" " + breadStyle.Render("agent: ") + selStyle.Render(agentName) + "\n")
	sb.WriteString(filterLine + "\n")
	sb.WriteString(div + "\n")

	nameW := 14
	for i, t := range tools {
		namePad := padW(truncW(t.name, nameW), nameW)
		desc := truncW(t.desc, w-nameW-8)
		if i == m.toolCursor {
			sb.WriteString(arrowStyle.Render("❯ ") + selStyle.Render(namePad) + "  " + dimStyle.Render(desc) + "\n")
		} else {
			sb.WriteString("  " + namePad + "  " + dimStyle.Render(desc) + "\n")
		}
	}

	if len(tools) == 0 {
		sb.WriteString("  " + dimStyle.Render("no tools found") + "\n")
	}

	sb.WriteString(div + "\n")
	sb.WriteString(" " + helpStyle.Render("↑↓ j k  type to filter  enter select  esc back  ctrl+c quit") + "\n")
	return sb.String()
}

func (m Model) viewScope() string {
w := maxInt(m.width, 40)
div := dimStyle.Render(strings.Repeat("─", w))

agentName := m.selectedAgent.Name
if m.selectedAgent.Emoji != "" {
agentName = m.selectedAgent.Emoji + " " + m.selectedAgent.Name
}
toolName := m.selectedTool
if c, ok := converter.All()[m.selectedTool]; ok {
toolName = c.Name()
}

var sb strings.Builder
sb.WriteString(" " + breadStyle.Render(agentName+" → ") + selStyle.Render(toolName) + "\n")
sb.WriteString(div + "\n")

labelW := 8
for i, s := range m.scopeOptions() {
labelPad := padW(s.label, labelW)
desc := truncW(s.desc, w-labelW-8)
if i == m.scopeCursor {
sb.WriteString(arrowStyle.Render("❯ ") + selStyle.Render(labelPad) + "  " + dimStyle.Render(desc) + "\n")
} else {
sb.WriteString("  " + labelPad + "  " + dimStyle.Render(desc) + "\n")
}
}

sb.WriteString(div + "\n")
sb.WriteString(" " + helpStyle.Render("↑↓ j k  enter confirm  esc back  ctrl+c quit") + "\n")
return sb.String()
}

func (m Model) viewDone() string {
var sb strings.Builder
sb.WriteString("\n  " + successStyle.Render("✓ Installed successfully") + "\n\n")
for _, f := range m.result {
sb.WriteString("  " + arrowStyle.Render("→") + " " + fileStyle.Render(f) + "\n")
}
sb.WriteString("\n  " + dimStyle.Render("press any key to exit") + "\n")
return sb.String()
}

// ─── Data helpers ─────────────────────────────────────────────────────────────

// visibleRows builds the display list: category headers + agent rows.
// When filter is active: all categories, grouped, sorted by match score.
// When filter is empty: only the current category.
func (m Model) visibleRows() []row {
if m.filter == "" {
cat := ""
if len(m.categories) > 0 && m.catIdx < len(m.categories) {
cat = m.categories[m.catIdx]
}
var rows []row
rows = append(rows, row{kind: rowCategory, cat: cat})
for _, a := range m.allAgents {
if a.Category == cat {
rows = append(rows, row{kind: rowAgent, agent: a})
}
}
return rows
}

// Filter mode: search all, group by category, sort by score
f := strings.ToLower(m.filter)
type group struct {
agents []row
}
groups := map[string]*group{}

for _, a := range m.allAgents {
sc := matchScore(a, f)
if sc < 0 {
continue
}
if groups[a.Category] == nil {
groups[a.Category] = &group{}
}
groups[a.Category].agents = append(groups[a.Category].agents, row{kind: rowAgent, agent: a, score: sc})
}

// Sort agents within each group by score then name
for _, g := range groups {
sort.Slice(g.agents, func(i, j int) bool {
if g.agents[i].score != g.agents[j].score {
return g.agents[i].score < g.agents[j].score
}
return g.agents[i].agent.Name < g.agents[j].agent.Name
})
}

// Build rows in canonical category order
var rows []row
for _, cat := range categoryOrder {
g, ok := groups[cat]
if !ok {
continue
}
rows = append(rows, row{kind: rowCategory, cat: cat})
rows = append(rows, g.agents...)
}
return rows
}

// matchScore returns how well agent a matches the filter f (lower = better).
// Returns -1 if no match.
func matchScore(a *agent.Agent, f string) int {
name := strings.ToLower(a.Name)
if name == f {
return 0
}
if strings.HasPrefix(name, f) {
return 1
}
if strings.Contains(name, f) {
return 2
}
haystack := strings.ToLower(a.Description + " " + a.Vibe + " " + a.Category)
if strings.Contains(haystack, f) {
return 3
}
return -1
}

// toolEntry holds display info for a tool.
type toolEntry struct {
	key  string
	name string
	desc string
}

// filteredTools returns tools matching the current toolFilter.
func (m Model) filteredTools() []toolEntry {
	convs := converter.All()
	f := strings.ToLower(m.toolFilter)
	var out []toolEntry
	for _, key := range converter.SupportedTools {
		c, ok := convs[key]
		if !ok {
			continue
		}
		if f != "" {
			hay := strings.ToLower(c.Name() + " " + c.Description() + " " + key)
			if !strings.Contains(hay, f) {
				continue
			}
		}
		out = append(out, toolEntry{key: key, name: c.Name(), desc: c.Description()})
	}
	return out
}

// firstAgentRow returns the index of the first agent row (skips category headers).
func (m Model) firstAgentRow(rows []row) int {
for i, r := range rows {
if r.kind == rowAgent {
return i
}
}
return 0
}

func (m Model) listHeight() int {
// catLine(1) + filterLine(1) + divider(1) + divider(1) + footer(1) = 5
h := m.height - 5
if h < 5 {
return 5
}
return h
}

func (m Model) scopeOptions() []scopeOpt {
conv, err := converter.Get(m.selectedTool)
if err != nil {
return nil
}
switch {
case !conv.IsProjectScoped():
return []scopeOpt{{label: "Global", desc: conv.Description(), global: true}}
case m.selectedTool == "copilot":
return []scopeOpt{
{label: "Local", desc: ".github/agents/ (current repository)", global: false},
{label: "Global", desc: "~/.copilot/agents/", global: true},
}
default:
return []scopeOpt{{label: "Local", desc: conv.Description(), global: false}}
}
}

func extractCategories(agents []*agent.Agent) []string {
seen := map[string]bool{}
for _, a := range agents {
seen[a.Category] = true
}
var cats []string
for _, c := range categoryOrder {
if seen[c] {
cats = append(cats, c)
}
}
return cats
}

// ─── String helpers ───────────────────────────────────────────────────────────

func truncW(s string, maxW int) string {
total := 0
for i, r := range s {
cw := rw.RuneWidth(r)
if total+cw > maxW {
if i > 0 {
return s[:i] + "…"
}
return ""
}
total += cw
}
return s
}

func padW(s string, n int) string {
w := 0
for _, r := range s {
w += rw.RuneWidth(r)
}
if w >= n {
return s
}
return s + strings.Repeat(" ", n-w)
}

func maxInt(a, b int) int {
if a > b {
return a
}
return b
}

// ─── Commands ─────────────────────────────────────────────────────────────────

func ensureRepo() tea.Msg {
dir, err := repo.EnsureRepo()
if err != nil {
return errMsg{err}
}
return repoReadyMsg(dir)
}

func loadAgents(repoDir string) tea.Cmd {
return func() tea.Msg {
agents, err := agent.ListAll(repoDir)
if err != nil {
return errMsg{err}
}
return agentsLoadedMsg(agents)
}
}

func doInstall(a *agent.Agent, tool string, global bool) tea.Cmd {
return func() tea.Msg {
files, err := installer.Install(a, tool, global)
if err != nil {
return errMsg{err}
}
return installDoneMsg(files)
}
}
