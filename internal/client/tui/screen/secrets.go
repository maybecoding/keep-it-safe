package screen

import (
	"context"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/maybecoding/keep-it-safe/generated/models"
	"github.com/maybecoding/keep-it-safe/internal/client/tui/state"
	"github.com/maybecoding/keep-it-safe/pkg/logger"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Secrets screen with secret list.
type Secrets struct {
	state      *state.State
	list       list.Model
	keys       *secretsKeyMap
	addInitCmd tea.Cmd
}

// additional buttons.
type secretsKeyMap struct {
	reload   key.Binding
	itemAdd  key.Binding
	itemView key.Binding
}

// NewSecrets creates new secrets.
func NewSecrets(st *state.State, addInitCmd tea.Cmd) *Secrets {
	s := &Secrets{state: st}
	items := []list.Item{}

	s.list = list.New(items, list.NewDefaultDelegate(), st.F.Width(), st.F.Height())
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
		itemView: key.NewBinding(
			key.WithKeys("v", "enter"),
			key.WithHelp("v/enter", "view"),
		),
	}
}

// type for list item.
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

// Init TUI model.
func (m *Secrets) Init() tea.Cmd {
	setTableSize := func() tea.Msg { return tea.WindowSizeMsg{Width: m.state.F.WidthFull(), Height: m.state.F.HeightFull()} }

	return tea.Batch(setTableSize, m.Reload)
}

// Update TUI model.
func (m *Secrets) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	var secret *models.Secret
	if i, ok := m.list.SelectedItem().(item); ok {
		secret = &i.Secret
	} else {
		secret = nil
	}

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
		case key.Matches(msg, m.keys.itemView) && secret != nil:
			cmds = append(cmds, m.ItemView(secret))
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
			logger.Debug().Msg("add to cmd")
			cmds = append(cmds, msg.Cmd)
		}
	case FormViewInit:
		return *m.state.FormView, func() tea.Msg { return msg }

	// if got data add secret
	case models.Data:
		cmds = append(cmds, m.Add(msg))
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View for TUI model.
func (m *Secrets) View() string {
	return m.state.F.Render(m.list.View(), "")
}

// Reload reloads secret list.
func (m *Secrets) Reload() tea.Msg {
	logger.Debug().Msg("reload")
	resp, err := m.state.C.SecretListWithResponse(context.Background(), &models.SecretListParams{Authorization: m.state.Token})
	if err != nil {
		logger.Error().Msgf("Reload error: %s", err.Error())
		return ActionResult{Result: err.Error()}
	}
	ar := ActionResult{}
	switch resp.StatusCode() {
	case http.StatusOK:
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
	case http.StatusBadRequest:
		ar.Result = "Incorrect request " + string(resp.Body)
	case http.StatusUnauthorized:
		ar.Result = "User unautorized"
	case http.StatusInternalServerError:
		ar.Result = "Internal server error"
	default:
		ar.Result = fmt.Sprintf("Unhandled result code %d\n", resp.StatusCode())
	}
	logger.Debug().Interface("result", ar.Result).Msg("Reload result")
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

// DataCmd func for converting data to tea cmd.
func DataCmd(d models.Data) tea.Cmd {
	return func() tea.Msg {
		return d
	}
}

// Add new item.
func (m *Secrets) Add(d models.Data) tea.Cmd {
	return func() tea.Msg {
		logger.Debug().Msg("add")
		resp, err := m.state.C.SecretSet(context.Background(), &models.SecretSetParams{Authorization: m.state.Token}, d)
		if err != nil {
			return ActionResult{Result: err.Error()}
		}
		defer func() {
			_ = resp.Body.Close()
		}()
		ar := ActionResult{}
		switch resp.StatusCode {
		case http.StatusOK:
			ar.Result = "Successefully added secret " + d.SecretName
			ar.Success = true
			ar.Cmd = m.Reload
		case http.StatusBadRequest:
			ar.Result = "Incorrect request"
		case http.StatusUnauthorized:
			ar.Result = "User unautorized"
		case http.StatusInternalServerError:
			ar.Result = "Internal server error"
		default:
			ar.Result = fmt.Sprintf("Unhandled result code %d\n", resp.StatusCode)
		}
		logger.Debug().Interface("result", ar.Result).Msg("Reload result")
		return ar
	}
}

// ItemView view item from secret list.
func (m *Secrets) ItemView(secret *models.Secret) tea.Cmd {
	logger.Debug().Msg("Start View")
	return func() tea.Msg {
		resp, err := m.state.C.SecretGetByIDWithResponse(context.Background(), secret.Id, &models.SecretGetByIDParams{Authorization: m.state.Token})
		if err != nil {
			return ActionResult{Result: err.Error()}
		}
		ar := ActionResult{}
		switch resp.StatusCode() {
		case http.StatusOK:
			if resp.JSON200 == nil {
				ar.Result = "Failed to load response"
			} else {
				d := *resp.JSON200
				ar.Result = "Successefully got secret"
				ar.Success = true
				ar.Cmd = func() tea.Msg {
					viewInit := FormViewInit{Name: fmt.Sprintf(`Secret :%q`, secret.Name), ModelBack: m.state.Secrets}
					switch secret.Type {
					case SecretTypeCredentials:
						viewInit.Components = []FormViewComponent{{Name: "Login", Value: d.Credentials.Login}, {Name: "Password", Value: d.Credentials.Password}}
						return viewInit
					case SecretTypeBankCard:
						viewInit.Components = []FormViewComponent{
							{Name: "Number", Value: d.BankCard.Number},
							{Name: "Holder", Value: d.BankCard.Holder},
							{Name: "Valid", Value: d.BankCard.Valid},
							{Name: "ValidationCode", Value: d.BankCard.ValidationCode},
						}
						return viewInit
					case SecretTypeText:
						viewInit.TextName = "Text"
						viewInit.Text = *d.Text
						return viewInit
					case SecretTypeBinary:
						viewInit.TextName = "Binary Hex"
						b := *d.Binary
						viewInit.Text = hex.EncodeToString(b)
						return viewInit
					}
					return nil
				}
			}
		case http.StatusUnauthorized:
			ar.Result = "User unautorized"
		case http.StatusNotFound:
			ar.Result = "Not found"
		case http.StatusInternalServerError:
			ar.Result = "Internal server error"
		default:
			ar.Result = fmt.Sprintf("Unhandled result code %d\n", resp.StatusCode())
		}
		logger.Debug().Interface("result", ar.Result).Msg("Reload result")
		return ar
	}
}
