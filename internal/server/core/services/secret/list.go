package secret

import (
	"context"
	"fmt"

	"github.com/maybecoding/keep-it-safe/internal/server/core/entity"
)

// List - get list of secrets by user.
func (s *Service) List(ctx context.Context, userID entity.UserID) ([]entity.Secret, error) {
	list, err := s.store.SecretList(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("secret - SecretList - store.SecretList: %w", err)
	}
	return list, nil
}
