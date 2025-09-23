package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Ambil secret key dari environment variable
var admin_jwt = os.Getenv("JWT_SECRET_ADMIN")

// Payload struct
type AdminClaims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// CreateToken buat generate token JWT, output langsung string token
func CreateToken(id uint, username string) string {
	claims := AdminClaims{
		ID:       id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(admin_jwt)) // ignore error
	return signedToken
}
