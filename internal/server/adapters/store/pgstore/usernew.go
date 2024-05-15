// Package pgstore storage.
package pgstore

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/maybecoding/keep-it-safe/internal/server/core/entity"
)

// UserNew creates new user.
func (s *Store) UserNew(ctx context.Context, login entity.UserLogin, hash entity.UserPasswordHash) (*entity.User, error) {
	usr := entity.User{}
	query := `insert into usr(login, hash) values(@login, @hash) returning id, login, hash;`
	err := s.pg.Pool(ctx).QueryRow(ctx, query, pgx.NamedArgs{
		"login": login,
		"hash":  hash,
	}).Scan(&usr.UserID, &usr.UserLogin, &usr.UserPasswordHash)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, entity.ErrUserNotAvailable
		}
		// var pgErr *pgconn.PgError replace with code if bug will be fixed
		// if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		// 	return nil, entity.ErrUserNotAvailable
		// }
		return nil, fmt.Errorf("pgstore - UserNew - ar.Pool(ctx).QueryRow: %w", err)
	}

	return &usr, nil
}
