package secret

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/maybecoding/keep-it-safe/internal/core/entity"
)

func (s *Service) Set(ctx context.Context, userID entity.UserID, data entity.Data) (entity.SecretID, error) {
	// convert data to bytes
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)

	var err error
	switch data.SecretType {
	case entity.SecretTypeCredentials:
		err = enc.Encode(&data.Credentials)
	case entity.SecretTypeText:
		err = enc.Encode(&data.Text)
	case entity.SecretTypeBinary:
		err = enc.Encode(&data.Binary)
	case entity.SecretTypeBankCard:
		err = enc.Encode(&data.BankCard)
	default:
		err = entity.ErrSecretTypeUndefined
	}
	if err != nil {
		return 0, fmt.Errorf("secret - Set - enc.Encode: %w", err)
	}
	dataB := buf.Bytes()

	// encrypt data
	dst, nonce, keyEncr, err := s.encr.Encrypt(dataB)
	if err != nil {
		return 0, fmt.Errorf("secret - Set - encr.Encrypt: %w", err)
	}

	// store to db
	scrDet := entity.SecretDetail{
		Secret: entity.Secret{
			UserID: userID,
			Type:   data.SecretType,
			Name:   data.SecretName,
		},
		Value:        dst,
		Nonce:        nonce,
		EncryptionSK: keyEncr,
		Meta:         data.SecretMeta,
	}
	secretID, err := s.store.SecretSet(ctx, scrDet)
	if err != nil {
		return 0, fmt.Errorf("secret - Set - store.SecretSet: %w", err)
	}
	return secretID, nil
}
