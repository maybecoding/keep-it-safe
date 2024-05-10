package secret

import (
	"context"

	"github.com/maybecoding/keep-it-safe/internal/server/core/entity"
)

// Store interface for storage used for service.
type Store interface {
	WithTx(ctx context.Context, fn func(context.Context) error) error
	SecretList(ctx context.Context, userID entity.UserID) ([]entity.Secret, error)
	SecretGet(ctx context.Context, secretID entity.SecretID) (*entity.SecretDetail, error)
	SecretSet(ctx context.Context, scrt entity.SecretDetail) (entity.SecretID, error)
	SecretAttrSet(ctx context.Context, secretID entity.SecretID, attr entity.SecretAttr) error
}

// Encrypter interface for encryption functions used in service.
type Encrypter interface {
	Decrypt(src []byte, nonce []byte, encrKeyEncrypted []byte) (dst []byte, err error)
	Encrypt(src []byte) (dst []byte, nonce []byte, keyEncr []byte, err error)
}

// Service struct of service.
type Service struct {
	store Store
	encr  Encrypter
}

// New creates new service.
func New(store Store, encr Encrypter) *Service {
	u := Service{store: store, encr: encr}
	return &u
}
