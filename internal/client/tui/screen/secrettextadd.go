package screen

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maybecoding/keep-it-safe/internal/client/api/v1/models"
	"github.com/maybecoding/keep-it-safe/internal/client/tui/state"
)

type SecretAddText struct {
	state      *state.State
	name       textinput.Model
	textarea   textarea.Model
	focusIndex int
}

func NewSecretAddText(state *state.State) *SecretAddText {
	name := textinput.New()
	name.Placeholder = "Name"
	name.Focus()
	name.CharLimit = 156
	name.Width = 20

	ti := textarea.New()
	ti.Placeholder = "Prepare your text here."

	return &SecretAddText{
		name:     name,
		state:    state,
		textarea: ti,
	}
}

func (m *SecretAddText) Init() tea.Cmd {
	setTableSize := func() tea.Msg { return tea.WindowSizeMsg{Width: m.state.F.WidthFull(), Height: m.state.F.HeightFull()} }
	m.name.SetValue("")
	m.textarea.SetValue("")
	m.focusIndex = 0
	return tea.Batch(textarea.Blink, setTableSize)
}

func (m *SecretAddText) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			if m.FocusName() {
				m.name.Blur()
				cmds = append(cmds, m.textarea.Focus())
			} else if m.FocusText() {
				m.textarea.Blur()
			} else {
				cmds = append(cmds, m.name.Focus())
			}
			m.FocusNext()

		case tea.KeyEnter:
			if m.FocusSubmit() && m.name.Value() != "" && m.textarea.Value() != "" {
				text := m.textarea.Value()
				data := DataCmd(models.Data{
					SecretName: m.name.Value(),
					SecretType: SecretTypeText,
					Text:       &text,
				})
				return *m.state.Secrets, data
			}

		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.state.F.WinSize(msg)
		m.textarea.SetHeight(m.state.F.Height() - 10)
		m.textarea.SetWidth(m.state.F.Width())
	}

	if m.FocusText() {
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.name, cmd = m.name.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *SecretAddText) View() string {
	view := ""
	submit := "Submit"
	if m.FocusSubmit() {
		submit = focusedStyle.Copy().Render("[" + submit + "]")
	}

	return m.state.F.Render(fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s", view, m.name.View(), m.textarea.View(), submit), "(ctrl+c to quit)")
}

func (m *SecretAddText) FocusName() bool {
	return m.focusIndex == 0
}

func (m *SecretAddText) FocusText() bool {
	return m.focusIndex == 1
}

func (m *SecretAddText) FocusSubmit() bool {
	return m.focusIndex == 2
}

func (m *SecretAddText) FocusNext() {
	m.focusIndex = (m.focusIndex + 1) % 3
}
