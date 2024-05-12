package encrypter

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/maybecoding/keep-it-safe/internal/server/core/entity"
)

// Encrypt data and returns encrypted data, nonce, encrypted encription key and error if it ocured.
func (e *Encrypter) Encrypt(src []byte) (result entity.EncryptionData, err error) {
	// Get rotated encryption keys
	var encrKey []byte
	encrKey, result.EncryptionKeyEncrypted, err = e.encrKey.Keys()
	if err != nil {
		return result, fmt.Errorf("encrypter - Encrypt - encrKey.Keys: %w", err)
	}

	// Encrypt using GCM with AES
	// New AES
	aesblock, err := aes.NewCipher(encrKey)
	if err != nil {
		return result, fmt.Errorf("encrypter - Encrypt - aes.NewCipher: %w", err)
	}

	// New GCM
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return result, fmt.Errorf("encrypter - Encrypt - cipher.NewGCM: %w", err)
	}

	// Initialization vector
	result.Nonce, err = rndGen(aesgcm.NonceSize())
	if err != nil {
		return result, fmt.Errorf("encrypter - Encrypt - rndGen: %w", err)
	}

	// Encrypt
	result.Bytes = aesgcm.Seal(nil, result.Nonce, src, nil)

	return result, nil
}
