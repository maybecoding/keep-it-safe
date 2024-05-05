package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/maybecoding/keep-it-safe/internal/core/entity"
	"github.com/maybecoding/keep-it-safe/pkg/pwd"
)

// Register - регистрация пользователя
//  1. Проверяем что имя пользователя доступно для добавления
//  2. Хэшируем пароль
//  3. Сохраняем пользователя
//  4. Формируем jwt-токен
func (s *Service) Register(ctx context.Context, login entity.UserLogin, password entity.UserPassword) (entity.Token, error) {
	// 1. Проверяем что имя пользователя доступно для добавления
	loginAvail, err := s.store.LoginAvailable(ctx, login)
	if err != nil {
		return "", fmt.Errorf("user - Register - order.LoginAvailable: %w", err)
	}
	if !loginAvail {
		return "", entity.ErrUserNotAvailable
	}

	// 2. Хэшируем пароль
	pwdHash, err := pwd.Hash([]byte(password))
	if err != nil {
		return "", fmt.Errorf("user - Register - pwd.Hash: %w", err)
	}

	// 3. Сохраняем пользователя
	user, err := s.store.UserNew(ctx, login, entity.UserPasswordHash(pwdHash))
	if err != nil {
		return "", fmt.Errorf("user - Register - order.UserNew: %w", err)
	}
	if user == nil {
		return "", errors.New("user - Register - user isn't created")
	}

	// 4. Формируем tokenData-токен
	tokenData := entity.TokenData{UserID: user.UserID}
	jwtToken, err := s.jwtEncode(tokenData)
	if err != nil {
		return "", fmt.Errorf("user - Register - pwd.Get: %w", err)
	}
	return jwtToken, nil
}
