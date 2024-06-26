package user

import (
	"fmt"

	"github.com/maybecoding/keep-it-safe/internal/server/core/entity"
)

// ValidateToken validates JWT token.
func (s *Service) ValidateToken(t entity.Token) (*entity.TokenData, error) {
	td, err := s.jwtDecode(t)
	if err != nil {
		return nil, fmt.Errorf("user - ValidateToken - uc.jwt.Use: %v", err)
	}
	return td, nil
}
