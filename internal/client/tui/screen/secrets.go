package screen

import (
	"context"
	"fmt"
	"log"

	"github.com/maybecoding/keep-it-safe/internal/client/api/v1/models"
	"github.com/maybecoding/keep-it-safe/internal/client/tui/state"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Secrets struct {
	state       *state.State
	list        list.Model
	keys        *secretsKeyMap
	addInitCmd  tea.Cmd
	viewInitCmd tea.Cmd
}

// additional buttons
type secretsKeyMap struct {
	reload  key.Binding
	itemAdd key.Binding
}

func NewSecrets(state *state.State, addInitCmd tea.Cmd, viewInitCmd tea.Cmd) *Secrets {
	s := &Secrets{state: state}
	items := []list.Item{}

	s.list = list.New(items, list.NewDefaultDelegate(), state.F.Width(), state.F.Height())
	s.list.Title = "Secrets"

	// add keys
	listKeys := newListKeyMap()

	s.list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.reload,
			listKeys.itemAdd,
		}
	}
	s.keys = listKeys

	s.addInitCmd = addInitCmd
	s.viewInitCmd = viewInitCmd

	return s
}

func newListKeyMap() *secretsKeyMap {
	return &secretsKeyMap{
		reload: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "reload items"),
		),
		itemAdd: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
	}
}

// type for list item
type item struct {
	models.Secret
}

func (i item) Title() string       { return i.Secret.Name }
func (i item) Description() string { return secretTypeName(i.Secret.Type) }
func (i item) FilterValue() string { return i.Secret.Name }

const (
	SecretTypeCredentials int32 = iota
	SecretTypeText
	SecretTypeBinary
	SecretTypeBankCard
)

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
	setTableSize := func() tea.Msg { return tea.WindowSizeMsg{Width: m.state.F.WidthFull(), Height: m.state.F.HeightFull()} }

	return tea.Batch(setTableSize, m.Reload)
}

func (m *Secrets) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		switch {
		case key.Matches(msg, m.keys.reload):
			cmds = append(cmds, m.Reload)
		case key.Matches(msg, m.keys.itemAdd):
			return *m.state.SecretChoose, m.addInitCmd
		}

	case tea.WindowSizeMsg:
		m.state.F.WinSize(msg)
		m.list.SetSize(m.state.F.Width(), m.state.F.Height())

	case ActionResult:
		if msg.Success {
			cmds = append(cmds, m.successRender())
		} else {
			cmds = append(cmds, m.failRender(msg.Result))
		}
		if msg.Cmd != nil {
			log.Println("add to cmd")
			cmds = append(cmds, msg.Cmd)
		}

	// if got data add secret
	case models.Data:
		cmds = append(cmds, m.Add(msg))
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Secrets) View() string {
	return m.state.F.Render(m.list.View(), "")
}

func (m *Secrets) Reload() tea.Msg {
	log.Println("reload")
	resp, err := m.state.C.SecretListWithResponse(context.Background(), &models.SecretListParams{Authorization: m.state.Token})
	if err != nil {
		log.Println("Reload error", err.Error())
		return ActionResult{Result: err.Error()}
	}
	ar := ActionResult{}
	switch resp.StatusCode() {
	case 200:
		if resp != nil && resp.JSON200 != nil {
			items := make([]list.Item, 0, len(*resp.JSON200))
			for _, s := range *resp.JSON200 {
				items = append(items, item{s})
			}
			ar.Cmd = m.list.SetItems(items)
			ar.Result = fmt.Sprintf("Loaded %d items", len(items))
			ar.Success = true
		} else {
			ar.Result = "Failed to fetch list of secrets"
		}
	case 400:
		ar.Result = "Incorrect request " + string(resp.Body)
	case 401:
		ar.Result = "User unautorized"
	case 500:
		ar.Result = "Internal server error"
	default:
		ar.Result = fmt.Sprintf("Unhandled result code %d\n", resp.StatusCode())
	}
	log.Println("Reload result", ar.Result)
	return ar
}

func (m *Secrets) successRender(s ...string) tea.Cmd {
	msg := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).Render(s...)
	return m.list.NewStatusMessage(msg)
}

func (m *Secrets) failRender(s ...string) tea.Cmd {
	msg := lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#b50436", Dark: "#b50436"}).Render(s...)
	return m.list.NewStatusMessage(msg)
}

func DataCmd(d models.Data) tea.Cmd {
	return func() tea.Msg {
		return d
	}
}

func (m *Secrets) Add(d models.Data) tea.Cmd {
	return func() tea.Msg {
		log.Println("add")
		resp, err := m.state.C.SecretSet(context.Background(), &models.SecretSetParams{Authorization: m.state.Token}, d)
		if err != nil {
			return ActionResult{Result: err.Error()}
		}
		ar := ActionResult{}
		switch resp.StatusCode {
		case 200:
			ar.Result = "Successefully added secret " + d.SecretName
			ar.Success = true
			ar.Cmd = m.Reload
		case 400:
			ar.Result = "Incorrect request"
		case 401:
			ar.Result = "User unautorized"
		case 500:
			ar.Result = "Internal server error"
		default:
			ar.Result = fmt.Sprintf("Unhandled result code %d\n", resp.StatusCode)
		}
		log.Println("Reload result", ar.Result)
		return ar
	}
}
