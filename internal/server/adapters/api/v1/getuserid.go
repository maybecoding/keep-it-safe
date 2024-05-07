package api

import (
	"fmt"
	"strings"

	"github.com/maybecoding/keep-it-safe/internal/server/adapters/api/v1/models"
	"github.com/maybecoding/keep-it-safe/internal/server/core/entity"
)

func (a *API) getUserID(authCookie string) (entity.UserID, *models.Error) {
	auth := strings.Split(authCookie, " ")
	if len(auth) != 2 || auth[0] != "Bearer" {
		return 0, &models.Error{Error: "authorization cookie isn't provided"}
	}
	token := entity.Token(auth[1])

	tokenData, err := a.user.ValidateToken(token)
	if err != nil {
		return 0, &models.Error{Error: fmt.Sprintf("token invalid: %s", err.Error())}
	}
	return tokenData.UserID, nil
}
