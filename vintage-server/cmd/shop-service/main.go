package main

import (
	"fmt"
	"log"
	"vintage-server/internal/database"
	"vintage-server/internal/shop/handler"
	"vintage-server/internal/shop/repository"
	"vintage-server/internal/shop/service"
	"vintage-server/pkg/auth"
	"vintage-server/pkg/config"
	"vintage-server/pkg/middleware"
	"vintage-server/pkg/uploader"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// 1. Koneksi Database (tidak berubah)
	// 1. Koneksi Database (tidak berubah)
	db, err := database.NewPostgres(cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Warning: failed to close DB: %v", err)
		}
	}()

	cloudinary, err := uploader.NewCloudinaryUploader(cfg.CloudinaryURL)
	if err != nil {
		log.Fatalf("Failed to connect to Cloudinary service: %v", err)
	}

	authService := auth.NewJWTService(cfg.JWTSecretKey)

	store := repository.NewShopStore(db)
	service := service.NewShopServices(store, cloudinary)
	handler := handler.NewShopHandler(service)

	router := gin.Default()

	api := router.Group("/api/v1/shop")
	{
		protected := api.Group("/protected")
		{
			protected.Use(middleware.AuthMiddleware(authService))
			{
				protected.POST("/create", middleware.AuthRoleMiddleware("seller", "customer"), handler.CreateShop)
				protected.PUT("/update", middleware.AuthRoleMiddleware("seller", "customer"), handler.UpdateShop)
			}
		}
	}

	log.Printf("Shop Service running on port :%s", cfg.ShopServicePort)
	if err := router.Run(fmt.Sprintf(":%s", cfg.ShopServicePort)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
