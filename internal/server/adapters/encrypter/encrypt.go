package encrypter

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
)

// Encrypt data and returns encrypted data, nonce, encrypted encription key and error if it ocured.
func (e *Encrypter) Encrypt(src []byte) (dst []byte, nonce []byte, keyEncr []byte, err error) {
	// Get rotated encryption keys
	encrKey, encrKeyEncrypted, err := e.encrKey.Keys()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("encrypter - Encrypt - encrKey.Keys: %w", err)
	}

	// Encrypt using GCM with AES
	// New AES
	aesblock, err := aes.NewCipher(encrKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("encrypter - Encrypt - aes.NewCipher: %w", err)
	}

	// New GCM
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("encrypter - Encrypt - cipher.NewGCM: %w", err)
	}

	// Initialization vector
	nonce, err = rndGen(aesgcm.NonceSize())
	if err != nil {
		return nil, nil, nil, fmt.Errorf("encrypter - Encrypt - rndGen: %w", err)
	}

	// Encrypt
	dst = aesgcm.Seal(nil, nonce, src, nil)

	return dst, nonce, encrKeyEncrypted, nil
}
