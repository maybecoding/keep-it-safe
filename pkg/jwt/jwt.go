// Package jwt used for JWT token check and generating.
package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Init returns 2 functions: 1. Encode - converts token data to token string 2. Decode - decodes token to token data.
func Init[T ~string, D any](secret string, expiresHours int) (encode func(D) (T, error), decode func(T) (*D, error)) {
	type claims struct {
		TokenData D
		jwt.RegisteredClaims
	}

	encode = func(tokenData D) (T, error) {
		c := claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expiresHours))),
			},
			TokenData: tokenData,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

		tokenStr, err := token.SignedString([]byte(secret))
		if err != nil {
			return "", fmt.Errorf("jwt - Encode - token.SignedString: %w", err)
		}
		return T(tokenStr), nil
	}
	decode = func(t T) (*D, error) {
		c := &claims{}
		token, err := jwt.ParseWithClaims(string(t), c, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("jwt - Decode - unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})
		if err != nil {
			return nil, fmt.Errorf("jwt - Decode - jwt.Parse: %w", err)
		}

		if !token.Valid {
			return nil, errors.New("jwt - Decode - token is not valid")
		}
		return &c.TokenData, nil
	}

	return
}
