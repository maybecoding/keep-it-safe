package tui

import (
	"fmt"

	"github.com/maybecoding/keep-it-safe/internal/client/tui/screen"
	"github.com/maybecoding/keep-it-safe/internal/client/tui/state"

	tea "github.com/charmbracelet/bubbletea"
	client "github.com/maybecoding/keep-it-safe/internal/client/api/v1"
)

func Run(c *client.ClientWithResponses, height int) error {
	s := &state.State{C: c, WindowHeight: height}

	s.Register = screen.NewRegister(s)
	s.Login = screen.NewLogin(s)
	s.Secrets = screen.NewSecrets(s)
	s.Welcome = screen.NewWelcome(s)

	p := tea.NewProgram(s.Welcome, tea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		return fmt.Errorf("tui - Run: %w", err)
	}
	return nil
}
