package entity

type SecretType int32

const (
	SecretTypeCredentials SecretType = iota
	SecretTypeText
	SecretTypeBinary
	SecretTypeBankCard
)

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
