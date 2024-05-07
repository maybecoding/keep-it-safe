package screen

import (
	"context"

	"github.com/maybecoding/keep-it-safe/internal/client/api/v1/models"
	"github.com/maybecoding/keep-it-safe/internal/client/tui/state"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Secrets struct {
	state    *state.State
	list     list.Model
	docStyle lipgloss.Style
}

func NewSecrets(state *state.State) *Secrets {
	s := &Secrets{state: state}
	s.docStyle = lipgloss.NewStyle().Margin(1, 2)
	items := []list.Item{
		item2{title: "Raspberry Pi’s", desc: "I have ’em all over my house"},
		item2{title: "Nutella", desc: "It's good on toast"},
		item2{title: "Bitter melon", desc: "It cools you down"},
	}
	s.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
	s.list.Title = "Secrets"

	return s
}

// type for list item
type item struct {
	models.Secret
}

func (i item) Title() string       { return i.Secret.Name }
func (i item) Description() string { return secretTypeName(i.Secret.Type) }
func (i item) FilterValue() string { return i.Secret.Name }

func secretTypeName(st int32) string {
	switch st {
	case 0:
		return "Credentials"
	case 1:
		return "Text"
	case 2:
		return "Binary"
	case 3:
		return "BankCard"
	default:
		return "Undefined"
	}
}

func (m *Secrets) Init() tea.Cmd {
	return m.Reload
}

func (m *Secrets) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
		// case ActionResult:
		// fmt.Println("result", msg.Result)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *Secrets) View() string {
	return docStyle.Render(m.list.View())
}

func (m *Secrets) Reload() tea.Msg {
	resp, err := m.state.C.SecretListWithResponse(context.Background(), &models.SecretListParams{Authorization: m.state.Token})
	if err != nil {
		return ActionResult{err.Error()}
	}
	ar := ActionResult{}
	switch resp.StatusCode() {
	case 200:
		if resp != nil && resp.JSON200 != nil {
			items := make([]list.Item, 0, len(*resp.JSON200))
			for _, s := range *resp.JSON200 {
				items = append(items, item{s})
			}
			m.list = list.New(items, list.NewDefaultDelegate(), 0, 0)
		} else {
			ar.Result = "Failed to fetch list of secrets"
		}
	case 400:
		ar.Result = string(resp.Body)
	case 401:
		ar.Result = "User unautorized"
	case 500:
		ar.Result = "Internal server error"
	}
	return ar
}

type item2 struct {
	title, desc string
}

func (i item2) Title() string       { return i.title }
func (i item2) Description() string { return i.desc }
func (i item2) FilterValue() string { return i.title }
