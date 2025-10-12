package middleware

import (
	"strings"
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthRoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, role, err := helper.ExtractAccountInfoFromToken(c)
		if err != nil {
			response.ErrorUnauthorized(c, "Unauthorized access.")
			c.Abort()
			return
		}

		c.Set("accountID", accountID)

		if len(allowedRoles) > 0 {
			if role == nil {
				response.ErrorForbidden(c)
				c.Abort()
				return
			}

			isAllowed := false
			for _, allowed := range allowedRoles {
				if strings.EqualFold(*role, allowed) {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				response.ErrorForbiddenRoles(c, allowedRoles...)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
