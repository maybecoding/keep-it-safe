// Package screen containing of TUI screens.
package screen

import (
	"github.com/maybecoding/keep-it-safe/internal/client/tui/state"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Form screen for form fields.
type Form struct {
	name        string
	modelBack   tea.Model
	modelSubmit tea.Model
	callback    func([]string) tea.Cmd
	state       *state.State

	inputs     []*textinput.Model
	focusIndex int
}

// NewForm creates form with fields.
func NewForm(st *state.State,
	name string,
	modelBack tea.Model,
	modelSubmit tea.Model,
	ips []InputParam,
	callback func([]string) tea.Cmd,
) *Form {
	// prepare inputs
	inputs := make([]*textinput.Model, 0, len(ips))
	for i, ip := range ips {
		ti := textinput.New()
		ti.Placeholder = ip.Placeholder
		if ip.Password {
			ti.EchoMode = textinput.EchoPassword
			ti.EchoCharacter = '•'
		}
		if i == 0 {
			ti.Focus()
			ti.PromptStyle = focusedStyle
			ti.TextStyle = focusedStyle
		}
		inputs = append(inputs, &ti)
	}

	return &Form{state: st, name: name, inputs: inputs, callback: callback, modelBack: modelBack, modelSubmit: modelSubmit}
}

type InputParam struct {
	Placeholder string
	Password    bool
}

type FormFields []string

var _ tea.Model = (*Form)(nil)

// Init TUI model.
func (m *Form) Init() tea.Cmd {
	return textinput.Blink
}

// Update TUI model.
func (m *Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		s := msg.String()
		switch {
		case s == tea.KeyLeft.String() && m.submitFocused() || s == tea.KeyCtrlLeft.String() && !m.submitFocused():
			return m.modelBack, nil

		case s == tea.KeyCtrlC.String() || s == tea.KeyEsc.String():
			return m, tea.Quit

			// Set focus to next input
		case s == tea.KeyTab.String() ||
			s == tea.KeyShiftTab.String() ||
			s == tea.KeyEnter.String() ||
			s == tea.KeyUp.String() ||
			s == tea.KeyDown.String():
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, Form.
			if s == tea.KeyEnter.String() && m.submitFocused() {
				if m.callback == nil {
					return m.modelSubmit, m.Fields
				}
				return m.modelSubmit, m.callback(m.FieldsStr())
			}

			// Cycle indexes
			if s == tea.KeyUp.String() || s == tea.KeyShiftTab.String() {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}
		}
	case tea.WindowSizeMsg:
		m.state.F.WinSize(msg)
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *Form) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		inp, cmd := m.inputs[i].Update(msg)
		m.inputs[i], cmds[i] = &inp, cmd
	}

	return tea.Batch(cmds...)
}

// View for TUI model.
func (m *Form) View() string {
	view := `╭────────────────────────────────────╮
` + m.state.F.SingleHeader(m.name) + `
╰────────────────────────────────────╯
`
	for i := range m.inputs {
		view += m.inputs[i].View() + "\n"
	}

	submit := submitText
	if m.submitFocused() {
		submit = focusedStyle.Render("[" + submit + "]")
	}
	view += "\n" + submit + "\n"

	backKey := leftText
	if !m.submitFocused() {
		backKey = "Ctrl + " + backKey
	}

	return m.state.F.Render(view, "ctrl+c quit • "+backKey+" back")
}

// Fields returns form fields as tea.msg.
func (m *Form) Fields() tea.Msg {
	ff := make(FormFields, 0, len(m.inputs))
	for _, input := range m.inputs {
		ff = append(ff, input.Value())
	}
	return ff
}

// FieldsStr returns form ffields as string slice.
func (m *Form) FieldsStr() []string {
	ff := make([]string, 0, len(m.inputs))
	for _, input := range m.inputs {
		ff = append(ff, input.Value())
	}
	return ff
}

func (m *Form) submitFocused() bool {
	return m.focusIndex == len(m.inputs)
}
