package auth

import "golang.org/x/crypto/bcrypt"

// HashSecret hashes a plaintext secret using bcrypt.
func HashSecret(secret string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(secret), 14)
	return string(bytes), err
}

// CheckSecretHash compares a plaintext secret with a bcrypt hash.
func CheckSecretHash(secret, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(secret))
	return err == nil
}
