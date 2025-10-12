package main

import (
	"fmt"

	_ "github.com/lib/pq"

	"log"
	"vintage-server/internal/database"
	handler "vintage-server/internal/handler/product"
	repo "vintage-server/internal/repository"
	service "vintage-server/internal/service"
	"vintage-server/pkg/auth"
	"vintage-server/pkg/config"
	"vintage-server/pkg/middleware"
	"vintage-server/pkg/uploader"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	db, err := database.NewPostgres(cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	cloudinaryService, err := uploader.NewCloudinaryUploader(cfg.CloudinaryURL)
	if err != nil {
		log.Fatalf("Failed to connect to Cloudinary service: %v", err)
	}
	authService := auth.NewJWTService(cfg.JWTSecretKey)

	productStore := repo.NewProductStore(db)
	productService := service.NewProductService(productStore, *authService, cloudinaryService)
	productHandler := handler.NewHandler(productService)

	// 4. Setup Router (tidak berubah)
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		productGroup := api.Group("/product")
		{
			// Rute Publik
			productGroup.GET("/category", productHandler.ReadCategories)
			productGroup.GET("/brand", productHandler.ReadBrand)
			productGroup.GET("/condition", productHandler.ReadConditions)
			productGroup.GET("/:id", productHandler.GetProuctByID)

			// Rute Terproteksi
			protected := productGroup.Group(
				"/protected",
				middleware.AuthMiddleware(authService))
			{
				// Category
				protected.POST("/category", middleware.AuthRoleMiddleware("admin"), productHandler.CreateCategory)
				protected.PUT("/category/:id", middleware.AuthRoleMiddleware("admin"), productHandler.UpdateCategory) // Gunakan path param untuk konsistensi
				protected.DELETE("/category/:id", middleware.AuthRoleMiddleware("admin"), productHandler.DeleteCategory)

				// Brand
				protected.POST("/brand", middleware.AuthRoleMiddleware("admin"), productHandler.CreateBrand)
				protected.PUT("/brand/:id", middleware.AuthRoleMiddleware("admin"), productHandler.UpdateBrand)
				protected.DELETE("/brand/:id", middleware.AuthRoleMiddleware("admin"), productHandler.DeleteBrand)

				// Condition
				protected.POST("/condition", middleware.AuthRoleMiddleware("admin"), productHandler.CreateCondition)
				protected.PUT("/condition/:id", middleware.AuthRoleMiddleware("admin"), productHandler.UpdateCondition)
				protected.DELETE("/condition/:id", middleware.AuthRoleMiddleware("admin"), productHandler.DeleteCondition)

				// Product
				protected.POST("/product/create", middleware.AuthRoleMiddleware("seller"), productHandler.CreateProduct)
				protected.PUT("/product/update", middleware.AuthRoleMiddleware("seller"), productHandler.UpdateProduct)

				// Size
				protected.POST("/size", middleware.AuthRoleMiddleware("admin"), productHandler.CreateProductSize)
			}
		}
	}

	// 5. Jalankan Server (tidak berubah)
	log.Printf("Product Service running on port :%s", cfg.ProductServicePort)
	if err := router.Run(fmt.Sprintf(":%s", cfg.ProductServicePort)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
