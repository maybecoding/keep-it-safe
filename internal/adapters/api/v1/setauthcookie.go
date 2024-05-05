package api

import (
	"net/http"

	"github.com/maybecoding/keep-it-safe/internal/core/entity"
)

func authCookie(token entity.Token) string {
	c := &http.Cookie{
		Name:  "Authorization",
		Value: "Bearer " + string(token),
	}
	return c.String()
}
