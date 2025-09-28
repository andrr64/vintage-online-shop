// File: pkg/hash/bcrypt.go
package hash

import "golang.org/x/crypto/bcrypt"

// Generate membuat hash dari password menggunakan bcrypt.
func Generate(password string) (string, error) {
	// bcrypt.GenerateFromPassword menerima password dalam bentuk byte slice.
	// bcrypt.DefaultCost adalah standar untuk tingkat kesulitan hashing (cost factor).
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Verify membandingkan password plaintext dengan hash yang sudah ada.
// Mengembalikan nil jika cocok, dan error jika tidak cocok.
func Verify(hashedPassword, password string) error {
	// bcrypt.CompareHashAndPassword secara aman membandingkan hash dan password
	// untuk mencegah timing attacks.
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
