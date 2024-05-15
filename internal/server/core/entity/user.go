package entity

import "errors"

type (
	// User Type for user.
	User struct {
		UserID
		UserLogin
		UserPasswordHash
	}

	// UserID Type for UserID.
	UserID int32
	// UserLogin Type for UserLogin.
	UserLogin string
	// UserPasswordHash type for user password hash.
	UserPasswordHash string
	// UserPassword type for user password.
	UserPassword string
	// Token type for user token.
	Token string
	// TokenData type for Token data.
	TokenData struct {
		UserID UserID
	}
)

var (
	ErrUserNotAvailable  = errors.New("login already taken")
	ErrUserNotFound      = errors.New("user not found")
	ErrIncorrectPassword = errors.New("password is incorrect")
)
