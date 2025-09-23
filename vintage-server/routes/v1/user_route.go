package v1

import (
	userController "vintage-server/controllers/user"
	"vintage-server/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup) {
	users := r.Group("/user")
	{
		users.POST("/auth/register", userController.Register)
		users.POST("/auth/login", userController.Login)
		users.POST("/auth/logout", userController.Logout)

		users.GET("/account", middlewares.AuthMiddleware(), userController.GetAccount)
		users.PUT("/account/profile", middlewares.AuthMiddleware(), userController.UpdateProfile)
		users.PUT("/account/password", middlewares.AuthMiddleware(), userController.UpdatePassword)
	}
}
