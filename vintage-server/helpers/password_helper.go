package helpers

import "golang.org/x/crypto/bcrypt"

// GeneratePasswordHash menerima password plaintext, mengembalikan hash string
func GeneratePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePassword menerima password plaintext dan hash, mengembalikan true kalau cocok
func ComparePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
