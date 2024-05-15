// Package secret - service for users.
package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/maybecoding/keep-it-safe/internal/server/core/entity"
	"github.com/maybecoding/keep-it-safe/pkg/pwd"
)

// Login - Аутентификация пользователя
// 1. Получаем по логину UserID и Hash из БД
// 2. Проверяем Хэш
// 3. Формируем jwt-токен.
func (s *Service) Login(ctx context.Context, login entity.UserLogin, password entity.UserPassword) (entity.Token, error) {
	// 1. Получаем по логину UserID и Hash из БД
	user, err := s.store.UserGet(ctx, login)
	if err != nil {
		if errors.Is(err, entity.ErrUserNotFound) {
			return "", entity.ErrUserNotFound
		}
		return "", fmt.Errorf("user - Login - store.UserGet: %w", err)
	}
	if user == nil {
		return "", entity.ErrUserNotFound
	}

	// 2. Проверяем Хэш
	pwdOk := pwd.Check([]byte(password), []byte(user.UserPasswordHash))
	if !pwdOk {
		return "", entity.ErrIncorrectPassword
	}

	// 3. Формируем jwt-токен
	tokenData := entity.TokenData{UserID: user.UserID}
	jwtToken, err := s.jwtEncode(tokenData)
	if err != nil {
		return "", fmt.Errorf("user - Login - pwd.Get: %w", err)
	}
	return jwtToken, nil
}
