// Package get used for creating main struts with automating dependencies.
package get

import (
	"github.com/maybecoding/keep-it-safe/internal/server/adapters/api/v1/httpserver"
	"github.com/maybecoding/keep-it-safe/internal/server/adapters/encrypter"
	"github.com/maybecoding/keep-it-safe/internal/server/adapters/store/pgstore"
	"github.com/maybecoding/keep-it-safe/internal/server/config"
	"github.com/maybecoding/keep-it-safe/internal/server/core/entity"
	"github.com/maybecoding/keep-it-safe/internal/server/core/services/secret"
	"github.com/maybecoding/keep-it-safe/internal/server/core/services/user"
	"github.com/maybecoding/keep-it-safe/pkg/jwt"
	"github.com/maybecoding/keep-it-safe/pkg/postgres"
)

// Get struct with main structs.
type Get struct {
	cfg  *config.Config
	pg   *postgres.Postgres
	encr *encrypter.Encrypter

	pgStore *pgstore.Store
	usrSrv  *user.Service
	scrtSrv *secret.Service
	server  *httpserver.Server
}

// New creates new Get struct.
func New(cfg *config.Config, pg *postgres.Postgres, encr *encrypter.Encrypter) *Get {
	return &Get{cfg: cfg, pg: pg, encr: encr}
}

// PGStore returns store.
func (g *Get) PGStore() *pgstore.Store {
	if g.pgStore == nil {
		g.pgStore = pgstore.New(g.pg)
	}
	return g.pgStore
}

// UsrSrv returns user service.
func (g *Get) UsrSrv() *user.Service {
	if g.usrSrv == nil {
		encode, decode := jwt.Init[entity.Token, entity.TokenData](g.cfg.JWT.Secret, g.cfg.JWT.ExpiresHours)
		g.usrSrv = user.New(g.PGStore(), encode, decode)
	}
	return g.usrSrv
}

// ScrtSrv returns secret service.
func (g *Get) ScrtSrv() *secret.Service {
	if g.scrtSrv == nil {
		g.scrtSrv = secret.New(g.PGStore(), g.encr)
	}
	return g.scrtSrv
}

// Server returns http server.
func (g *Get) Server() *httpserver.Server {
	if g.server == nil {
		g.server = httpserver.New(&g.cfg.HTTP, g.UsrSrv(), g.ScrtSrv())
	}
	return g.server
}
