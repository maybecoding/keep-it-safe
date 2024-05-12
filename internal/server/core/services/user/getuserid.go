package user

import (
	"strings"

	"github.com/maybecoding/keep-it-safe/generated/models"
	"github.com/maybecoding/keep-it-safe/internal/server/core/entity"
)

func (s *Service) GetUserID(authCookie string) (entity.UserID, *models.Error) {
	auth := strings.Split(authCookie, " ")
	if len(auth) != 2 || auth[0] != "Bearer" {
		return 0, &models.Error{Error: "authorization cookie isn't provided"}
	}
	token := entity.Token(auth[1])

	tokenData, err := s.ValidateToken(token)
	if err != nil {
		return 0, &models.Error{Error: "token invalid:" + err.Error()}
	}
	return tokenData.UserID, nil
}
