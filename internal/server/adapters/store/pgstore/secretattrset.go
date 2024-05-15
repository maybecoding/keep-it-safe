package pgstore

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/maybecoding/keep-it-safe/internal/server/core/entity"
)

// SecretAttrSet - sets secret attr.
func (s *Store) SecretAttrSet(ctx context.Context, secretID entity.SecretID, attr entity.SecretAttr) error {
	// create meta attr if not exists
	var attrID int32
	err := s.pg.Pool(ctx).QueryRow(ctx, attrSelectQuery, attr.Attr).Scan(&attrID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("pgstore - SecretAttrSet - failed to get attr_id: %w", err)
		}
		err = s.pg.Pool(ctx).QueryRow(ctx, attrInsertQuery, attr.Attr).Scan(&attrID)
		if err != nil {
			return fmt.Errorf("pgstore - SecretAttrSet - failed to insert attr: %w", err)
		}
	}
	_, err = s.pg.Pool(ctx).Exec(ctx, attrValueInsertQuery, secretID, attrID, attr.Value)
	if err != nil {
		return fmt.Errorf("pgstore - SecretAttrSet - failed insert into attr_value: %w", err)
	}
	return nil
}

var (
	attrSelectQuery = `select id from attr where name = $1;`
	attrInsertQuery = `insert into attr (name) values ($1)
on conflict(name) do nothing returning id;`
)

var attrValueInsertQuery = `insert into scrt_attr (scrt_id, attr_id, value)
values ($1, $2, $3)
on conflict (scrt_id, attr_id) do update set value = $3`
