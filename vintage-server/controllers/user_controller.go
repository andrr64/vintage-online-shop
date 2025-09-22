package controllers

import (
	"net/http"
	"vintage-server/dto"
	"vintage-server/repositories"
	"vintage-server/services"

	"github.com/gin-gonic/gin"
)

var userService = services.NewUserService(
	repositories.NewUserRepository(),
)

func Register(c *gin.Context) {
	var input dto.RegisterUserDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.CommonResponse[string]{
			Message: "Invalid input",
			Success: false,
			Data:    nil,
		})
		return
	}

	userDTO, err := userService.Register(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.CommonResponse[string]{
			Message: err.Error(),
			Success: false,
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, dto.CommonResponse[dto.ResponseRegisterUserDTO]{
		Message: "User registered successfully",
		Success: true,
		Data:    &userDTO,
	})
}

func Login(c *gin.Context) {
	var input dto.LoginUserDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.CommonResponse[string]{
			Message: "Invalid input",
			Success: false,
			Data:    nil,
		})
		return
	}

	token, err := userService.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.CommonResponse[string]{
			Message: err.Error(),
			Success: false,
			Data:    nil,
		})
		return
	}

	// Set cookie JWT
	c.SetCookie(
		"access_token", // nama cookie
		token,          // value
		3600*24,        // umur (1 hari, dalam detik)
		"/",            // path
		"",             // domain (kosong = current domain)
		true,           // secure (true kalau pakai HTTPS)
		true,           // httpOnly
	)

	// Response tetap bisa kirim token kalau mau
	c.JSON(http.StatusOK, dto.CommonResponse[map[string]string]{
		Message: "Login successful",
		Success: true,
		Data:    &map[string]string{"token": token},
	})
}
