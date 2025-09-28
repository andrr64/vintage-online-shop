package middleware

import (
	"net/http"
	"vintage-server/pkg/auth"
	"vintage-server/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSvc *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("access_token")

		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		claims, err := jwtSvc.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		// Simpan info user dari token di context untuk dipakai oleh handler
		c.Set("accountID", claims.AccountID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
