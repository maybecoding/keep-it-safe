// Package app - DI package.
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/maybecoding/keep-it-safe/internal/server/adapters/encrypter"
	"github.com/maybecoding/keep-it-safe/internal/server/app/get"
	"github.com/maybecoding/keep-it-safe/internal/server/config"
	"github.com/maybecoding/keep-it-safe/pkg/logger"
	"github.com/maybecoding/keep-it-safe/pkg/postgres"
	"github.com/maybecoding/keep-it-safe/pkg/starter"
)

// App struct with main app structural structs.
type App struct {
	cfg  *config.Config
	pg   *postgres.Postgres
	encr *encrypter.Encrypter

	starter *starter.Starter
	get     *get.Get
}

// New returns new struct.
func New(cfg *config.Config) *App {
	a := &App{cfg: cfg}
	return a
}

// Init - initialize of components witch needs initialization with error return.
func (a *App) Init() error {
	// init logger
	err := logger.Init(a.cfg.Log.Level, false)
	if err != nil {
		return fmt.Errorf("app - Init - logger.Init: %w", err)
	}
	// init pg
	pg, err := postgres.New(a.cfg.DB.Path)
	if err != nil {
		logger.Error().Err(err).Msg("app - Init - postgres.New")
		return fmt.Errorf("app - Init - postgres.New: %w", err)
	}
	a.pg = pg

	// init encrypter
	encr := encrypter.New(a.cfg.Encryption)
	err = encr.Init()
	if err != nil {
		logger.Error().Err(err).Msg("app - Init - encr.Init")
		return fmt.Errorf("app - Init - encr.Init: %w", err)
	}
	a.encr = encr

	return nil
}

// Run runs application.
func (a *App) Run() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	a.starter = starter.New(ctx)
	a.get = get.New(a.cfg, a.pg, a.encr)

	// Starting of HTTP-server
	a.starter.OnRun(a.get.Server().Run)

	// On terminate
	a.starter.OnShutdown(a.get.Server().Shutdown)

	err := a.starter.Run()
	if err != nil {
		logger.Error().Err(err).Msg("app - Init - starter.Run")
	}
}
