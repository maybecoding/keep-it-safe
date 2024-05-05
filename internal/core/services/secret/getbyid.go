package secret

import (
	"context"
	"fmt"

	"github.com/maybecoding/keep-it-safe/internal/core/entity"
)

// GetByID - get secret of user by id
func (s *Service) GetByID(ctx context.Context, userID entity.UserID, secretID entity.SecretID) (*entity.Data, error) {
	// get secret from DB
	scrDet, err := s.store.SecretGet(ctx, secretID)
	if err != nil {
		return nil, fmt.Errorf("secret - GetByID - store.SecretGet: %w", err)
	}
	if scrDet.UserID != userID {
		return nil, entity.ErrSecretForbiden
	}
	return s.getByScrDet(*scrDet)
}
