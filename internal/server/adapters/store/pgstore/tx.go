package pgstore

import (
	"context"
)

// WithTx injects transaction into tx.
func (s *Store) WithTx(ctx context.Context, fn func(context.Context) error) error {
	return s.pg.WithTx(ctx, fn)
}
