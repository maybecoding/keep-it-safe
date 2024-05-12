package encrypter

import (
	"fmt"

	"github.com/maybecoding/keep-it-safe/pkg/logger"
)

func padMasterKey(msk []byte, desiredSizes []int) ([]byte, error) {
	currentLength := len(msk)

	for _, size := range desiredSizes {
		if currentLength == size {
			return msk, nil
		}
		if currentLength < size {
			padding := make([]byte, size-currentLength)
			msk = append(msk, padding...)
			logger.Info().Msgf("master key size %d was padded to %d\n", currentLength, size)
			return msk, nil
		}
	}

	// if master key more then 32 bytes - error
	return nil, fmt.Errorf("encrypter - padMasterKey - master key size %d: must be max 32 byte", len(msk))
}
