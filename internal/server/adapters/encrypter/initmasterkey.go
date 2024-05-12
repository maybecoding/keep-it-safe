package encrypter

import (
	"encoding/hex"
	"fmt"
	"slices"
)

func (e *Encrypter) initMasterKey() error {
	desiredSizes := []int{16, 24, 32}

	// Check encryption key size.
	if !slices.Contains(desiredSizes, e.cfg.KeySize) {
		return fmt.Errorf("encrypter - InitMasterKey - incorrect Key Size %d. Must be 16, 24 or 32", e.cfg.KeySize)
	}

	// Convert hex msk to bytes.
	msk, err := hex.DecodeString(e.cfg.MasterKeyHex)
	if err != nil {
		return fmt.Errorf("encrypter - InitMasterKey - hex.DecodeString: %w", err)
	}

	// Pad master key if size not desired.
	msk, err = padMasterKey(msk, desiredSizes)
	if err != nil {
		return fmt.Errorf("encrypter - InitMasterKey: %w", err)
	}
	e.masterKey = msk

	return nil
}
