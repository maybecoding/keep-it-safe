package screen

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/maybecoding/keep-it-safe/generated/models"
	"github.com/maybecoding/keep-it-safe/internal/client/tui/state"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Register struct for registration form.
type Register struct {
	state *state.State

	buttonFocused string
	buttonBlurred string
	errorMessage  string
	inputs        []textinput.Model
	focusIndex    int
}

// NewRegister returns registration form.
func NewRegister(st *state.State) *Register {
	const inputCount = 2
	inputs := make([]textinput.Model, inputCount)
	var t textinput.Model
	for i := range inputs {
		t = textinput.New()
		t.Cursor.Style = focusedStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Login"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}
		inputs[i] = t
	}

	buttonFocused := focusedStyle.Copy().Render("[ " + submitText + " ]")
	buttonBlurred := fmt.Sprintf("[ %s ]", blurredStyle.Render(submitText))
	return &Register{state: st, inputs: inputs, buttonFocused: buttonFocused, buttonBlurred: buttonBlurred}
}

var _ tea.Model = (*Register)(nil)

func (m *Register) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Register) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		s := msg.String()
		switch {
		case s == tea.KeyLeft.String():
			return m.state.Welcome, nil

		case s == tea.KeyEsc.String() || s == tea.KeyCtrlC.String():
			return m, tea.Quit

			// Set focus to next input
		case s == tea.KeyTab.String() || s == tea.KeyShiftTab.String() ||
			s == tea.KeyEnter.String() || s == tea.KeyUp.String() || s == tea.KeyDown.String():
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, register.
			if s == tea.KeyEnter.String() && m.focusIndex == len(m.inputs) {
				return m, m.Register
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
	case ActionResult:
		if msg.Result == "" {
			return m.state.Welcome, func() tea.Msg { return msg }
		}
		m.errorMessage = msg.Result
	case tea.WindowSizeMsg:
		m.state.F.WinSize(msg)
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *Register) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

// View for TUI model.
func (m *Register) View() string {
	view := `
╭────────────────────────────────────╮
│    Register                        │
╰────────────────────────────────────╯
`
	for i := range m.inputs {
		view += m.inputs[i].View()
		if i < len(m.inputs)-1 {
			view += "\n"
		}
	}

	if m.focusIndex == len(m.inputs) {
		view += "\n\n" + m.buttonFocused + "\n\n"
	} else {
		view += "\n\n" + m.buttonBlurred + "\n\n"
	}
	if m.errorMessage != "" {
		view += errorStyle.Copy().Render(m.errorMessage) + "\n"
	}

	return m.state.F.Render(view, "← back • esc quit")
}

func (m *Register) Register() tea.Msg {
	resp, err := m.state.C.RegisterWithResponse(context.Background(), models.Credential{Login: m.inputs[0].Value(), Password: m.inputs[1].Value()})
	if err != nil {
		return ActionResult{Result: err.Error()}
	}

	ar := ActionResult{}
	switch resp.StatusCode() {
	case http.StatusOK:
		if resp != nil && resp.HTTPResponse != nil {
			auth := strings.Split(resp.HTTPResponse.Header.Get("Set-Cookie"), "=")
			if len(auth) > 1 && auth[0] != "" {
				m.state.Token = auth[1]
			} else {
				ar.Result = "Failed to get authorization token"
			}
		}
	case http.StatusBadRequest:
		ar.Result = "Bad request"
	case http.StatusConflict:
		ar.Result = "User already exists"
	case http.StatusInternalServerError:
		ar.Result = "Internal server Error"
	}
	return ar
}
