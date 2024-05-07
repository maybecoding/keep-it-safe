package main

import (
	"fmt"
	"net/http"

	"github.com/maybecoding/keep-it-safe/internal/client/tui"

	"github.com/maybecoding/keep-it-safe/internal/client/config"

	client "github.com/maybecoding/keep-it-safe/internal/client/api/v1"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	// init client
	hc := http.Client{}
	c, err := client.NewClientWithResponses(cfg.Server.Address, client.WithHTTPClient(&hc))
	if err != nil {
		fmt.Println(err)
		return
	}

	tui.Run(c, cfg.TUI.WindowHeight)
}
