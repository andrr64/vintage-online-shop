package v1

import (
	"github.com/gin-gonic/gin"
	"vintage-server/controllers"
)

func RegisterUserRoutes(r *gin.RouterGroup) {
	users := r.Group("/users")
	{
		users.POST("/register", controllers.Register)
		users.POST("/login", controllers.Login)
	}
}
