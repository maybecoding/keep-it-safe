package state

import (
	tea "github.com/charmbracelet/bubbletea"
	client "github.com/maybecoding/keep-it-safe/internal/client/api/v1"
)

type State struct {
	C            *client.ClientWithResponses
	Token        string
	WindowHeight int

	Register tea.Model
	Login    tea.Model
	Secrets  tea.Model
	Welcome  tea.Model
}
