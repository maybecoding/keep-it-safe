package encrypter

import (
	"encoding/hex"
	"fmt"

	"github.com/maybecoding/keep-it-safe/internal/server/adapters/encrypter/rotator"
)

// Init - initialize Encrypter
func (e *Encrypter) Init() error {
	// Check encryption key size
	switch e.cfg.KeySize {
	default:
		return fmt.Errorf("encrypter - Init - incorrect Key Size %d. Must be 16, 24 or 32", e.cfg.KeySize)
	case 16, 24, 32:
	}

	// Convert hex msk to bytes
	msk, err := hex.DecodeString(e.cfg.MasterKeyHex)
	if err != nil {
		return fmt.Errorf("encrypter - Init - hex.DecodeString: %w", err)
	}

	// if master key more then 32 bytes - error
	if len(msk) > 32 {
		return fmt.Errorf("encrypter - Init - master key size %d: must be max 32 byte", len(msk))
	}

	// if key len not in 16, 24 or 32 bytes adds padding
	switch len(msk) {
	default:
		for _, size := range []int{16, 24, 32} {
			if len(msk) < size {
				msk = append(msk, make([]byte, size-len(msk))...)
				fmt.Printf("master key size %d was padded to %d", len(msk), size)
				break
			}
		}
	case 16, 24, 32:
	}
	e.masterKey = msk

	// initialize encryptionKey rotator
	rotate := func() ([]byte, []byte, error) {
		key, err := rndGen(e.cfg.KeySize)
		if err != nil {
			return nil, nil, fmt.Errorf("encrypter - rotate - rndGen: %w", err)
		}
		keyEncr, err := encryptFix(e.masterKey, key)
		if err != nil {
			return nil, nil, fmt.Errorf("encrypter - rotate - encryptFix: %w", err)
		}
		return key, keyEncr, nil
	}
	e.encrKey, err = rotator.New(e.cfg.RotateDuration, rotate)
	if err != nil {
		return fmt.Errorf("encrypter - Init - rotator.New: %w", err)
	}
	return nil
}
