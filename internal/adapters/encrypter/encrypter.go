package encrypter

import (
	"crypto/aes"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/maybecoding/keep-it-safe/internal/adapters/encrypter/rotator"
)

// Encrypter - struct for encrypt and decrypt
type Encrypter struct {
	keySize      int
	masterKeyHex string
	encrKeyDur   time.Duration

	masterKey []byte
	encrKey   *rotator.Rotator
}

func New(msk string, encrKeyDur time.Duration) *Encrypter {
	return &Encrypter{masterKeyHex: msk, encrKeyDur: encrKeyDur}
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
	for i := 0; i < blockCnt; i += 1 {
		currSrc := src[i*aes.BlockSize : (i+1)*aes.BlockSize]
		currDst := dst[i*aes.BlockSize : (i+1)*aes.BlockSize]
		aesblock.Encrypt(currDst, currSrc)
	}

	// dst := make([]byte, len(src))
	// aesblock.Encrypt(dst, src)
	return dst, nil
}

func decryptFix(key []byte, src []byte) ([]byte, error) {
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("fixDecrypt - aes.NewCipher: %w", err)
	}

	srcLen := len(src)

	dst := make([]byte, srcLen)
	for i := 0; i < srcLen/aes.BlockSize; i += 1 {
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
