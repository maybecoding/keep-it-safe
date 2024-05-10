package screen

import (
	"github.com/maybecoding/keep-it-safe/internal/client/tui/state"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type secretChooseKeyMap struct {
	Back key.Binding
	Quit key.Binding
}

func (k secretChooseKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit}
}

func (k secretChooseKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Back}, {k.Quit}}
}

type SecretChoose struct {
	state *state.State
	keys  secretChooseKeyMap
	help  help.Model

	secretTypesLen   int
	focusIndex       int
	modelsSecretType []*tea.Model
	modelBack        *tea.Model
}

func NewSecretChoose(state *state.State) *SecretChoose {
	keyMap := secretChooseKeyMap{
		Back: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc", "q", "ctrl+c"),
			key.WithHelp("esc/q", "quit"),
		),
	}
	help := help.New()

	return &SecretChoose{state: state, keys: keyMap, help: help, secretTypesLen: 4}
}

var _ tea.Model = (*SecretChoose)(nil)

func (m *SecretChoose) Init() tea.Cmd {
	return textinput.Blink
}

func (m *SecretChoose) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SecretChooseInit:
		m.modelBack = msg.Back
		m.modelsSecretType = msg.SecretTypes
	case tea.KeyMsg:
		s := msg.String()
		switch {
		case key.Matches(msg, m.keys.Back):
			return *m.state.Welcome, nil

		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		// if back
		case s == "left":
			return *m.modelBack, nil
		// if choosed secret type
		case s == "enter":
			nxt := m.modelsSecretType[m.focusIndex]
			return *nxt, (*nxt).Init()

			// Set focus to next input
		case s == "tab" || s == "shift+tab" || s == "up" || s == "down":
			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex >= m.secretTypesLen {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = m.secretTypesLen - 1
			}
		}
	case tea.WindowSizeMsg:
		m.state.F.WinSize(msg)
	}

	return m, nil
}

func (m *SecretChoose) View() string {
	view := `
╭────────────────────────────────────╮
│    Choose Secret Type              │
╰────────────────────────────────────╯
`

	for i := 0; i < m.secretTypesLen; i += 1 {
		item := secretTypeName(int32(i))
		if i == m.focusIndex {
			item = focusedStyle.Copy().Render("[" + item + "]")
		}
		view += "\n" + item + "\n"
	}

	return m.state.F.Render(view, m.help.View(m.keys))
}

func (m *SecretChoose) ModelsSet(back *tea.Model, nxt []*tea.Model) {
	m.modelBack = back
	m.modelsSecretType = nxt
}

type SecretChooseInit struct {
	Back        *tea.Model
	SecretTypes []*tea.Model
}
