package pgstore

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/maybecoding/keep-it-safe/internal/server/core/entity"
)

// UserGet gets user by login.
func (s *Store) UserGet(ctx context.Context, login entity.UserLogin) (*entity.User, error) {
	usr := entity.User{}
	query := `select id, login, hash from  usr where login = $1`
	err := s.pg.Pool(ctx).QueryRow(ctx, query, login).Scan(&usr.UserID, &usr.UserLogin, &usr.UserPasswordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrUserNotFound
		}
		return nil, fmt.Errorf("pgstore - UserNew - ar.Pool(ctx).QueryRow: %w", err)
	}
	return &usr, nil
}
