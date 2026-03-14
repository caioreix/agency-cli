package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

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

// ─── List items ───────────────────────────────────────────────────────────────

type agentItem struct{ a *agent.Agent }

func (i agentItem) Title() string {
	if i.a.Emoji != "" {
		return i.a.Emoji + " " + i.a.Name
	}
	return i.a.Name
}
func (i agentItem) Description() string {
	vibe := i.a.Vibe
	if vibe == "" {
		vibe = i.a.Description
	}
	return "[" + i.a.Category + "] " + vibe
}
func (i agentItem) FilterValue() string {
	return i.a.Name + " " + i.a.Category + " " + i.a.Vibe + " " + i.a.Description
}

type toolItem struct {
	key  string
	name string
	desc string
}

func (i toolItem) Title() string       { return i.name }
func (i toolItem) Description() string { return i.desc }
func (i toolItem) FilterValue() string { return i.name + " " + i.key }

type scopeItem struct {
	label  string
	desc   string
	global bool
}

func (i scopeItem) Title() string       { return i.label }
func (i scopeItem) Description() string { return i.desc }
func (i scopeItem) FilterValue() string { return i.label }

// ─── Styles ───────────────────────────────────────────────────────────────────

var (
	accentColor = lipgloss.Color("205")
	dimColor    = lipgloss.Color("240")

	titleStyle      = lipgloss.NewStyle().Bold(true).Foreground(accentColor)
	breadcrumbStyle = lipgloss.NewStyle().Foreground(dimColor).Italic(true)
	successStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	errorStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	dimStyle        = lipgloss.NewStyle().Foreground(dimColor)
	fileStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	arrowStyle      = lipgloss.NewStyle().Foreground(accentColor)
)

// ─── Model ────────────────────────────────────────────────────────────────────

type Model struct {
	step    step
	spinner spinner.Model

	agentList list.Model
	toolList  list.Model
	scopeList list.Model

	selectedAgent  *agent.Agent
	selectedTool   string
	selectedGlobal bool

	result []string
	err    error
	width  int
	height int
}

func New() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(accentColor)
	return Model{
		step:    stepLoading,
		spinner: s,
		width:   80,
		height:  24,
	}
}

// ─── Tea interface ────────────────────────────────────────────────────────────

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, ensureRepo)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		h := listHeight(m.height)
		if m.step >= stepAgent {
			m.agentList.SetSize(m.width, h)
		}
		if m.step >= stepTool {
			m.toolList.SetSize(m.width, h)
		}
		if m.step >= stepScope {
			m.scopeList.SetSize(m.width, scopeListHeight(m.height))
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			// don't quit if a filter is being typed
			if m.step == stepAgent && m.agentList.FilterState() == list.Filtering {
				break
			}
			if m.step != stepLoading && m.step != stepInstalling {
				return m, tea.Quit
			}

		case "esc":
			switch m.step {
			case stepTool:
				// only go back if not currently typing a filter
				if m.toolList.FilterState() != list.Filtering {
					m.step = stepAgent
					return m, nil
				}
			case stepScope:
				m.step = stepTool
				return m, nil
			}

		case "enter":
			switch m.step {
			case stepAgent:
				if m.agentList.FilterState() == list.Filtering {
					break
				}
				if item, ok := m.agentList.SelectedItem().(agentItem); ok {
					m.selectedAgent = item.a
					m.step = stepTool
					m.toolList = makeToolList(m.width, m.height)
					return m, nil
				}

			case stepTool:
				if m.toolList.FilterState() == list.Filtering {
					break
				}
				if item, ok := m.toolList.SelectedItem().(toolItem); ok {
					m.selectedTool = item.key
					m.step = stepScope
					m.scopeList = makeScopeList(m.selectedTool, m.width, m.height)
					return m, nil
				}

			case stepScope:
				if item, ok := m.scopeList.SelectedItem().(scopeItem); ok {
					m.selectedGlobal = item.global
					m.step = stepInstalling
					return m, tea.Batch(m.spinner.Tick, doInstall(m.selectedAgent, m.selectedTool, m.selectedGlobal))
				}

			case stepDone, stepErr:
				return m, tea.Quit
			}
		}

	case repoReadyMsg:
		return m, loadAgents(string(msg))

	case agentsLoadedMsg:
		items := make([]list.Item, len(msg))
		for i, a := range msg {
			items[i] = agentItem{a}
		}
		d := styledDelegate()
		l := list.New(items, d, m.width, listHeight(m.height))
		l.Title = "Select an Agent"
		l.Styles.Title = titleStyle
		l.SetFilteringEnabled(true)
		m.agentList = l
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

	// Delegate remaining key/mouse events to the active list
	var cmd tea.Cmd
	switch m.step {
	case stepAgent:
		m.agentList, cmd = m.agentList.Update(msg)
	case stepTool:
		m.toolList, cmd = m.toolList.Update(msg)
	case stepScope:
		m.scopeList, cmd = m.scopeList.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	switch m.step {
	case stepLoading:
		return "\n  " + m.spinner.View() + " Loading agents…\n"

	case stepAgent:
		return m.agentList.View()

	case stepTool:
		bc := "  " + breadcrumbStyle.Render("Agent: "+m.selectedAgent.Name) + "\n"
		return bc + m.toolList.View()

	case stepScope:
		name := m.selectedAgent.Name
		if m.selectedAgent.Emoji != "" {
			name = m.selectedAgent.Emoji + " " + name
		}
		toolDisplay := m.selectedTool
		if c, ok := converter.All()[m.selectedTool]; ok {
			toolDisplay = c.Name()
		}
		bc := "  " + breadcrumbStyle.Render(name+" → "+toolDisplay) + "\n"
		return bc + m.scopeList.View()

	case stepInstalling:
		return "\n  " + m.spinner.View() + " Installing…\n"

	case stepDone:
		var sb strings.Builder
		sb.WriteString("\n")
		sb.WriteString("  " + successStyle.Render("✓ Installed successfully") + "\n\n")
		for _, f := range m.result {
			sb.WriteString("  " + arrowStyle.Render("→") + " " + fileStyle.Render(f) + "\n")
		}
		sb.WriteString("\n  " + dimStyle.Render("press any key to exit") + "\n")
		return sb.String()

	case stepErr:
		return "\n  " + errorStyle.Render("✗ "+m.err.Error()) + "\n\n" +
			"  " + dimStyle.Render("press any key to exit") + "\n"
	}
	return ""
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func listHeight(h int) int {
	v := h - 4
	if v < 8 {
		return 8
	}
	return v
}

func scopeListHeight(h int) int {
	v := listHeight(h)
	if v > 10 {
		return 10
	}
	return v
}

func styledDelegate() list.DefaultDelegate {
	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.
		Foreground(accentColor).BorderLeftForeground(accentColor)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.
		Foreground(accentColor).BorderLeftForeground(accentColor)
	return d
}

func makeToolList(width, height int) list.Model {
	convs := converter.All()
	items := make([]list.Item, 0, len(converter.SupportedTools))
	for _, key := range converter.SupportedTools {
		if c, ok := convs[key]; ok {
			items = append(items, toolItem{key: key, name: c.Name(), desc: c.Description()})
		}
	}
	d := styledDelegate()
	l := list.New(items, d, width, listHeight(height))
	l.Title = "Select a Tool"
	l.Styles.Title = titleStyle
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	return l
}

func makeScopeList(toolName string, width, height int) list.Model {
	conv, _ := converter.Get(toolName)

	var items []list.Item
	switch {
	case !conv.IsProjectScoped():
		items = []list.Item{
			scopeItem{label: "Global", desc: conv.Description(), global: true},
		}
	case toolName == "copilot":
		items = []list.Item{
			scopeItem{label: "Local", desc: ".github/agents/ (current repository)", global: false},
			scopeItem{label: "Global", desc: "~/.copilot/agents/", global: true},
		}
	default:
		items = []list.Item{
			scopeItem{label: "Local", desc: conv.Description(), global: false},
		}
	}

	d := styledDelegate()
	l := list.New(items, d, width, scopeListHeight(height))
	l.Title = "Select Scope"
	l.Styles.Title = titleStyle
	l.SetFilteringEnabled(false)
	l.SetShowStatusBar(false)
	return l
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
