package user

import (
	"context"

	"github.com/maybecoding/keep-it-safe/internal/core/entity"
)

type Store interface {
	UserGet(ctx context.Context, login entity.UserLogin) (*entity.User, error)
	LoginAvailable(ctx context.Context, login entity.UserLogin) (loginAvail bool, err error)
	UserNew(ctx context.Context, login entity.UserLogin, hash entity.UserPasswordHash) (*entity.User, error)
}

type Service struct {
	store     Store
	jwtEncode func(entity.TokenData) (entity.Token, error)
	jwtDecode func(entity.Token) (*entity.TokenData, error)
}

func New(
	store Store,
	encode func(entity.TokenData) (entity.Token, error),
	decode func(entity.Token) (*entity.TokenData, error),
) *Service {
	u := Service{store: store, jwtEncode: encode, jwtDecode: decode}
	return &u
}
