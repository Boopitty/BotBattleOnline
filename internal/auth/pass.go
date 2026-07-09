package auth

import (
	"github.com/alexedwards/argon2id"
)

// Change a password into a hash before it's stored in the database
func HashPass(password string) (string, error) {
	param := argon2id.DefaultParams
	hash, err := argon2id.CreateHash(password, param)
	if err != nil {
		return "", err
	}
	return hash, nil
}

// Check the validity of a hashed password
func CheckPassHash(password, hash string) (bool, error) {
	check, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return check, nil
}
