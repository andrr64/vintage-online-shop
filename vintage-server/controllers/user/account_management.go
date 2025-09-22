package user

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

// Register godoc
// @Summary Register a new user
// @Description Create a new user account with email, username, and password
// @Tags User/Account Management
// @Accept json
// @Produce json
// @Param input body dto.RegisterUserDTO true "Register user data"
// @Success 200 {object} dto.CommonResponse[dto.ResponseRegisterUserDTO]
// @Failure 400 {object} dto.CommonResponse[string]
// @Router /api/v1/user/register [post]
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

// Login godoc
// @Summary User login
// @Description Login with email and password. Returns JWT token in response and sets httpOnly cookie.
// @Tags User/Account Management
// @Accept json
// @Produce json
// @Param input body dto.LoginUserDTO true "Login credentials"
// @Success 200 {object} dto.CommonResponse[map[string]string]
// @Failure 400 {object} dto.CommonResponse[string]
// @Failure 401 {object} dto.CommonResponse[string]
// @Router /api/v1/user/login [post]
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
		false,          // secure (true kalau pakai HTTPS)
		true,           // httpOnly
	)

	// Response tetap bisa kirim token kalau mau
	c.JSON(http.StatusOK, dto.CommonResponse[map[string]string]{
		Message: "Login successful",
		Success: true,
		Data:    &map[string]string{"token": token},
	})
}

// Logout godoc
// @Summary Logout user
// @Description Invalidate user session by clearing JWT cookie
// @Tags User/Account Management
// @Accept json
// @Produce json
// @Success 200 {object} dto.CommonResponse[string]
// @Router /api/v1/user/logout [post]
func Logout(c *gin.Context) {
	// Hapus cookie JWT
	c.SetCookie(
		"access_token",
		"",
		-1, // expire segera
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, dto.CommonResponse[string]{
		Message: "User logged out successfully",
		Success: true,
		Data:    nil,
	})
}

func GetAccount(c *gin.Context) {
	// Ambil user dari context / JWT middleware
	currentUserI, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.CommonResponse[string]{
			Message: "Unauthorized",
			Success: false,
			Data:    nil,
		})
		return
	}

	currentUser := currentUserI.(map[string]interface{})

	userID := currentUser["id"].(uint)
	user, err := userService.FindByID(userID)

	response := dto.ResponseUserInfoDTO{
		Username: user.Username,
		Email:    user.Email,
		Fullname: user.Fullname,
	}

	if err != nil {
		c.JSON(http.StatusNotFound, dto.CommonResponse[any]{
			Message: "User not found",
			Success: false,
		})
		return
	}

	c.JSON(http.StatusOK, dto.CommonResponse[dto.ResponseUserInfoDTO]{
		Message: "OK",
		Success: true,
		Data:    &response,
	})
}
