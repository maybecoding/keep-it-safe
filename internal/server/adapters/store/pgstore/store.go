package pgstore

import "github.com/maybecoding/keep-it-safe/pkg/postgres"

type Store struct {
	pg *postgres.Postgres
}

func New(pg *postgres.Postgres) *Store {
	return &Store{pg: pg}
}
