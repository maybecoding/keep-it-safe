package pgstore

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgxutil"
	"github.com/maybecoding/keep-it-safe/internal/server/core/entity"
)

func (s *Store) SecretGet(ctx context.Context, secretID entity.SecretID) (*entity.SecretDetail, error) {
	// time.Sleep(1 * time.Minute)
	secret, err := pgxutil.SelectRow(ctx, s.pg.Pool(ctx), `
	select id, usr_id, type, name, created_at, updated_at, value, nonce, encryption_sk, (
		select array_agg(row(attr.name, val.value) order by attr.id)
		from scrt_attr val
			inner join attr
				on attr.id = val.attr_id
		where val.scrt_id = $1
	)
	from scrt
	where id = $1
	`,
		[]any{secretID},
		pgx.RowToStructByPos[entity.SecretDetail])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, entity.ErrSecretNotFound
		}
		return nil, fmt.Errorf("pgstore - SecretList - pgx.CollectRows: %w", err)
	}

	return &secret, nil
}
