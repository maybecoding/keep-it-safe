package screen

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/maybecoding/keep-it-safe/internal/client/tui/state"
)

// FormView - struct for readonly form.
type FormView struct {
	state        *state.State
	d            *FormViewInit
	vpTitleStyle func(strs ...string) string
	vpInfoStyle  func(strs ...string) string
	viewport     viewport.Model
}

// NewFormView returns new form view.
func NewFormView(st *state.State) *FormView {
	// titleStyle
	b := lipgloss.RoundedBorder()
	b.Right = "├"
	titleStyle := lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)

	// infoStyle
	bb := lipgloss.RoundedBorder()
	bb.Left = "┤"
	infoStyle := titleStyle.Copy().BorderStyle(bb)

	m := &FormView{state: st, vpTitleStyle: titleStyle.Render, vpInfoStyle: infoStyle.Render}
	m.viewport = viewport.New(st.F.Width(), 0)
	m.viewport.HighPerformanceRendering = false

	return m
}

// FormViewInit struct with initialization for screen FormView.
type FormViewInit struct {
	Name       string
	ModelBack  tea.Model
	TextName   string
	Text       string
	Components []FormViewComponent
}

// FormViewComponent settings for input component.
type FormViewComponent struct {
	Name, Value string
}

// type FormViewFields []string.
var _ tea.Model = (*FormView)(nil)

// Init TUI model.
func (m *FormView) Init() tea.Cmd {
	return nil
}

// Update TUI model.
func (m *FormView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		s := msg.String()
		switch {
		case s == tea.KeyCtrlC.String() || s == "q" || s == tea.KeyEsc.String():
			return m, tea.Quit

		case s == tea.KeyLeft.String() && m.d != nil:
			return m.d.ModelBack, nil
		}

	case FormViewInit:
		m.d = &msg
		m.viewport.SetContent(msg.Text)

	case tea.WindowSizeMsg:
		m.state.F.WinSize(msg)
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View for TUI model.
func (m *FormView) View() string {
	if m.d == nil {
		return m.state.F.Render("Component not initialized", "ctrl+c quit")
	}

	top := `╭────────────────────────────────────╮
` + m.state.F.SingleHeader(m.d.Name) + `
╰────────────────────────────────────╯
`

	for _, comp := range m.d.Components {
		top += fmt.Sprintf("\n\n%s: %s", comp.Name, comp.Value)
	}
	top += "\n"

	if m.d.TextName != "" {
		m.viewport.Width = m.state.F.Width()

		top += m.viewportHeaderView()

		footer := m.viewportFooterView()
		free := m.state.F.FreeSpace(top+"\n\n"+footer, "")

		m.viewport.Height = free
		m.viewport.YPosition = lipgloss.Height(top)
		return m.state.F.Render(top+"\n"+m.viewport.View()+"\n"+footer, "ctrl+c/q quit • ← back")
	}
	return m.state.F.Render(top, "ctrl+c/q quit • ← back")
}

func (m *FormView) viewportHeaderView() string {
	title := m.vpTitleStyle(m.d.TextName)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m *FormView) viewportFooterView() string {
	const wFullPrc = 100
	info := m.vpInfoStyle(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*wFullPrc))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
