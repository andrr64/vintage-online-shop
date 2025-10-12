package handler

import (
	"vintage-server/internal/account/service"

	"github.com/gin-gonic/gin"
)

type AccountHandler interface {
	// auth
	RegisterCustomer(c *gin.Context)
	LoginCustomer(c *gin.Context)
	LoginAdmin(c *gin.Context)
	LoginSeller(c *gin.Context)
	Logout(c *gin.Context)

	// Profile Management
	UpdateProfile(c *gin.Context)
	UpdateAvatar(c *gin.Context)

	// Address Management
	CreateAddress(c *gin.Context)
	GetAddresses(c *gin.Context)
	UpdateAddress(c *gin.Context)
	DeleteAddress(c *gin.Context)
	SetPrimaryAddress(c *gin.Context)

	// Wishlist Management
	AddToWishlist(c *gin.Context)
	GetWishlist(c *gin.Context)
	RemoveFromWishlist(c *gin.Context)
}

type accountHandler struct {
	services service.AccountServices
}

// Logout implements AccountHandler.
func (a *accountHandler) Logout(c *gin.Context) {
	panic("unimplemented")
}

func NewAccountHandler(services service.AccountServices) AccountHandler {
	return &accountHandler{
		services: services,
	}
}
