package admin

import (
	"net/http"
	"vintage-server/dto"
	"vintage-server/helpers/response_helper"

	"github.com/gin-gonic/gin"

	jwt "vintage-server/security/jwt"

	repo "vintage-server/repositories/impl"
	adminService "vintage-server/services/admin/impl"
)

var authService = adminService.NewAdminAuthService(
	repo.NewAdminRepository(),
)

func Login(c *gin.Context) {
	var input dto.InputAdminLoginDTO
	if err := c.ShouldBind(&input); err != nil {
		response_helper.Failed[any](c, http.StatusUnauthorized, "Username or Password is Wrong", nil)
		return
	}
	// TESTING
	response := dto.ResponseAdminLogin{
		Username: "test",
		Token:    jwt.CreateToken(1, "test"),
	}
	response_helper.Success(c, &response, "OK")
}

func Register(c *gin.Context) {
	var input dto.InputAdminRegisterDTO

	if err := c.ShouldBind(&input); err != nil {
		response_helper.Failed[any](c, http.StatusBadRequest, "Invalid Input", nil)
		return
	}
	result, error := authService.Register(input)
	if error != nil {
		response_helper.Failed[any](c, http.StatusBadRequest, error.Error(), nil)
		return
	}
	response_helper.Success(c, &result, "OK")

}
