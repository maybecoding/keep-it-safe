package screen

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/maybecoding/keep-it-safe/generated/models"
	"github.com/maybecoding/keep-it-safe/internal/client/tui/state"
	"github.com/maybecoding/keep-it-safe/pkg/logger"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type loginKeyMap struct {
	Back key.Binding
	Quit key.Binding
}

// ShortHelp returns short help for help component.
func (k loginKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Back, k.Quit}
}

// FullHelp returns short help for help component.
func (k loginKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Back}, {k.Quit}}
}

// Login struct for login form.
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

// NewLogin returns new login form.
func NewLogin(state *state.State) *Login {
	keyMap := loginKeyMap{
		Back: key.NewBinding(
			key.WithKeys(tea.KeyLeft.String()),
			key.WithHelp("←", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("esc", "ctrl+c"),
			key.WithHelp("esc", "quit"),
		),
	}
	hlp := help.New()

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
	return &Login{state: state, keys: keyMap, help: hlp, inputs: inputs, buttonFocused: buttonFocused, buttonBlurred: buttonBlurred}
}

var _ tea.Model = (*Login)(nil)

// Init TUI model.
func (m *Login) Init() tea.Cmd {
	return textinput.Blink
}

// Update TUI model.
func (m *Login) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		s := msg.String()
		switch {
		case key.Matches(msg, m.keys.Back):
			return *m.state.Welcome, nil

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
			return *m.state.Secrets, (*m.state.Secrets).Init() // nil // func() tea.Msg { return msg }
		}
		m.errorMessage = msg.Result
	case tea.WindowSizeMsg:
		m.state.F.WinSize(msg)
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

// View for TUI model.
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

	return m.state.F.Render(view, m.help.View(m.keys))
}

// Login login on server.
func (m *Login) Login() tea.Msg {
	resp, err := m.state.C.LoginWithResponse(context.Background(),
		models.Credential{
			Login:    m.inputs[0].Value(),
			Password: m.inputs[1].Value(),
		})
	if err != nil {
		return ActionResult{Result: err.Error()}
	}

	ar := ActionResult{}
	switch resp.StatusCode() {
	case http.StatusOK:
		if resp != nil && resp.HTTPResponse != nil {
			auth := strings.Split(resp.HTTPResponse.Header.Get("Set-Cookie"), "=")
			if len(auth) > 1 && auth[0] != "" {
				m.state.Token = strings.ReplaceAll(auth[1], "\"", "")
				logger.Debug().Str("token", m.state.Token).Msg("set token")
			} else {
				logger.Error().Msgf("Failed to get authorization token %s", m.state.Token)
				ar.Result = "Failed to get authorization token"
			}
		}
	case http.StatusBadRequest:
		ar.Result = "Bad request"
	case http.StatusUnauthorized:
		ar.Result = "Login or password are incorrect"
	case http.StatusInternalServerError:
		ar.Result = "Internal Server error"
	}
	return ar
}
