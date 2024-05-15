package tui

import (
	"encoding/hex"
	"fmt"

	frame "github.com/maybecoding/keep-it-safe/internal/client/tui/render"
	"github.com/maybecoding/keep-it-safe/internal/client/tui/screen"
	"github.com/maybecoding/keep-it-safe/internal/client/tui/state"
	"github.com/maybecoding/keep-it-safe/pkg/logger"

	tea "github.com/charmbracelet/bubbletea"
	client "github.com/maybecoding/keep-it-safe/generated/client"
	"github.com/maybecoding/keep-it-safe/generated/models"
)

func Run(c *client.ClientWithResponses, buildVersion, buildTime string) error {
	s := &state.State{C: c, F: frame.New().MarginSet(1, 2)}

	var secretAddInitCmd tea.Cmd = func() tea.Msg {
		return screen.SecretChooseInit{
			Back:        s.Secrets,
			SecretTypes: []tea.Model{s.SecretAdd.Credential, s.SecretAdd.Text, s.SecretAdd.Binary, s.SecretAdd.BankCard},
		}
	}

	s.Welcome = screen.NewWelcome(s, buildVersion, buildTime)
	s.Register = screen.NewRegister(s)
	s.Login = screen.NewLogin(s)
	s.Secrets = screen.NewSecrets(s, secretAddInitCmd)
	s.SecretChoose = screen.NewSecretChoose(s)
	s.FormView = screen.NewFormView(s)

	s.SecretAdd.Credential = screen.NewForm(s,
		"Add credentials",
		s.Secrets,
		s.Secrets,
		[]screen.InputParam{{Placeholder: "Secret Name"}, {Placeholder: "Login"}, {Placeholder: "Password", Password: true}},
		func(s []string) tea.Cmd {
			return screen.DataCmd(models.Data{
				SecretName: s[0],
				SecretType: screen.SecretTypeCredentials,
				Credentials: &models.DataCredentials{
					Login:    s[1],
					Password: s[2],
				},
			})
		},
	)

	s.SecretAdd.Text = screen.NewSecretAddText(s)
	s.SecretAdd.Binary = screen.NewFormText(s,
		"Add Binary Data",
		"Use Hex symbols for writing binary data.",
		s.Secrets,
		s.Secrets,
		func(s []string) tea.Cmd {
			hexBytes, err := hex.DecodeString(s[1])
			logger.Debug().Bytes("tui - Run - Add Binary Data - Bytes", hexBytes)
			if err != nil {
				logger.Error().Err(err).Msg("tui - Run - Add Binary Data")
			}
			return screen.DataCmd(models.Data{
				SecretName: s[0],
				SecretType: screen.SecretTypeBinary,
				Binary:     &hexBytes,
			})
		},
	)

	s.SecretAdd.BankCard = screen.NewForm(s,
		"Add bank card",
		s.Secrets,
		s.Secrets,
		[]screen.InputParam{{Placeholder: "Secret Name"}, {Placeholder: "Holder"}, {Placeholder: "Number"}, {Placeholder: "Valid"}, {Placeholder: "Code"}},
		func(s []string) tea.Cmd {
			return screen.DataCmd(models.Data{
				SecretName: s[0],
				SecretType: screen.SecretTypeBankCard,
				BankCard: &models.DataBankCard{
					Holder:         s[1],
					Number:         s[2],
					Valid:          s[3],
					ValidationCode: s[4],
				},
			})
		},
	)

	p := tea.NewProgram(s.Welcome)
	// p := tea.NewProgram(screen.NewForm(s, ), tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		return fmt.Errorf("tui - Run: %w", err)
	}
	return nil
}
