package secret

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/maybecoding/keep-it-safe/internal/server/core/entity"
)

// Set saves secret data.
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
	encrData, err := s.encr.Encrypt(dataB)
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
		Value:        encrData.Bytes,
		Nonce:        encrData.Nonce,
		EncryptionSK: encrData.EncryptionKeyEncrypted,
		Meta:         data.SecretMeta,
	}

	var secretID entity.SecretID
	err = s.store.WithTx(ctx, func(ctx context.Context) error {
		var err error
		secretID, err = s.store.SecretSet(ctx, scrDet)
		if err != nil {
			return fmt.Errorf("secret - Set - store.SecretSet: %w", err)
		}
		for _, attr := range scrDet.Meta {
			err = s.store.SecretAttrSet(ctx, secretID, attr)
			if err != nil {
				return fmt.Errorf("secret - Set - store.SecretAttrSet: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("secret - Set: %w", err)
	}

	return secretID, nil
}
