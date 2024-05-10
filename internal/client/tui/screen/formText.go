package screen

import (
	"log"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/maybecoding/keep-it-safe/internal/client/tui/state"
)

type FormText struct {
	state       *state.State
	name        string
	modelBack   *tea.Model
	modelSubmit *tea.Model
	callback    func([]string) tea.Cmd

	input      textinput.Model
	textarea   textarea.Model
	focusIndex int
}

func NewFormText(state *state.State, name, placeholder string, modelBack *tea.Model, modelSubmit *tea.Model, callback func([]string) tea.Cmd) *FormText {
	input := textinput.New()
	input.Placeholder = "Name"
	input.Focus()
	input.CharLimit = 156
	input.Width = 20

	ta := textarea.New()
	ta.CharLimit = 10000
	ta.Placeholder = placeholder

	return &FormText{
		input:       input,
		state:       state,
		textarea:    ta,
		modelBack:   modelBack,
		modelSubmit: modelSubmit,
		name:        name,
		callback:    callback,
	}
}

var _ tea.Model = (*Form)(nil)

func (m *FormText) Init() tea.Cmd {
	setTableSize := func() tea.Msg { return tea.WindowSizeMsg{Width: m.state.F.WidthFull(), Height: m.state.F.HeightFull()} }
	m.input.SetValue("")
	m.textarea.SetValue("")
	m.focusIndex = 0
	return tea.Batch(textarea.Blink, setTableSize)
}

func (m *FormText) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft:
			if m.FocusSubmit() {
				return *m.modelBack, nil
			}
		case tea.KeyTab:
			if m.FocusInput() {
				m.input.Blur()
				cmds = append(cmds, m.textarea.Focus())
			} else if m.FocusText() {
				m.textarea.Blur()
			} else {
				cmds = append(cmds, m.input.Focus())
			}
			m.FocusNext()

		case tea.KeyEnter:
			log.Println("text enter")
			if m.FocusSubmit() && m.input.Value() != "" && m.textarea.Value() != "" {
				log.Println("text ok")
				if m.callback == nil {
					log.Println("text callback nil")
					return *m.modelSubmit, m.Fields
				}
				log.Println("text before callback")
				return *m.modelSubmit, m.callback(m.FieldsStr())
			}

		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.state.F.WinSize(msg)
		// height will set on view method according to other components
		m.textarea.SetWidth(m.state.F.Width())
	}

	if m.FocusText() {
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *FormText) View() string {
	top := `╭────────────────────────────────────╮
` + m.state.F.SingleHeader(m.name) + `
╰────────────────────────────────────╯
`

	top += "\n\n" + m.input.View() + "\n\n"

	submit := "Submit"
	if m.FocusSubmit() {
		submit = focusedStyle.Copy().Render("[" + submit + "]")
	}

	bottom := "\n\n" + submit + "\n\nctrl+c quit"
	if m.FocusSubmit() {
		bottom += " • ← back"
	}

	m.textarea.SetHeight(m.state.F.FreeSpace(top, bottom))

	return m.state.F.Render(top+m.textarea.View(), bottom)
}

func (m *FormText) FocusInput() bool {
	return m.focusIndex == 0
}

func (m *FormText) FocusText() bool {
	return m.focusIndex == 1
}

func (m *FormText) FocusSubmit() bool {
	return m.focusIndex == 2
}

func (m *FormText) FocusNext() {
	m.focusIndex = (m.focusIndex + 1) % 3
}

func (m *FormText) Fields() tea.Msg {
	return FormFields{m.input.Value(), m.textarea.Value()}
}

func (m *FormText) FieldsStr() []string {
	return []string{m.input.Value(), m.textarea.Value()}
}
