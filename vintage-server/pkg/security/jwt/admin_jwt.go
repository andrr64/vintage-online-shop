package jwt

import (
	"fmt"
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

// CreateAdminAccessToken buat generate token JWT, output langsung string token
func CreateAdminAccessToken(id uint, username string) string {
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

func ParseAdminAccessToken(tokenStr string) (*AdminClaims, error) {
	// parse token
	token, err := jwt.ParseWithClaims(tokenStr, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		// pastikan algoritmanya sesuai (HMAC)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(admin_jwt), nil
	})

	if err != nil {
		return nil, err
	}

	// cek valid dan tipe claims
	if claims, ok := token.Claims.(*AdminClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}
