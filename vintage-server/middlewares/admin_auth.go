package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"vintage-server/dto"
	"vintage-server/helpers/response_helper"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("access_token")
		if err != nil {
			response_helper.Failed[any](c, http.StatusUnauthorized, "Unauthorized", nil)
			c.Abort()
			return
		}

		secret := os.Getenv("JWT_SECRET_ADMIN")
		if secret == "" {
			c.JSON(http.StatusInternalServerError, dto.CommonResponse[string]{
				Message: "JWT secret not set",
				Success: false,
				Data:    nil,
			})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// pastikan algoritmanya sesuai
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, dto.CommonResponse[string]{
				Message: "Invalid token",
				Success: false,
				Data:    nil,
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, dto.CommonResponse[string]{
				Message: "Invalid claims",
				Success: false,
				Data:    nil,
			})
			c.Abort()
			return
		}

		// simpan user di context
		c.Set("currentUser", map[string]interface{}{
			"id":       uint(claims["id"].(float64)),
			"username": claims["username"].(string),
		})

		c.Next() // jangan lupa ini
	}
}
