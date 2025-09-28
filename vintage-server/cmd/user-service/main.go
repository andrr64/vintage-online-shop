package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	user "vintage-server/internal/service/account" // Sesuaikan path
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
	cloudinaryService, err := uploader.NewCloudinaryUploader(cfg.CloudinaryURL)
	if err != nil {
		log.Fatalf("Failed to connec to Cloudinary service: %v", err)
		return
	}

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo, cfg.JWTSecretKey, cloudinaryService)
	userHandler := user.NewHandler(userService)
	authService := auth.NewJWTService(cfg.JWTSecretKey)

	// 3. Setup Router Gin
	router := gin.Default()

	// 4. Daftar rute ke method di Handler
	api := router.Group("/api/v1") // Grup rute untuk versioning
	{
		account := api.Group("/account")
		{
			protected := account.Group("/protected")
			{
				protected.Use(middleware.AuthMiddleware(authService))
				{
					protected.POST("/logout", userHandler.Logout)
					protected.PUT("/update-profile", userHandler.UpdateProfile)
				}
			}

			customer := account.Group("/customer")
			{
				customer.POST("/register", userHandler.RegisterCustomer)
				customer.POST("/login", userHandler.LoginCustomer)

			}
			admin := account.Group("/admin")
			{
				admin.POST("/login", userHandler.LoginAdmin)
			}
		}

	}

	// 5. Jalankan server
	log.Println("User Service running on port :8081")
	router.Run(":8081")
}
