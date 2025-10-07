package main

import (
	"fmt"
	"log"
	"vintage-server/internal/database"
	repo "vintage-server/internal/repository"
	handler "vintage-server/internal/handler/shop"
	service "vintage-server/internal/service/shop"
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
	db, err := database.NewPostgres(cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	cloudinary, err := uploader.NewCloudinaryUploader(cfg.CloudinaryURL)
	if err != nil {
		log.Fatalf("Failed to connect to Cloudinary service: %v", err)
	}
	
	authService := auth.NewJWTService(cfg.JWTSecretKey)

	shopStore := repo.NewShopStore(db)
	shopService := service.NewShopService(shopStore, *authService, cloudinary)
	shopHandler := handler.NewHandler(shopService)

	router := gin.Default()

	api := router.Group("/api/v1/shop")
	{
		management := api.Group("/management")
		{
			protected := management.Group("/protected")
			{
				protected.Use(middleware.AuthMiddleware(authService))
				{
					protected.POST("/create", shopHandler.CreateShop)
					protected.PUT("/update", shopHandler.UpdateShop)
				}
			}
		}
	}

	log.Printf("Shop Service running on port :%s", cfg.ShopServicePort)
	if err := router.Run(fmt.Sprintf(":%s", cfg.ShopServicePort)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
