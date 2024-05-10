package screen

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maybecoding/keep-it-safe/internal/client/tui/state"
)

type FormView struct {
	state *state.State
	d     *FormViewInit
}

func NewFormView(state *state.State) *FormView {
	return &FormView{state: state}
}

type FormViewInit struct {
	Name       string
	Components []FormViewComponent
	ModelBack  *tea.Model
}

type FormViewComponent struct {
	Name, Value string
}

type FormViewFields []string

var _ tea.Model = (*FormView)(nil)

func (m *FormView) Init() tea.Cmd {
	return nil
}

func (m *FormView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		s := msg.String()
		switch {
		case s == "ctrl+c" || s == "esc":
			return m, tea.Quit

		case s == "left" && m.d != nil:
			return *m.d.ModelBack, nil
		}

	case tea.WindowSizeMsg:
		m.state.F.WinSize(msg)
	}

	return m, nil
}

func (m *FormView) View() string {
	if m.d == nil {
		return m.state.F.Render("Component not initialized", "ctrl+c quit")
	}

	top := `╭────────────────────────────────────╮
` + m.state.F.SingleHeader(m.d.Name) + `
╰────────────────────────────────────╯`

	for _, comp := range m.d.Components {
		top += fmt.Sprintf("\n\n%s: %s\n\n", comp.Name, comp.Value)
	}

	return m.state.F.Render(top, "ctrl+c quit • ← back")
}
