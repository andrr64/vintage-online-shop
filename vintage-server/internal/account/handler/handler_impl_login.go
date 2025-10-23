package handler

import (
	"vintage-server/internal/account/dto"
	"vintage-server/pkg/helper"
	"vintage-server/pkg/response"

	"github.com/gin-gonic/gin"
)

// LoginAdmin implements AccountHandler.
func (a *accountHandler) LoginAdmin(c *gin.Context) {
	var req dto.LoginRequest
	if !helper.BindJSON(c, &req) {
		return
	}
	res, err := a.services.Login.LoginAs(c.Request.Context(), "admin", req)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	c.SetCookie("access_token", res.AccessToken, 3600*24, "/", "", false, true)
	response.SuccessOK(c, res)
}

// LoginCustomer implements AccountHandler.
func (a *accountHandler) LoginCustomer(c *gin.Context) {
	var req dto.LoginRequest
	if !helper.BindJSON(c, &req) {
		return
	}
	res, err := a.services.Login.LoginAs(c.Request.Context(), "customer", req)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	c.SetCookie("access_token", res.AccessToken, 3600*24, "/", "", false, true)
	response.SuccessOK(c, res)
}

// LoginSeller implements AccountHandler.
func (a *accountHandler) LoginSeller(c *gin.Context) {
	var req dto.LoginRequest
	if !helper.BindJSON(c, &req) {
		return
	}
	res, err := a.services.Login.LoginAs(c.Request.Context(), "seller", req)
	if err != nil {
		helper.HandleError(c, err)
		return
	}
	c.SetCookie("access_token", res.AccessToken, 3600*24, "/", "", false, true)
	response.SuccessOK(c, res)
}
