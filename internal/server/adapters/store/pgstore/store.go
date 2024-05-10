package pgstore

import "github.com/maybecoding/keep-it-safe/pkg/postgres"

// Store struct for store.
type Store struct {
	pg *postgres.Postgres
}

// New creates new store.
func New(pg *postgres.Postgres) *Store {
	return &Store{pg: pg}
}
