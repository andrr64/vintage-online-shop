package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Ambil secret dari environment
var user_jwt = os.Getenv("JWT_SECRET")

// Payload untuk User
type UserClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// CreateUserAccessToken generate token JWT untuk user
func CreateUserAccessToken(id uint, username string) (string, error) {
	claims := UserClaims{
		UserID:   id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(user_jwt))
}

// ParseUserAccessToken untuk verifikasi token dan ambil claims user
func ParseUserAccessToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// pastikan algoritmanya sesuai (HMAC)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(user_jwt), nil
	})

	if err != nil {
		return nil, err
	}

	// cek valid dan tipe claims
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}