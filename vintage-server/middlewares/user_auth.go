package middlewares

import (
	"net/http"
	"vintage-server/dto"
	"vintage-server/helpers/response_helper"
	security_jwt "vintage-server/security/jwt" // sesuaikan path package jwt

	"github.com/gin-gonic/gin"
)

func UserAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ambil token dari cookie
		tokenStr, err := c.Cookie("access_token")
		if err != nil {
			response_helper.Failed[any](c, http.StatusUnauthorized, "Unauthorized", nil)
			c.Abort()
			return
		}

		// parse token via helper
		claims, err := security_jwt.ParseUserAccessToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.CommonResponse[string]{
				Message: "Invalid token",
				Success: false,
				Data:    nil,
			})
			c.Abort()
			return
		}

		// inject user ke context
		c.Set("currentUser", map[string]interface{}{
			"id":       claims.UserID,
			"username": claims.Username,
		})

		c.Next()
	}
}
