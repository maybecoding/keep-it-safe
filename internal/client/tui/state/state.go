package state

import (
	tea "github.com/charmbracelet/bubbletea"
	client "github.com/maybecoding/keep-it-safe/internal/client/api/v1"
	frame "github.com/maybecoding/keep-it-safe/internal/client/tui/render"
)

type State struct {
	C     *client.ClientWithResponses
	F     *frame.Frame
	Token string

	Welcome      *tea.Model
	Register     *tea.Model
	Login        *tea.Model
	Secrets      *tea.Model
	SecretChoose *tea.Model
	SecretAdd    SecretTypes
	FormView     *tea.Model
}

type SecretTypes struct {
	Credential *tea.Model
	Text       *tea.Model
	Binary     *tea.Model
	BankCard   *tea.Model
}
