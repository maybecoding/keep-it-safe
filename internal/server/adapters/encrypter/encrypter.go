package encrypter

import (
	"crypto/aes"
	"crypto/rand"
	"fmt"

	"github.com/maybecoding/keep-it-safe/internal/server/adapters/encrypter/rotator"
	"github.com/maybecoding/keep-it-safe/internal/server/config"
)

// Encrypter - struct for encrypt and decrypt.
type Encrypter struct {
	masterKey []byte
	encrKey   *rotator.Rotator

	cfg config.Encryption
}

func New(cfg config.Encryption) *Encrypter {
	return &Encrypter{cfg: cfg}
}

func encryptFix(key []byte, src []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher: %w", err)
	}

	srcLen := len(src)
	blockCnt := (srcLen + aes.BlockSize - 1) / aes.BlockSize
	src = append(src, make([]byte, aes.BlockSize*blockCnt-srcLen)...)

	dst := make([]byte, blockCnt*aes.BlockSize) // зашифровываем
	for i := range blockCnt {
		currSrc := src[i*aes.BlockSize : (i+1)*aes.BlockSize]
		currDst := dst[i*aes.BlockSize : (i+1)*aes.BlockSize]
		aesblock.Encrypt(currDst, currSrc)
	}

	return dst, nil
}

func decryptFix(key []byte, src []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("fixDecrypt - aes.NewCipher: %w", err)
	}

	srcLen := len(src)

	dst := make([]byte, srcLen)
	for i := range srcLen / aes.BlockSize {
		currSrc := src[i*aes.BlockSize : (i+1)*aes.BlockSize]
		currDst := dst[i*aes.BlockSize : (i+1)*aes.BlockSize]
		aesblock.Decrypt(currDst, currSrc)
	}

	// dst := make([]byte, len(src))
	// aesblock.Decrypt(dst, src)
	return dst, nil
}

func rndGen(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("rndGen - rand.Read: %w", err)
	}
	return b, nil
}
