package v1

import (
	adminController "vintage-server/controllers/admin"
	"vintage-server/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(r *gin.RouterGroup) {
	admin := r.Group("/admin")
	{
		admin.POST("/auth/register/:key", adminController.Register)
		admin.POST("/auth/login", adminController.Login)
		admin.POST("/auth/logout", adminController.Logout)

		admin.GET("/product_management", middlewares.AdminAuthMiddleware(), adminController.GetProducts)
	}
}
