package pgstore

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgxutil"
	"github.com/maybecoding/keep-it-safe/internal/server/core/entity"
)

func (s *Store) SecretList(ctx context.Context, userID entity.UserID) ([]entity.Secret, error) {
	secrets, err := pgxutil.Select(ctx,
		s.pg.Pool(ctx),
		`select id, usr_id, type, name, updated_at, created_at from scrt where usr_id = $1`,
		[]any{userID},
		pgx.RowToStructByPos[entity.Secret])
	if err != nil {
		return nil, fmt.Errorf("pgstore - SecretList - pgx.CollectRows: %w", err)
	}
	return secrets, nil
}
