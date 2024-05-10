package encrypter

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncryptDecryptFix(t *testing.T) {
	// for i := 0; i < 10; i++ {
	for i := range 10 {
		for _, size := range []int{16, 24, 32} {
			key, err := rndGen(size)
			require.NoError(t, err)

			src, err := rndGen(size)
			require.NoError(t, err)

			dest, err := encryptFix(key, src)
			require.NoError(t, err)
			require.NotEqual(t, src, dest)

			srcDecr, err := decryptFix(key, dest)
			require.NoError(t, err)
			srcDecr = srcDecr[:size]
			require.Equal(t, src, srcDecr, "source and decrypted destination are not equal for iteration %d and size %d", i, size)
		}
	}
}
