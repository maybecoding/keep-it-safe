// Package entity stores main entities of server app.
package entity

// DataCredentials stores login and password.
type DataCredentials struct {
	Login    string
	Password string
}

// DataText stores text.
type DataText string

// DataBinary stores bytes.
type DataBinary []byte

// DataBankCard stores bank card data.
type DataBankCard struct {
	Number         string
	Valid          string
	Holder         string
	ValidationCode string
}

// Data stores secret with data of SecretType.
type Data struct {
	SecretName SecretName
	SecretMeta SecretMeta
	SecretType SecretType

	Credentials *DataCredentials
	Text        *DataText
	Binary      *DataBinary
	BankCard    *DataBankCard
}
