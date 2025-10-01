package main

import (
	"fmt"

	_ "github.com/lib/pq"

	"log"
	"vintage-server/internal/service/product"
	"vintage-server/pkg/auth"
	"vintage-server/pkg/config"
	"vintage-server/pkg/middleware"
	"vintage-server/pkg/uploader"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func main() {
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

	cloudinaryService, err := uploader.NewCloudinaryUploader(cfg.CloudinaryURL)
	if err != nil {
		log.Fatalf("Failed to connect to Cloudinary service: %v", err)
		return
	}
	// 1. Buat instance JWT service terlebih dahulu
	authService := auth.NewJWTService(cfg.JWTSecretKey)

	// 2. Buat instance repository
	productRepo := product.NewRepository(db)

	// 3. Suntikkan (inject) repository dan authService ke dalam product service
	productService := product.NewService(productRepo, *authService, cloudinaryService)

	// 4. Suntikkan product service ke dalam handler
	productHandler := product.NewHandler(productService)
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		product := api.Group("/product")
		{
			product.GET(("/category"), productHandler.ReadCategories)
			protected := product.Group("/protected", middleware.AuthMiddleware(authService))
			{
				protected.POST(("/category"), productHandler.CreateCategory)
				protected.PUT(("/category"), productHandler.UpdateCategory)
				protected.DELETE(("/category"), productHandler.DeleteCategory)
			}

		}
	}
	log.Printf("Product Service running on port :%s", cfg.ProductServicePort)
	router.Run(fmt.Sprintf(":%s", cfg.ProductServicePort))
}
