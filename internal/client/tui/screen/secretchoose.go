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

func (k *secretChooseKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit}
}

func (k *secretChooseKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Back}, {k.Quit}}
}

type SecretChoose struct {
	state *state.State
	keys  *secretChooseKeyMap
	help  help.Model

	modelBack        tea.Model
	modelsSecretType []tea.Model
	secretTypesLen   int
	focusIndex       int
}

func NewSecretChoose(st *state.State) *SecretChoose {
	keyMap := secretChooseKeyMap{
		Back: key.NewBinding(
			key.WithKeys(tea.KeyLeft.String()),
			key.WithHelp(leftText, "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys(tea.KeyEsc.String(), "q", tea.KeyCtrlC.String()),
			key.WithHelp("esc/q", "quit"),
		),
	}
	hlp := help.New()
	const secretTypesCnt = 4
	return &SecretChoose{state: st, keys: &keyMap, help: hlp, secretTypesLen: secretTypesCnt}
}

var _ tea.Model = (*SecretChoose)(nil)

// Init TUI model.
func (m *SecretChoose) Init() tea.Cmd {
	return textinput.Blink
}

// Update TUI model.
func (m *SecretChoose) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case SecretChooseInit:
		m.modelBack = msg.Back
		m.modelsSecretType = msg.SecretTypes
	case tea.KeyMsg:
		s := msg.String()
		switch {
		case key.Matches(msg, m.keys.Back):
			return m.state.Welcome, nil

		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		// if back
		case s == tea.KeyLeft.String():
			return m.modelBack, nil
		// if choosed secret type
		case s == tea.KeyEnter.String():
			nxt := m.modelsSecretType[m.focusIndex]
			return nxt, nxt.Init()

			// Set focus to next input
		case s == tea.KeyTab.String() || s == tea.KeyShiftTab.String() || s == tea.KeyUp.String() || s == tea.KeyDown.String():
			// Cycle indexes
			if s == tea.KeyUp.String() || s == tea.KeyShiftTab.String() {
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

// View for TUI model.
func (m *SecretChoose) View() string {
	view := `
╭────────────────────────────────────╮
│    Choose Secret Type              │
╰────────────────────────────────────╯
`

	// for i := 0; i < m.secretTypesLen; i++ {
	for i := range m.secretTypesLen {
		item := secretTypeName(int32(i))
		if i == m.focusIndex {
			item = focusedStyle.Copy().Render("[" + item + "]")
		}
		view += "\n" + item + "\n"
	}

	return m.state.F.Render(view, m.help.View(m.keys))
}

type SecretChooseInit struct {
	Back        tea.Model
	SecretTypes []tea.Model
}
