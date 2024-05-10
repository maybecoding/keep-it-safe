package screen

import (
	"github.com/maybecoding/keep-it-safe/internal/client/tui/state"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	noStyle      = lipgloss.NewStyle()
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
)

type ActionResult struct {
	Cmd     tea.Cmd
	Result  string
	Success bool
}

type welcomeKeyMap struct {
	Login    key.Binding
	Register key.Binding
	Help     key.Binding
	Quit     key.Binding
}

func (k welcomeKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k welcomeKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Login, k.Register},
		{k.Help, k.Quit},
	}
}

type Welcome struct {
	state *state.State
	keys  welcomeKeyMap
	help  help.Model
}

func NewWelcome(state *state.State) *Welcome {
	keyMap := welcomeKeyMap{
		Login: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "login"),
		),
		Register: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "register"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc", "q", "ctrl+c"),
			key.WithHelp("esc/q", "quit"),
		),
	}
	return &Welcome{state: state, keys: keyMap, help: help.New()}
}

var _ tea.Model = (*Welcome)(nil)

func (m *Welcome) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m *Welcome) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Register):
			return *m.state.Register, (*m.state.Register).Init()
		case key.Matches(msg, m.keys.Login):
			return *m.state.Login, (*m.state.Login).Init()
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		}
	// if registration or login successful
	case ActionResult:
		return *m.state.Secrets, nil
	case tea.WindowSizeMsg:
		m.state.F.WinSize(msg)
	}

	return m, nil
}

func (m *Welcome) View() string {
	welcomeT := `╭────────────────────────────────────╮
│    Welcome to Keep IT Safe!        │
│                                    │
│     Please Register or Login       │
│ to Start keeping your secrets Safe │
╰────────────────────────────────────╯
`

	hT := m.help.View(m.keys)

	return m.state.F.Render(welcomeT, hT)
}
