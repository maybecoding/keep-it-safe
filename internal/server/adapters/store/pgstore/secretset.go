package pgstore

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/maybecoding/keep-it-safe/internal/server/core/entity"
)

// SecretSet idempotental saves secret.
func (s *Store) SecretSet(ctx context.Context, scrt entity.SecretDetail) (entity.SecretID, error) {
	now := time.Now()
	var secretID entity.SecretID
	err := s.pg.Pool(ctx).QueryRow(ctx, `insert into scrt (usr_id, type, name, value, nonce, encryption_sk, created_at, updated_at)
	values (@usr_id, (@type), @name, @value, @nonce, @encryption_sk, @created_at, @updated_at)
	on conflict (usr_id, name) do update set type = @type, value = @value, nonce = @nonce, encryption_sk = @encryption_sk, updated_at = @updated_at
	returning id`, pgx.NamedArgs{
		"usr_id":        scrt.UserID,
		"type":          scrt.Type,
		"name":          scrt.Name,
		"value":         scrt.Value,
		"nonce":         scrt.Nonce,
		"encryption_sk": scrt.EncryptionSK,
		"created_at":    now,
		"updated_at":    now,
	}).Scan(&secretID)
	if err != nil {
		return secretID, fmt.Errorf("pgstore - SecretSet: %w", err)
	}
	return secretID, nil
}
