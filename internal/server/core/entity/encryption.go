package entity

// EncryptionData struct for: source for decryption / decryption result.
type EncryptionData struct {
	Bytes                  []byte
	Nonce                  []byte
	EncryptionKeyEncrypted []byte
}
