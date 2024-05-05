package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Postgres struct {
	pool *pgxpool.Pool
}

// New - Подключение к БД создание пула соединений
func New(uri string) (*Postgres, error) {
	pg := &Postgres{}
	ctx := context.Background()
	pgPool, err := pgxpool.New(context.Background(), uri)
	if err != nil {
		return nil, fmt.Errorf("postgres - New - pgxpool.New: %w", err)
	}
	pg.pool = pgPool
	err = pgPool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("postgres - New - Ping: %w", err)
	}
	return pg, nil
}

// Close - Закрытие соединения
func (pg *Postgres) Close() {
	if pg.pool != nil {
		pg.pool.Close()
	}
}
