package entity

import (
	"errors"
	"time"
)

type (
	SecretID   int32
	SecretName string
	SecretMeta []SecretAttr
)

type Secret struct {
	ID      SecretID
	UserID  UserID
	Type    SecretType
	Name    SecretName
	Created time.Time
	Updated time.Time
}

type SecretDetail struct {
	Secret
	Value        []byte
	Nonce        []byte
	EncryptionSK []byte
	Meta         SecretMeta
}

type SecretAttr struct {
	Attr  string
	Value string
}

var (
	ErrSecretNotFound      = errors.New("secret not found")
	ErrSecretForbiden      = errors.New("secret forbiden")
	ErrSecretTypeUndefined = errors.New("secret type undefined")
)
