package pgstore

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/maybecoding/keep-it-safe/internal/server/core/entity"
)

func (s *Store) LoginAvailable(ctx context.Context, login entity.UserLogin) (loginAvail bool, err error) {
	err = s.pg.Pool(ctx).QueryRow(ctx, `select not exists (select 1 from usr where login = @login)`, pgx.NamedArgs{
		"login": login,
	}).Scan(&loginAvail)
	if err != nil {
		return false, fmt.Errorf("pgstore - LoginAvailable - ar.Pool(ctx).QueryRow: %w", err)
	}
	return loginAvail, nil
}
