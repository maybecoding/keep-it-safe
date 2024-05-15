package encrypter

import (
	"fmt"

	"github.com/maybecoding/keep-it-safe/internal/server/adapters/encrypter/rotator"
)

// Init - initialize Encrypter.
func (e *Encrypter) Init() error {
	// initialize master key.
	err := e.initMasterKey()
	if err != nil {
		return fmt.Errorf("encrypter - Init: %w", err)
	}

	// initialize encryptionKey rotator.
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
