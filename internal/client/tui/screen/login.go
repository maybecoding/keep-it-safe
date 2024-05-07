package screen

import (
	"context"
	"fmt"
	"strings"

	"github.com/maybecoding/keep-it-safe/internal/client/api/v1/models"
	"github.com/maybecoding/keep-it-safe/internal/client/tui/state"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type loginKeyMap struct {
	Back key.Binding
	Quit key.Binding
}

func (k loginKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit}
}

func (k loginKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Back}, {k.Quit}}
}

type Login struct {
	state *state.State
	keys  loginKeyMap
	help  help.Model

	inputs     []textinput.Model
	focusIndex int

	buttonFocused string
	buttonBlurred string
	errorMessage  string
}

func NewLogin(state *state.State) *Login {
	keyMap := loginKeyMap{
		Back: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc", "quit"),
		),
	}
	help := help.New()

	// prepare inputs
	inputs := make([]textinput.Model, 2)
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

	buttonFocused := focusedStyle.Copy().Render("[ Submit ]")
	buttonBlurred := fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	return &Login{state: state, keys: keyMap, help: help, inputs: inputs, buttonFocused: buttonFocused, buttonBlurred: buttonBlurred}
}

var _ tea.Model = (*Login)(nil)

func (m *Login) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Login) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		s := msg.String()
		switch {
		case key.Matches(msg, m.keys.Back):
			return m.state.Welcome, nil

		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

			// Set focus to next input
		case s == "tab" || s == "shift+tab" || s == "enter" || s == "up" || s == "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, Login.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, m.Login
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
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
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *Login) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m *Login) View() string {
	view := `
╭────────────────────────────────────╮
│    Login                           │
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

	helpView := m.help.View(m.keys)
	height := m.state.WindowHeight - strings.Count(view, "\n") - strings.Count(helpView, "\n")

	return view + strings.Repeat("\n", height) + helpView
}

func (m Login) Login() tea.Msg {
	resp, err := m.state.C.LoginWithResponse(context.Background(), models.Credential{Login: m.inputs[0].Value(), Password: m.inputs[1].Value()})
	if err != nil {
		return ActionResult{err.Error()}
	}

	ar := ActionResult{}
	switch resp.StatusCode() {
	case 200:
		if resp != nil && resp.HTTPResponse != nil {
			auth := strings.Split(resp.HTTPResponse.Header.Get("Set-Cookie"), "=")
			if len(auth) > 1 && auth[0] != "" {
				m.state.Token = auth[1]
			} else {
				ar.Result = "Failed to get authorization token"
			}
		}
	case 400:
		ar.Result = "Bad request"
	case 401:
		ar.Result = "Login or password are incorrect"
	case 500:
		ar.Result = "Internal server error"
	}
	return ar
}
