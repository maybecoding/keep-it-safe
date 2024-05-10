// Package pwd - helps hash password and check if password matches.
package pwd

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Hash - hash password.
func Hash(pwd []byte, cost ...int) ([]byte, error) {
	c := bcrypt.DefaultCost
	if len(cost) > 0 {
		c = cost[0]
	}
	hash, err := bcrypt.GenerateFromPassword(pwd, c)
	if err != nil {
		return nil, fmt.Errorf("pwd - Hash - bcrypt.GenerateFromPassword: %w", err)
	}
	return hash, nil
}

// Check - checks password corresponds to its hash.
func Check(pwd []byte, hash []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, pwd)
	return err == nil
}
