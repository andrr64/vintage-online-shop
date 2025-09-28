// File: pkg/auth/jwt.go
package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTService adalah service untuk mengelola JWT.
type JWTService struct {
	secretKey string
}

// NewJWTService adalah constructor untuk JWTService.
func NewJWTService(secretKey string) *JWTService {
	return &JWTService{secretKey: secretKey}
}

// Claims adalah data yang kita simpan di dalam token.
// Update: Role jadi array string (user bisa punya banyak role).
type Claims struct {
	AccountID uuid.UUID `json:"account_id"`
	Role     string    `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken membuat token JWT baru untuk user.
func (s *JWTService) GenerateToken(userID uuid.UUID, role string) (string, error) {
	// Tentukan masa berlaku token
	expirationTime := time.Now().Add(24 * time.Hour)

	// Buat claims
	claims := &Claims{
		AccountID: userID,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Buat token baru dengan claims dan metode signing
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tandatangani token dengan secret key
	return token.SignedString([]byte(s.secretKey))
}

// ValidateToken memvalidasi token JWT.
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Pastikan metode signing-nya adalah HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	// Cek apakah token valid dan ambil claims-nya
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}