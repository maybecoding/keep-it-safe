package get

import (
	"github.com/maybecoding/keep-it-safe/internal/adapters/api/v1/httpserver"
	"github.com/maybecoding/keep-it-safe/internal/adapters/encrypter"
	"github.com/maybecoding/keep-it-safe/internal/adapters/store/pgstore"
	"github.com/maybecoding/keep-it-safe/internal/config"
	"github.com/maybecoding/keep-it-safe/internal/core/entity"
	"github.com/maybecoding/keep-it-safe/internal/core/services/secret"
	"github.com/maybecoding/keep-it-safe/internal/core/services/user"
	"github.com/maybecoding/keep-it-safe/pkg/jwt"
	"github.com/maybecoding/keep-it-safe/pkg/postgres"
)

type Get struct {
	cfg  *config.Config
	pg   *postgres.Postgres
	encr *encrypter.Encrypter

	pgStore *pgstore.Store
	usrSrv  *user.Service
	scrtSrv *secret.Service
	server  *httpserver.Server
}

func New(cfg *config.Config, pg *postgres.Postgres, encr *encrypter.Encrypter) *Get {
	return &Get{cfg: cfg, pg: pg, encr: encr}
}

func (g *Get) PGStore() *pgstore.Store {
	if g.pgStore == nil {
		g.pgStore = pgstore.New(g.pg)
	}
	return g.pgStore
}

func (g *Get) UsrSrv() *user.Service {
	if g.usrSrv == nil {
		encode, decode := jwt.Init[entity.Token, entity.TokenData](g.cfg.JWT.Secret, g.cfg.JWT.ExpiresHours)
		g.usrSrv = user.New(g.PGStore(), encode, decode)
	}
	return g.usrSrv
}

func (g *Get) ScrtSrv() *secret.Service {
	if g.scrtSrv == nil {
		g.scrtSrv = secret.New(g.PGStore(), g.encr)
	}
	return g.scrtSrv
}

func (g *Get) Server() *httpserver.Server {
	if g.server == nil {
		g.server = httpserver.New(&g.cfg.HTTP, g.UsrSrv(), g.ScrtSrv())
	}
	return g.server
}
