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
					protected.PUT("/update-profile", accountHandler.UpdateProfile)
					protected.PUT("/update-avatar", accountHandler.UpdateAvatar)

					protected.POST("/address", accountHandler.CreateAddress)
					protected.PUT("/address", accountHandler.UpdateAddress)
					protected.GET("/address", accountHandler.GetAddresses)
					protected.DELETE("/address", accountHandler.DeleteAddress)

					protected.PUT("/address/set-primary", accountHandler.SetPrimaryAddress)
					
					protected.POST("/wishlist/:product-id", middleware.AuthRoleMiddleware("customer"), accountHandler.AddToWishlist)
				}
			}

			customer := account.Group("/customer")
			{
				customer.POST("/register", accountHandler.RegisterCustomer)
				customer.POST("/login", accountHandler.LoginCustomer)

			}
			admin := account.Group("/admin")
			{
				admin.POST("/login", accountHandler.LoginAdmin)
			}

			seller := account.Group("/seller")
			{
				seller.POST("/login", accountHandler.LoginSeller)
			}
		}

	}
	// 5. Jalankan server
	log.Printf("User Service running on port : %s", cfg.UserServicePort)
	router.Run(fmt.Sprint(":", cfg.UserServicePort))
}
