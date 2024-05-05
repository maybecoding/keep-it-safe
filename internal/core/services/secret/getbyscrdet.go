package secret

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/maybecoding/keep-it-safe/internal/core/entity"
)

func (s *Service) getByScrDet(scrDet entity.SecretDetail) (*entity.Data, error) {
	// decrypt
	dataB, err := s.encr.Decrypt(scrDet.Value, scrDet.Nonce, scrDet.EncryptionSK)
	if err != nil {
		return nil, fmt.Errorf("secret - get - encr.Decrypt: %w", err)
	}

	// prepare result
	data := entity.Data{
		SecretMeta: scrDet.Meta,
		SecretType: scrDet.Type,
	}
	// parse from bytes
	dec := gob.NewDecoder(bytes.NewBuffer(dataB))

	switch scrDet.Type {
	case entity.SecretTypeCredentials:
		err = dec.Decode(&data.Credentials)
	case entity.SecretTypeText:
		err = dec.Decode(&data.Text)
	case entity.SecretTypeBinary:
		err = dec.Decode(&data.Binary)
	case entity.SecretTypeBankCard:
		err = dec.Decode(&data.BankCard)
	default:
		err = entity.ErrSecretTypeUndefined
	}
	if err != nil {
		return nil, fmt.Errorf("secret - get - dec.Decode: %w", err)
	}

	return &data, nil
}
