package encrypter

import (
	"encoding/hex"
	"fmt"

	"github.com/maybecoding/keep-it-safe/internal/adapters/encrypter/rotator"
)

// Init - initialize Encrypter
func (e *Encrypter) Init() error {
	// Convert hex msk to bytes
	msk, err := hex.DecodeString(e.masterKeyHex)
	if err != nil {
		return fmt.Errorf("encrypter - Init - hex.DecodeString: %w", err)
	}
	if len(msk) != 32 && len(msk) != 24 && len(msk) != 16 {
		return fmt.Errorf("provided key have size %d but it must be 32 or 24 or 16", len(msk))
	}
	e.masterKey = msk
	e.keySize = len(msk)

	// initialize encryptionKey rotator
	rotate := func() ([]byte, []byte, error) {
		key, err := rndGen(e.keySize)
		if err != nil {
			return nil, nil, fmt.Errorf("encrypter - rotate - rndGen: %w", err)
		}
		keyEncr, err := encryptFix(e.masterKey, key)
		if err != nil {
			return nil, nil, fmt.Errorf("encrypter - rotate - encryptFix: %w", err)
		}
		return key, keyEncr, nil
	}
	e.encrKey, err = rotator.New(e.encrKeyDur, rotate)
	if err != nil {
		return fmt.Errorf("encrypter - Init - rotator.New: %w", err)
	}
	return nil
}
