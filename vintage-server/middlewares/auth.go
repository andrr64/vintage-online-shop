package middlewares

import (
	"net/http"
	"os"
	"vintage-server/dto"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"vintage-server/helpers/response_helper"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("access_token")
		if err != nil {
			response_helper.Failed[any](c, http.StatusUnauthorized, "Unauthorized", nil)
			c.Abort()
			return
		}

		secret := os.Getenv("JWT_SECRET")
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

		claims := token.Claims.(jwt.MapClaims)

		// simpan user_id di context
		currentUser := map[string]interface{}{
			"id":       uint(claims["user_id"].(float64)),
			"username": claims["username"].(string),
			// email & fullname bisa diambil dari DB nanti
		}
		c.Set("currentUser", currentUser)
		c.Next()
	}
}
