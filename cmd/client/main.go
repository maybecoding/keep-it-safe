package main

import (
	"fmt"
	"net/http"

	"github.com/maybecoding/keep-it-safe/internal/client/tui"
	"github.com/maybecoding/keep-it-safe/pkg/logger"

	"github.com/maybecoding/keep-it-safe/internal/client/config"

	client "github.com/maybecoding/keep-it-safe/generated/client"
)

var (
	buildVersion = "N/A"
	buildTime    = "N/A"
)

func main() {
	// init logger
	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = logger.Init(cfg.Log.Level, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	// init client
	hc := http.Client{}
	c, err := client.NewClientWithResponses(cfg.Server.Address, client.WithHTTPClient(&hc))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to get Client for HTTP requests")
		return
	}

	err = tui.Run(c, buildVersion, buildTime)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to start TOI client")
		return
	}
}
