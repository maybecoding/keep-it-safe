// Package encrypter used for encryption on data.
package encrypter

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/maybecoding/keep-it-safe/internal/server/core/entity"
)

// Decrypt using src data, src data Vector and encrypted encryption key, returns decrypted data or error.
func (e *Encrypter) Decrypt(src entity.EncryptionData) (dst []byte, err error) {
	// using master key we must decrypt enctyption key
	encrKey, err := decryptFix(e.masterKey, src.EncryptionKeyEncrypted)
	if err != nil {
		return nil, fmt.Errorf("encrypter - Decrypt - decryptFix: %w", err)
	}

	// Decrypt using GCM with AES
	// New AES
	aesblock, err := aes.NewCipher(encrKey)
	if err != nil {
		return nil, fmt.Errorf("encrypter - Decrypt - aes.NewCipher: %w", err)
	}

	// New GCM
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, fmt.Errorf("encrypter - Decrypt - cipher.NewGCM: %w", err)
	}

	// Decrypt
	dst, err = aesgcm.Open(nil, src.Nonce, src.Bytes, nil)
	if err != nil {
		return nil, fmt.Errorf("encrypter - Decrypt - aesgcm.Open: %w", err)
	}

	return dst, nil
}
