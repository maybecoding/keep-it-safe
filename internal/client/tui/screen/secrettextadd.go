package screen

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maybecoding/keep-it-safe/generated/models"
	"github.com/maybecoding/keep-it-safe/internal/client/tui/state"
)

type SecretAddText struct {
	state      *state.State
	textarea   textarea.Model
	name       textinput.Model
	focusIndex int
}

// NewSecretAddText creates screen for input long text.
func NewSecretAddText(st *state.State) *SecretAddText {
	name := textinput.New()
	name.Placeholder = "Name"
	name.Focus()
	name.CharLimit = 156
	name.Width = 20

	ti := textarea.New()
	ti.Placeholder = "Prepare your text here."

	return &SecretAddText{
		name:     name,
		state:    st,
		textarea: ti,
	}
}

// Init TUI model.
func (m *SecretAddText) Init() tea.Cmd {
	setTableSize := func() tea.Msg { return tea.WindowSizeMsg{Width: m.state.F.WidthFull(), Height: m.state.F.HeightFull()} }
	m.name.SetValue("")
	m.textarea.SetValue("")
	m.focusIndex = 0
	return tea.Batch(textarea.Blink, setTableSize)
}

// Update TUI model.
func (m *SecretAddText) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			switch {
			case m.FocusName():
				m.name.Blur()
				cmds = append(cmds, m.textarea.Focus())
			case m.FocusText():
				m.textarea.Blur()
			default:
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
				return m.state.Secrets, data
			}

		case tea.KeyCtrlC:
			return m, tea.Quit
		default:
		}
	case tea.WindowSizeMsg:
		const padding = 10
		m.state.F.WinSize(msg)
		m.textarea.SetHeight(m.state.F.Height() - padding)
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

// View for TUI model.
func (m *SecretAddText) View() string {
	view := ""
	submit := submitText
	if m.FocusSubmit() {
		submit = focusedStyle.Copy().Render("[" + submit + "]")
	}

	return m.state.F.Render(fmt.Sprintf("%s\n\n%s\n\n%s\n\n%s", view, m.name.View(), m.textarea.View(), submit), "(ctrl+c to quit)")
}

// FocusName is name focused.
func (m *SecretAddText) FocusName() bool {
	return m.focusIndex == 0
}

// FocusText is text focused.
func (m *SecretAddText) FocusText() bool {
	return m.focusIndex == 1
}

// FocusSubmit - is submit button focused.
func (m *SecretAddText) FocusSubmit() bool {
	const submitNumber = 2
	return m.focusIndex == submitNumber
}

// FocusNext - moves focus to next element.
func (m *SecretAddText) FocusNext() {
	const elemCnt = 3
	m.focusIndex = (m.focusIndex + 1) % elemCnt
}
