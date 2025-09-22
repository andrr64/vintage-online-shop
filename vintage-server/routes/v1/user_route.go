package v1

import (
	"github.com/gin-gonic/gin"
	userController "vintage-server/controllers/user"
)

func RegisterUserRoutes(r *gin.RouterGroup) {
	users := r.Group("/user")
	{
		users.POST("/register", userController.Register)
		users.POST("/login", userController.Login)
	}
}
