package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/maybecoding/keep-it-safe/generated/server"
	"github.com/maybecoding/keep-it-safe/internal/server/adapters/api/v1"
	"github.com/maybecoding/keep-it-safe/internal/server/config"
	"github.com/maybecoding/keep-it-safe/internal/server/core/services/secret"
	"github.com/maybecoding/keep-it-safe/internal/server/core/services/user"
	"github.com/maybecoding/keep-it-safe/pkg/logger"

	middleware "github.com/oapi-codegen/nethttp-middleware"
)

type Server struct {
	cfg    *config.HTTP
	user   *user.Service
	secret *secret.Service

	server *http.Server
}

func New(cfg *config.HTTP, u *user.Service, s *secret.Service) *Server {
	return &Server{cfg: cfg, user: u, secret: s}
}

func (s *Server) Run(_ context.Context) error {
	logger.Info().Msg("Starting HTTP server")
	swagger, err := server.GetSwagger()
	if err != nil {
		return fmt.Errorf("http - New - api.GetSwagger: %w", err)
	}

	swgAPI := api.New(s.user, s.secret)

	strictHandler := server.NewStrictHandler(swgAPI, nil)
	h := server.HandlerFromMux(strictHandler, http.NewServeMux())

	handler := middleware.OapiRequestValidator(swagger)(h)
	s.server = &http.Server{Addr: s.cfg.Address, Handler: handler}

	return fmt.Errorf("http - Run - server.ListenAndServe: %w", s.server.ListenAndServe())
}

func (s *Server) Shutdown(_ context.Context) error {
	logger.Info().Msg("Stopping HTTP server")
	return fmt.Errorf("http - Shutdown - server.Shutdown:%w", s.server.Shutdown(context.Background()))
}
