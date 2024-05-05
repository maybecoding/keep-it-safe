package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type CustomClaims struct {
	Username string
	Admin    bool
}

// TestEncodeDecode_Successful tests successful token encoding and decoding.
func TestEncodeDecode_Successful(t *testing.T) {
	secret := "secret"

	// Initialize the encode and decode functions with the example secret and expiry.
	encode, decode := Init[string, CustomClaims](secret, 24)

	// Create an instance of our claims.
	claims := CustomClaims{
		Username: "exampleUser",
		Admin:    false,
	}

	// Test encoding.
	tokenString, err := encode(claims)
	require.NoError(t, err, "Encoding should be successful")
	require.NotEmpty(t, tokenString, "Token string should not be empty")

	// Test decoding.
	decodedClaims, err := decode(tokenString)
	require.NoError(t, err, "Decoding should be successful")
	require.NotNil(t, decodedClaims, "Decoded claims should not be nil")
	assert.Equal(t, claims, *decodedClaims, "Decoded claims should match the original claims")
}

// TestDecode_InvalidToken simulates decoding failure due to invalid token.
func TestDecode_InvalidTokenString(t *testing.T) {
	secret := "secret"

	_, decode := Init[string, CustomClaims](secret, 24)

	// Provide an intentionally invalid token string.
	invalidTokenString := "invalidToken"

	// Test decoding.
	decodedClaims, err := decode(invalidTokenString)
	assert.Error(t, err, "Decoding should fail due to invalid token")
	assert.Nil(t, decodedClaims, "Decoded claims should be nil when token is invalid")
}
