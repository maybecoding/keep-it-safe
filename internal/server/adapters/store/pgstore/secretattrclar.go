package pgstore

import (
	"context"
	"fmt"

	"github.com/maybecoding/keep-it-safe/internal/server/core/entity"
)

func (s *Store) SecretAttrClear(ctx context.Context, secretID entity.SecretID) error {
	_, err := s.pg.Pool(ctx).Exec(ctx, attrClearQuery, secretID)
	if err != nil {
		return fmt.Errorf("pgstore - SecretAttrSet - failed insert into attr_value: %w", err)
	}
	return nil
}

var attrClearQuery = `delete from attr_value where scrt_id = $1`
