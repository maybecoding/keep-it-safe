package entity

// SecretType type of secret.
type SecretType int32

const (
	// SecretTypeCredentials Type for Credentials.
	SecretTypeCredentials SecretType = iota
	// SecretTypeText Type for Credentials.
	SecretTypeText
	// SecretTypeBinary Type for Credentials.
	SecretTypeBinary
	// SecretTypeBankCard Type for Credentials.
	SecretTypeBankCard
)

// String returns text presentation of secret type.
func (st SecretType) String() string {
	switch st {
	case SecretTypeCredentials:
		return "Credentials"
	case SecretTypeText:
		return "Text"
	case SecretTypeBinary:
		return "Binary"
	case SecretTypeBankCard:
		return "BankCard"
	default:
		return "Undefined"
	}
}
