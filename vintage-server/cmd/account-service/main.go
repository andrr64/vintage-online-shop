package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	handler "vintage-server/internal/handler/account"
	repository "vintage-server/internal/repository/account"
	service "vintage-server/internal/service/account" // Sesuaikan path

	handlerV2 "vintage-server/internal/account/handler"
	repoV2 "vintage-server/internal/account/repository"
	serviceV2 "vintage-server/internal/account/service"

	"vintage-server/pkg/auth"
	"vintage-server/pkg/config"
	"vintage-server/pkg/middleware"
	"vintage-server/pkg/uploader"
)

func main() {
	// 1. Muat Konfigurasi
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
		return
	}

	// 2. Koneksi Database menggunakan config
	db, err := sqlx.Connect("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
		return
	}

	// 2. Merakit semua lapisan (Wiring)
	cloudinary, err := uploader.NewCloudinaryUploader(cfg.CloudinaryURL)
	if err != nil {
		log.Fatalf("Failed to connec to Cloudinary service: %v", err)
		return
	}

	storev2 := repoV2.NewAccountStore(db)
	servicev2 := serviceV2.NewAccountServices(storev2, cfg.JWTSecretKey, cloudinary)
	handlerv2 := handlerV2.NewAccountHandler(*servicev2)

	accountStore := repository.NewAccountStore(db)
	accountService := service.NewService(accountStore, cfg.JWTSecretKey, cloudinary)
	accountHandler := handler.NewAccountHandler(accountService)

	// 3. Setup Router Gin
	router := gin.Default()

	// 4. Daftar rute ke method di Handler
	api := router.Group("/api/v1") // Grup rute untuk versioning
	{
		account := api.Group("/account")
		{
			protected := account.Group("/protected")
			{
				protected.Use(middleware.AuthMiddleware(auth.NewJWTService(cfg.JWTSecretKey)))
				{
					protected.POST("/logout", accountHandler.Logout)

					protected.PUT("/profile/update-profile", accountHandler.UpdateProfile)
					protected.PUT("/profile/update-avatar", accountHandler.UpdateAvatar)

					protected.POST("/address", handlerv2.CreateAddress)
					protected.PUT("/address", accountHandler.UpdateAddress)
					protected.GET("/address", handlerv2.GetAddresses)
					protected.DELETE("/address/:address-id", handlerv2.DeleteAddress)

					protected.PUT("/address/set-primary/:address-id", handlerv2.SetPrimaryAddress)

					protected.POST("/wishlist/:product-id", middleware.AuthRoleMiddleware("customer"), handlerv2.AddToWishlist)
					protected.DELETE("/wishlist/:product-id", middleware.AuthRoleMiddleware("customer"), handlerv2.RemoveFromWishlist)
					protected.GET("/wishlist", middleware.AuthRoleMiddleware("customer"), handlerv2.GetWishlist)
				}
			}
			login := account.Group("/login")
			{
				login.POST("/customer", handlerv2.LoginCustomer)
				login.POST("/admin", handlerv2.LoginAdmin)
				login.POST("/seller", handlerv2.LoginSeller)
			}
			register := account.Group("/register")
			{
				register.POST("/customer", accountHandler.RegisterCustomer)
			}
		}

	}
	// 5. Jalankan server
	log.Printf("User Service running on port : %s", cfg.UserServicePort)
	router.Run(fmt.Sprint(":", cfg.UserServicePort))
}
